/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"
	apiutils "kubedb.dev/apimachinery/pkg/utils"
	raftutils "kubedb.dev/apimachinery/pkg/utils/raft"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	metautil "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SystemReplicationStatus struct {
	Status        string
	Details       string
	ReplayBacklog string
}

type SystemReplicationHealthSummary struct {
	AllHealthy bool
	HasActive  bool
	HasSyncing bool
	HasError   bool
	Summary    string
}

const (
	SystemReplicationStatusColumn        = "REPLICATION_STATUS"
	SystemReplicationStatusDetailsColumn = "REPLICATION_STATUS_DETAILS"
	SystemReplicationReplayBacklogColumn = "REPLAY_BACKLOG"
)

const SystemReplicationStatusQuery = `
SELECT REPLICATION_STATUS, REPLICATION_STATUS_DETAILS,
       (LAST_LOG_POSITION - REPLAYED_LOG_POSITION) AS REPLAY_BACKLOG
FROM SYS.M_SERVICE_REPLICATION`

func (HanaDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralHanaDB))
}

func (h *HanaDB) ResourceKind() string {
	return ResourceKindHanaDB
}

func (h *HanaDB) ResourcePlural() string {
	return ResourcePluralHanaDB
}

func (h *HanaDB) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", h.ResourcePlural(), SchemeGroupVersion.Group)
}

func (h *HanaDB) ResourceShortCode() string {
	return ResourceCodeHanaDB
}

func (h *HanaDB) OffshootName() string {
	return h.Name
}

func (h *HanaDB) ServiceName() string {
	return h.OffshootName()
}

func (h *HanaDB) SecondaryServiceName() string {
	return metautil.NameWithPrefix(h.ServiceName(), string(SecondaryServiceAlias))
}

func (h *HanaDB) GoverningServiceName() string {
	return metautil.NameWithSuffix(h.ServiceName(), "pods")
}

func (h *HanaDB) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", h.ServiceName(), h.Namespace)
}

func (h *HanaDB) GoverningServiceDNS(podName string) string {
	return fmt.Sprintf("%s.%s.%s.svc.%s", podName, h.GoverningServiceName(), h.Namespace, apiutils.FindDomain())
}

type hanaRaftProvider struct {
	dnsSuffix    string
	offshootName string
}

func (p hanaRaftProvider) GoverningServiceDNS(podName string) string {
	return podName + p.dnsSuffix
}

func (p hanaRaftProvider) OffshootName() string {
	return p.offshootName
}

func getHanaRaftProvider(db *HanaDB) hanaRaftProvider {
	return hanaRaftProvider{
		dnsSuffix:    fmt.Sprintf(".%s.%s.svc.%s", db.GoverningServiceName(), db.Namespace, apiutils.FindDomain()),
		offshootName: db.OffshootName(),
	}
}

// GetCurrentLeaderID queries raft leader id from a coordinator pod.
func GetCurrentLeaderID(db *HanaDB, podName, user, pass string) (uint64, error) {
	return raftutils.GetCurrentLeaderID(kubedb.HanaDBCoordinatorClientPort, db.GoverningServiceDNS(podName), user, pass)
}

// AddNodeToRaft requests raft membership add via coordinator /add-node endpoint.
func AddNodeToRaft(db *HanaDB, primaryPodName, podName string, nodeID int, user, pass string) (string, error) {
	return raftutils.AddNodeToRaft(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, kubedb.HanaDBCoordinatorPort, primaryPodName, podName, nodeID, user, pass)
}

// RemoveNodeFromRaft requests raft membership remove via coordinator /remove-node endpoint.
func RemoveNodeFromRaft(db *HanaDB, primaryPodName string, nodeID int, user, pass string) (string, error) {
	return raftutils.RemoveNodeFromRaft(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, primaryPodName, nodeID, user, pass)
}

// GetCurrentLeaderPodName returns current leader pod name by resolving raft leader id.
func GetCurrentLeaderPodName(db *HanaDB, podName, user, pass string) (string, error) {
	return raftutils.GetCurrentLeaderPodName(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, podName, user, pass)
}

func GetRaftLeaderIDWithRetries(db *HanaDB, dbPodName, user, pass string, maxTries int, retryDelay time.Duration) (int, error) {
	return raftutils.GetRaftLeaderIDWithRetries(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, dbPodName, user, pass, maxTries, retryDelay)
}

func GetRaftPrimaryNode(db *HanaDB, replicas int, user, pass string, maxTries int, retryDelay time.Duration) (int, error) {
	return raftutils.GetRaftPrimaryNode(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, replicas, user, pass, maxTries, retryDelay)
}

func AddRaftNodeWithRetries(db *HanaDB, primaryPodName, podName string, nodeID int, user, pass string, maxTries int, retryDelay time.Duration) error {
	return raftutils.AddRaftNodeWithRetries(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, kubedb.HanaDBCoordinatorPort, primaryPodName, podName, nodeID, user, pass, maxTries, retryDelay)
}

func RemoveRaftNodeWithRetries(db *HanaDB, primaryPodName string, nodeID int, user, pass string, maxTries int, retryDelay time.Duration) error {
	return raftutils.RemoveRaftNodeWithRetries(getHanaRaftProvider(db), kubedb.HanaDBCoordinatorClientPort, primaryPodName, nodeID, user, pass, maxTries, retryDelay)
}

func NewSystemReplicationStatus(status, details, replayBacklog string) SystemReplicationStatus {
	return SystemReplicationStatus{
		Status:        strings.ToUpper(strings.TrimSpace(status)),
		Details:       strings.TrimSpace(details),
		ReplayBacklog: strings.TrimSpace(replayBacklog),
	}
}

func ParseSystemReplicationStatuses(rows []map[string]string) []SystemReplicationStatus {
	statuses := make([]SystemReplicationStatus, 0, len(rows))
	for _, row := range rows {
		status := NewSystemReplicationStatus(
			row[SystemReplicationStatusColumn],
			row[SystemReplicationStatusDetailsColumn],
			row[SystemReplicationReplayBacklogColumn],
		)
		if status.Status == "" {
			continue
		}
		statuses = append(statuses, status)
	}
	return statuses
}

func EvaluateSystemReplicationHealth(statuses []SystemReplicationStatus) SystemReplicationHealthSummary {
	summary := SystemReplicationHealthSummary{
		AllHealthy: true,
	}
	if len(statuses) == 0 {
		summary.AllHealthy = false
		summary.Summary = "no replication status found"
		return summary
	}

	statusParts := make([]string, 0, len(statuses))
	for _, status := range statuses {
		replStatus := strings.ToUpper(strings.TrimSpace(status.Status))
		replDetails := strings.TrimSpace(status.Details)
		backlog := strings.TrimSpace(status.ReplayBacklog)
		if replStatus == "" {
			continue
		}

		statusPart := replStatus
		if backlog != "" && backlog != "0" {
			statusPart += "(backlog:" + backlog + ")"
		}
		if replDetails != "" && replStatus != "ACTIVE" {
			statusPart += "[" + replDetails + "]"
		}
		statusParts = append(statusParts, statusPart)

		switch replStatus {
		case "ACTIVE":
			summary.HasActive = true
		case "SYNCING", "INITIALIZING", "UNKNOWN":
			summary.HasSyncing = true
		case "ERROR":
			summary.HasError = true
		default:
			summary.HasSyncing = true
		}

		if !isSystemReplicationMemberHealthy(replStatus, replDetails) {
			summary.AllHealthy = false
		}
	}

	if len(statusParts) == 0 {
		summary.AllHealthy = false
		summary.Summary = "no replication status found"
		return summary
	}

	summary.Summary = strings.Join(statusParts, ", ")
	return summary
}

func isSystemReplicationMemberHealthy(status, details string) bool {
	if status != "ACTIVE" {
		return false
	}

	if details == "" {
		return true
	}

	normalizedDetails := strings.ToUpper(details)
	if strings.Contains(normalizedDetails, "DISCONNECT") ||
		strings.Contains(normalizedDetails, "ERROR") ||
		strings.Contains(normalizedDetails, "FAIL") ||
		strings.Contains(normalizedDetails, "SYNCING") ||
		strings.Contains(normalizedDetails, "INITIALIZ") ||
		strings.Contains(normalizedDetails, "UNKNOWN") {
		return false
	}

	// If details mention connectivity state, require connected.
	if strings.Contains(normalizedDetails, "CONNECT") &&
		!strings.Contains(normalizedDetails, "CONNECTED") {
		return false
	}

	return true
}

// GetAuthCredentialsFromSecret reads SYSTEM user/password from the auth secret.
func GetAuthCredentialsFromSecret(ctx context.Context, kc client.Client, db *HanaDB) (string, string, error) {
	secret := &core.Secret{}
	if err := kc.Get(ctx, types.NamespacedName{
		Namespace: db.Namespace,
		Name:      db.GetAuthSecretName(),
	}, secret); err != nil {
		return "", "", err
	}

	user := kubedb.HanaDBSystemUser
	if usernameBytes, ok := secret.Data[core.BasicAuthUsernameKey]; ok && len(usernameBytes) > 0 {
		user = string(usernameBytes)
	}

	if passwordBytes, ok := secret.Data[core.BasicAuthPasswordKey]; ok && len(passwordBytes) > 0 {
		return user, string(passwordBytes), nil
	}

	passwordJSON, ok := secret.Data[kubedb.HanaDBPasswordFileKey]
	if !ok {
		return "", "", fmt.Errorf("secret %s/%s missing %s key", secret.Namespace, secret.Name, kubedb.HanaDBPasswordFileKey)
	}

	var passwordData struct {
		MasterPassword string `json:"master_password"`
	}
	if err := json.Unmarshal(passwordJSON, &passwordData); err != nil {
		return "", "", fmt.Errorf("failed to parse %s in secret %s/%s: %v", kubedb.HanaDBPasswordFileKey, secret.Namespace, secret.Name, err)
	}
	if passwordData.MasterPassword == "" {
		return "", "", fmt.Errorf("master password not specified in secret %s/%s", secret.Namespace, secret.Name)
	}

	return user, passwordData.MasterPassword, nil
}

func (h *HanaDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = kubedb.ComponentDatabase
	return metautil.FilterKeys(SchemeGroupVersion.Group, selector, metautil.OverwriteKeys(nil, h.Labels, override))
}

func (h *HanaDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(h.Spec.ServiceTemplates, alias)
	return h.offshootLabels(metautil.OverwriteKeys(h.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (h *HanaDB) OffshootLabels() map[string]string {
	return h.offshootLabels(h.OffshootSelectors(), nil)
}

func (h *HanaDB) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      h.ResourceFQN(),
		metautil.InstanceLabelKey:  h.Name,
		metautil.ManagedByLabelKey: SchemeGroupVersion.Group,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (h *HanaDB) OffshootPodSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      h.ResourceFQN(),
		metautil.InstanceLabelKey:  h.Name,
		metautil.ManagedByLabelKey: SchemeGroupVersion.Group,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (h *HanaDB) PodControllerLabels(podTemplate *ofst.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Controller.Labels != nil {
		return h.offshootLabels(metautil.OverwriteKeys(h.OffshootSelectors(), extraLabels...), podTemplate.Controller.Labels)
	}
	return h.offshootLabels(metautil.OverwriteKeys(h.OffshootSelectors(), extraLabels...), nil)
}

func (h *HanaDB) PodLabels(podTemplate *ofst.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return h.offshootLabels(metautil.OverwriteKeys(h.OffshootSelectors(), extraLabels...), podTemplate.Labels)
	}
	return h.offshootLabels(metautil.OverwriteKeys(h.OffshootSelectors(), extraLabels...), nil)
}

func (h *HanaDB) ServiceAccountName() string {
	return h.OffshootName()
}

// Owner returns owner reference to resources
func (h *HanaDB) Owner() *metav1.OwnerReference {
	return metav1.NewControllerRef(h, SchemeGroupVersion.WithKind(h.ResourceKind()))
}

func (h *HanaDB) GetAuthSecretName() string {
	if h.Spec.AuthSecret != nil && h.Spec.AuthSecret.Name != "" {
		return h.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(h.OffshootName(), "auth")
}

func (h *HanaDB) GetPersistentSecrets() []string {
	var secrets []string
	secrets = append(secrets, h.GetAuthSecretName())
	return secrets
}

func (h *HanaDB) GetNameSpacedName() string {
	return h.Namespace + "/" + h.Name
}

func (r *HanaDB) ResourceSingular() string {
	return ResourceSingularHanaDB
}

type hanadbStatsService struct {
	*HanaDB
}

func (os hanadbStatsService) GetNamespace() string {
	return os.HanaDB.GetNamespace()
}

func (os hanadbStatsService) ServiceName() string {
	return os.OffshootName() + "-stats"
}

func (os hanadbStatsService) ServiceMonitorName() string {
	return os.ServiceName()
}

func (os hanadbStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return os.OffshootLabels()
}

func (os hanadbStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (os hanadbStatsService) Scheme() string {
	sc := promapi.SchemeHTTP
	return sc.String()
}

func (h *HanaDB) StatsService() mona.StatsAccessor {
	return &hanadbStatsService{h}
}

type hanadbApp struct {
	*HanaDB
}

func (r hanadbApp) Name() string {
	return r.HanaDB.Name
}

func (r hanadbApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", SchemeGroupVersion.Group, ResourceSingularHanaDB))
}

func (os hanadbStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (h HanaDB) AppBindingMeta() appcat.AppBindingMeta {
	return &hanadbApp{&h}
}

func (h *HanaDB) StatsServiceLabels() map[string]string {
	return h.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (h *HanaDB) PetSetName() string {
	return h.OffshootName()
}

func (h *HanaDB) ObserverPetSetName() string {
	return fmt.Sprintf("%s-observer", h.PetSetName())
}

func (h *HanaDB) ConfigSecretName() string {
	uid := string(h.UID)
	return metautil.NameWithSuffix(h.OffshootName(), uid[len(uid)-6:])
}

func (h *HanaDB) IsStandalone() bool {
	return h.Spec.Topology == nil || (h.Spec.Topology.Mode != nil && *h.Spec.Topology.Mode == HanaDBModeStandalone)
}

func (h *HanaDB) IsCluster() bool {
	return h.Spec.Topology != nil
}

func (h *HanaDB) IsSystemReplication() bool {
	return h.Spec.Topology != nil && h.Spec.Topology.Mode != nil &&
		*h.Spec.Topology.Mode == HanaDBModeSystemReplication
}

func (h *HanaDB) SetHealthCheckerDefaults() {
	if h.Spec.HealthChecker.PeriodSeconds == nil {
		h.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(120)
	}
	if h.Spec.HealthChecker.TimeoutSeconds == nil {
		h.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(120)
	}
	if h.Spec.HealthChecker.FailureThreshold == nil {
		h.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (h *HanaDB) SetDefaults(kc client.Client) {
	if h == nil {
		return
	}
	if h.Spec.StorageType == "" {
		h.Spec.StorageType = StorageTypeDurable
	}
	if h.Spec.DeletionPolicy == "" {
		h.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if h.Spec.PodTemplate == nil {
		h.Spec.PodTemplate = &ofst.PodTemplateSpec{}
	}

	var hanadbVersion catalog.HanaDBVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: h.Spec.Version,
	}, &hanadbVersion)
	if err != nil {
		klog.Errorf("can't get the HanaDB version object %s for %s \n", h.Spec.Version, err.Error())
		return
	}

	if h.IsStandalone() {
		if h.Spec.Replicas == nil {
			h.Spec.Replicas = pointer.Int32P(1)
		}
	}
	if h.IsSystemReplication() {
		if h.Spec.Topology.SystemReplication == nil {
			h.Spec.Topology.SystemReplication = &HanaDBSystemReplicationSpec{}
		}
		if h.Spec.Topology.SystemReplication.ReplicationMode == "" {
			h.Spec.Topology.SystemReplication.ReplicationMode = ReplicationModeSync
		}
		if h.Spec.Topology.SystemReplication.OperationMode == "" {
			h.Spec.Topology.SystemReplication.OperationMode = OperationModeLogReplay
		}
	}

	h.setDefaultContainerSecurityContext(&hanadbVersion, h.Spec.PodTemplate)

	h.SetHealthCheckerDefaults()

	h.setDefaultContainerResourceLimits(h.Spec.PodTemplate)
}

func (h *HanaDB) setDefaultContainerSecurityContext(hanadbVersion *catalog.HanaDBVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = hanadbVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.HanaDBContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.HanaDBContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}

	h.assignDefaultContainerSecurityContext(hanadbVersion, container.SecurityContext)

	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (h *HanaDB) assignDefaultContainerSecurityContext(hanadbVersion *catalog.HanaDBVersion, sc *core.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}

	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = hanadbVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = hanadbVersion.Spec.SecurityContext.RunAsGroup
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (h *HanaDB) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.HanaDBContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesHanaDB)
	}
}

func (h *HanaDB) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(h.Namespace), labels.SelectorFromSet(h.OffshootLabels()), expectedItems)
}
