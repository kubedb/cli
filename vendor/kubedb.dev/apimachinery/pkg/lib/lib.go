/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"context"
	"time"

	"kubedb.dev/apimachinery/apis/kubedb"

	cm_api "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cm_cs "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/typed/certmanager/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
)

const (
	DefaultCertRenewBefore = time.Minute * 5
	DefaultAuthRenewBefore = time.Minute * 10
)

type LabeledObject interface {
	metav1.Object
	OffshootLabels() map[string]string
}

func SecretExists(kc kubernetes.Interface, meta metav1.ObjectMeta) bool {
	_, err := kc.CoreV1().Secrets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
	return err == nil
}

func AddOwnerReferenceToSecret(kc kubernetes.Interface, gvk schema.GroupVersionKind, db LabeledObject, secret metav1.ObjectMeta) error {
	ref := metav1.NewControllerRef(db, gvk)

	_, _, err := core_util.CreateOrPatchSecret(context.TODO(), kc, secret, func(in *core.Secret) *core.Secret {
		in.Labels = meta.OverwriteKeys(in.Labels, db.OffshootLabels())
		core_util.EnsureOwnerReference(in, ref)

		return in
	}, metav1.PatchOptions{})

	return err
}

func IsCertificateConditionTrue(conditions []cm_api.CertificateCondition, condType cm_api.CertificateConditionType) bool {
	for i := range conditions {
		if conditions[i].Type == condType && conditions[i].Status == cmmeta.ConditionTrue {
			return true
		}
	}
	return false
}

func HasCertificateCondition(conditions []cm_api.CertificateCondition, condType cm_api.CertificateConditionType) bool {
	for i := range conditions {
		if conditions[i].Type == condType {
			return true
		}
	}
	return false
}

func GetCertificateCondition(conditions []cm_api.CertificateCondition, condType cm_api.CertificateConditionType) (int, *cm_api.CertificateCondition) {
	for i := range conditions {
		c := conditions[i]
		if c.Type == condType {
			return i, &c
		}
	}
	return -1, nil
}

func UpsertCertificateCondition(conditions []cm_api.CertificateCondition, newCond cm_api.CertificateCondition) []cm_api.CertificateCondition {
	for idx, cond := range conditions {
		if cond.Type != newCond.Type {
			continue
		}
		if cond.Status == newCond.Status {
			newCond.LastTransitionTime = cond.LastTransitionTime
		}

		conditions[idx] = newCond
		return conditions
	}

	conditions = append(conditions, newCond)
	return conditions
}

func GetIssuerObjectRef(tlsConfig *kmapi.TLSConfig, alias string) cmmeta.ObjectReference {
	issuer := tlsConfig.IssuerRef
	if _, cert := kmapi.GetCertificate(tlsConfig.Certificates, alias); cert != nil {
		if cert.IssuerRef != nil {
			issuer = cert.IssuerRef
		}
	}

	return cmmeta.ObjectReference{
		Name:  issuer.Name,
		Kind:  issuer.Kind,
		Group: pointer.String(issuer.APIGroup),
	}
}

func getDBIssuers(tlsConfig kmapi.TLSConfig) []*core.TypedLocalObjectReference {
	issuers := make([]*core.TypedLocalObjectReference, 0)

	if tlsConfig.IssuerRef != nil {
		issuers = append(issuers, tlsConfig.IssuerRef)
	}

	for _, cert := range tlsConfig.Certificates {
		if cert.IssuerRef != nil {
			issuers = append(issuers, cert.IssuerRef)
		}
	}

	return issuers
}

func DBsForIssuer(dc dynamic.Interface, gvr schema.GroupVersionResource, issuer *cm_api.Issuer) ([]cache.ExplicitKey, error) {
	var out []cache.ExplicitKey

	objs, err := dc.Resource(gvr).Namespace(issuer.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range objs.Items {
		f, _, err := unstructured.NestedFieldNoCopy(obj.Object, "spec", "tls")
		if err != nil {
			return nil, err
		}

		if tlsObject, ok := f.(map[string]any); ok {
			var tls kmapi.TLSConfig
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(tlsObject, &tls)
			if err != nil {
				return nil, err
			}
			issuers := getDBIssuers(tls)
			for _, cur := range issuers {
				if issuer.Namespace == obj.GetNamespace() &&
					issuer.Name == cur.Name &&
					cur.APIGroup != nil &&
					*cur.APIGroup == cm_api.SchemeGroupVersion.Group &&
					cur.Kind == cm_api.IssuerKind {
					out = append(out, cache.ExplicitKey(obj.GetNamespace()+"/"+obj.GetName()))
					break
				}
			}
		}
	}

	return out, nil
}

func DBsForClusterIssuer(dc dynamic.Interface, gvr schema.GroupVersionResource, clusterIssuer *cm_api.ClusterIssuer) ([]cache.ExplicitKey, error) {
	var out []cache.ExplicitKey

	objs, err := dc.Resource(gvr).Namespace(core.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, obj := range objs.Items {
		f, _, err := unstructured.NestedFieldNoCopy(obj.Object, "spec", "tls")
		if err != nil {
			return nil, err
		}
		if tlsObject, ok := f.(map[string]any); ok {
			var tls kmapi.TLSConfig
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(tlsObject, &tls)
			if err != nil {
				return nil, err
			}
			issuers := getDBIssuers(tls)
			for _, cur := range issuers {
				if clusterIssuer.Name == cur.Name &&
					cur.APIGroup != nil &&
					*cur.APIGroup == cm_api.SchemeGroupVersion.Group &&
					cur.Kind == cm_api.ClusterIssuerKind {
					out = append(out, cache.ExplicitKey(obj.GetName()))
					break
				}
			}
		}
	}

	return out, nil
}

func DBForSecret(certGetter cm_cs.CertificatesGetter, kind string, s *core.Secret) (cache.ExplicitKey, error) {
	ctrl := metav1.GetControllerOf(s)
	if ctrl != nil {
		ok, err := core_util.IsOwnerOfGroupKind(ctrl, kubedb.GroupName, kind)
		if err != nil || !ok {
			return "", err
		}
		return cache.ExplicitKey(s.Namespace + "/" + ctrl.Name), nil
	}

	certName, ok := s.Annotations[cm_api.CertificateNameKey]
	if !ok {
		return "", nil
	}

	cert, err := certGetter.Certificates(s.Namespace).Get(context.TODO(), certName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if cert.Spec.SecretName != s.Name {
		return "", nil
	}

	certCtrl := metav1.GetControllerOf(cert)
	ok, err = core_util.IsOwnerOfGroupKind(certCtrl, kubedb.GroupName, kind)
	if err != nil || !ok {
		return "", err
	}
	return cache.ExplicitKey(s.Namespace + "/" + certCtrl.Name), nil
}

func DBForService(kind string, s *core.Service) cache.ExplicitKey {
	ctrl := metav1.GetControllerOf(s)
	ok, err := core_util.IsOwnerOfGroupKind(ctrl, kubedb.GroupName, kind)
	if err != nil || !ok {
		return ""
	}
	return cache.ExplicitKey(s.Namespace + "/" + ctrl.Name)
}

func ServiceDNS(svc metav1.ObjectMeta) []string {
	return []string{
		svc.Name + "." + svc.Namespace + ".svc",
		svc.Name,
	}
}

func ServiceHosts(getter v1.ServicesGetter, svc metav1.ObjectMeta) (sets.Set[string], sets.Set[string], error) {
	s, err := getter.Services(svc.Namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	dnsNames := sets.New[string](ServiceDNS(svc)...)
	ipAddresses := sets.New[string]()
	serviceIngress := s.Status.LoadBalancer.Ingress
	for _, ingres := range serviceIngress {
		if ingres.Hostname != "" {
			dnsNames.Insert(ingres.Hostname)
		} else if ingres.IP != "" {
			ipAddresses.Insert(ingres.IP)
		}
	}
	return dnsNames, ipAddresses, nil
}

func IsServiceReady(getter v1.ServicesGetter, svc metav1.ObjectMeta) bool {
	s, err := getter.Services(svc.Namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	if s.Spec.Type == core.ServiceTypeLoadBalancer {
		return len(s.Status.LoadBalancer.Ingress) > 0
	}
	return true
}

const VolumeExpansionAnnotationCoreKey = "volumeexpansion.ops.kubedb.com"

func VolumeExpansionAnnotationKey(podName string) string {
	return VolumeExpansionAnnotationCoreKey + "/" + podName
}

const StorageMigrationAnnotationCoreKey = "storagemigration.ops.kubedb.com"

func StorageMigrationAnnotationKey(podName string) string {
	return StorageMigrationAnnotationCoreKey + "/" + podName
}

func BackupOrRestoreRunningForDB(KBClient client.Client, stashClient scs.StashV1beta1Interface, dbMeta metav1.ObjectMeta, kind string) (bool, string, error) {
	if stashOperatorExist(KBClient) {
		if skip, msg, err := stashBackupOrRestoreRunningForDB(stashClient, dbMeta); err != nil || skip {
			return skip, msg, err
		}
	}
	if kubeStashOperatorExist(KBClient) {
		if skip, msg, err := kubeStashBackupOrRestoreRunningForDB(KBClient, dbMeta, kind); err != nil || skip {
			return skip, msg, err
		}
	}

	return false, "", nil
}

func PauseBackupConfiguration(KBClient client.Client, stashClient scs.StashV1beta1Interface, dbMeta metav1.ObjectMeta, pausedBackups []kmapi.TypedObjectReference, kind string, opsGeneration int64) ([]kmapi.TypedObjectReference, []kmapi.Condition, error) {
	stashBackups := make([]kmapi.TypedObjectReference, 0)
	stashConditions := make([]kmapi.Condition, 0)
	kubeStashBackups := make([]kmapi.TypedObjectReference, 0)
	kubeStashConditions := make([]kmapi.Condition, 0)
	var err error

	if stashOperatorExist(KBClient) {
		stashBackups, stashConditions, err = pauseStashBackupConfiguration(stashClient, dbMeta, pausedBackups, opsGeneration)
		if err != nil {
			return nil, nil, err
		}
	}
	if kubeStashOperatorExist(KBClient) {
		kubeStashBackups, kubeStashConditions, err = pauseKubeStashBackupConfiguration(KBClient, dbMeta, pausedBackups, kind, opsGeneration)
		if err != nil {
			return nil, nil, err
		}
	}

	return append(stashBackups, kubeStashBackups...), append(stashConditions, kubeStashConditions...), nil
}

func ResumeBackupConfiguration(KBClient client.Client, stashClient scs.StashV1beta1Interface, pausedBackups []kmapi.TypedObjectReference, opsGeneration int64) ([]kmapi.Condition, error) {
	stashConditions := make([]kmapi.Condition, 0)
	kubeStashConditions := make([]kmapi.Condition, 0)
	var err error

	if kubeStashOperatorExist(KBClient) {
		stashConditions, err = resumeStashBackupConfiguration(stashClient, pausedBackups, opsGeneration)
		if err != nil {
			return nil, err
		}
	}

	if kubeStashOperatorExist(KBClient) {
		kubeStashConditions, err = resumeKubeStashBackupConfiguration(KBClient, pausedBackups, opsGeneration)
		if err != nil {
			return nil, err
		}
	}

	return append(stashConditions, kubeStashConditions...), nil
}
