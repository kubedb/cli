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

package v1alpha1

import (
	"context"
	"fmt"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kubestash.dev/apimachinery/apis"
	storageapi "kubestash.dev/apimachinery/apis/storage/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sync"
)

// log is for logging in this package.
var backupconfigurationlog = logf.Log.WithName("backupconfiguration-resource")

func (b *BackupConfiguration) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(b).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-core-kubestash-com-v1alpha1-backupconfiguration,mutating=true,failurePolicy=fail,sideEffects=None,groups=core.kubestash.com,resources=backupconfigurations,verbs=create;update,versions=v1alpha1,name=mbackupconfiguration.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &BackupConfiguration{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (b *BackupConfiguration) Default() {
	backupconfigurationlog.Info("default", "name", b.Name)

	c := apis.GetRuntimeClient()

	if len(b.Spec.Backends) == 0 {
		b.setDefaultBackend(context.Background(), c)
	}

	b.setDefaultRetentionPolicy(context.Background(), c)
}

func (b *BackupConfiguration) setDefaultBackend(ctx context.Context, c client.Client) {
	bs := b.getDefaultStorage(ctx, c)
	if bs == nil {
		return
	}

	b.Spec.Backends = []BackendReference{
		{
			StorageRef: &kmapi.ObjectReference{
				Name:      bs.Name,
				Namespace: bs.Namespace,
			},
		},
	}
}

func (b *BackupConfiguration) setDefaultRetentionPolicy(ctx context.Context, c client.Client) {
	for i, backend := range b.Spec.Backends {
		if backend.RetentionPolicy == nil {
			rp := b.getDefaultRetentionPolicy(ctx, c)
			if rp == nil {
				return
			}

			b.Spec.Backends[i].RetentionPolicy = &kmapi.ObjectReference{
				Name:      rp.Name,
				Namespace: rp.Namespace,
			}
		}
	}
}

func (b *BackupConfiguration) getDefaultStorage(ctx context.Context, c client.Client) *storageapi.BackupStorage {
	bsList := &storageapi.BackupStorageList{}
	if err := c.List(ctx, bsList); err != nil {
		backupconfigurationlog.Error(err, "failed to list BackupStorage")
		return nil
	}

	// Check if there is any same namespace level default BackupStorage
	for _, bs := range bsList.Items {
		if bs.Namespace == b.Namespace &&
			bs.Spec.Default &&
			*bs.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromSame {
			return &bs
		}
	}

	// Check if there is any selector level default BackupStorage
	ns := &core.Namespace{ObjectMeta: v1.ObjectMeta{Name: b.Namespace}}
	if err := c.Get(ctx, client.ObjectKeyFromObject(ns), ns); err != nil {
		backupconfigurationlog.Error(err, "failed to get namespace")
		return nil
	}
	for _, bs := range bsList.Items {
		if bs.Spec.Default &&
			*bs.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromSelector &&
			selectorMatches(bs.Spec.UsagePolicy.AllowedNamespaces.Selector, ns.Labels) {
			return &bs
		}
	}

	// Check if there is any all namespace level default BackupStorage
	for _, bs := range bsList.Items {
		if bs.Spec.Default &&
			*bs.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromAll {
			return &bs
		}
	}

	backupconfigurationlog.Error(fmt.Errorf("no default BackupStorage found"), "")
	return nil
}

func (b *BackupConfiguration) getDefaultRetentionPolicy(ctx context.Context, c client.Client) *storageapi.RetentionPolicy {
	rpList := &storageapi.RetentionPolicyList{}
	if err := c.List(ctx, rpList); err != nil {
		backupconfigurationlog.Error(err, "failed to list RetentionPolicy")
		return nil
	}

	// Check if there is any same namespace level default RetentionPolicy
	for _, rp := range rpList.Items {
		if rp.Namespace == b.Namespace &&
			rp.Spec.Default &&
			*rp.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromSame {
			return &rp
		}
	}

	// Check if there is any selector level default RetentionPolicy
	ns := &core.Namespace{ObjectMeta: v1.ObjectMeta{Name: b.Namespace}}
	if err := c.Get(ctx, client.ObjectKeyFromObject(ns), ns); err != nil {
		backupconfigurationlog.Error(err, "failed to get namespace")
		return nil
	}
	for _, rp := range rpList.Items {
		if rp.Spec.Default &&
			*rp.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromSelector &&
			selectorMatches(rp.Spec.UsagePolicy.AllowedNamespaces.Selector, ns.Labels) {
			return &rp
		}
	}

	// Check if there is any all namespace level default RetentionPolicy
	for _, rp := range rpList.Items {
		if rp.Spec.Default &&
			*rp.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromAll {
			return &rp
		}
	}

	backupconfigurationlog.Error(fmt.Errorf("no default RetentionPolicy found"), "")
	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-core-kubestash-com-v1alpha1-backupconfiguration,mutating=false,failurePolicy=fail,sideEffects=None,groups=core.kubestash.com,resources=backupconfigurations,verbs=create;update,versions=v1alpha1,name=vbackupconfiguration.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &BackupConfiguration{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (b *BackupConfiguration) ValidateCreate() (admission.Warnings, error) {
	backupconfigurationlog.Info("validate create", apis.KeyName, b.Name)

	c := apis.GetRuntimeClient()

	if err := b.validateBackends(); err != nil {
		return nil, err
	}

	if err := b.validateSessions(context.Background(), c); err != nil {
		return nil, err
	}

	if err := b.validateBackendsAgainstUsagePolicy(context.Background(), c); err != nil {
		return nil, err
	}

	return nil, b.validateHookTemplatesAgainstUsagePolicy(context.Background(), c)
}

var (
	rc   client.Client
	once sync.Once
)

func (b *BackupConfiguration) validateBackends() error {
	if len(b.Spec.Backends) == 0 {
		return fmt.Errorf("backend can not be empty")
	}

	if err := b.validateBackendNameUnique(); err != nil {
		return err
	}

	return b.validateBackendReferences()
}

func (b *BackupConfiguration) validateBackendNameUnique() error {
	backendMap := make(map[string]struct{})

	for _, backend := range b.Spec.Backends {
		if _, ok := backendMap[backend.Name]; ok {
			return fmt.Errorf("duplicate backend name found: %q. Please choose a different backend name", backend.Name)
		}
		backendMap[backend.Name] = struct{}{}
	}
	return nil
}

func (b *BackupConfiguration) validateBackendReferences() error {
	for _, backend := range b.Spec.Backends {
		if backend.RetentionPolicy == nil {
			return fmt.Errorf("no RetentionPolicy is found for backend: %q. Please add a RetentionPolicy", backend.Name)
		}

		if backend.StorageRef == nil {
			return fmt.Errorf("no storage reference is found for backend: %q. Please provide a storage reference", backend.Name)
		}
	}
	return nil
}

func (b *BackupConfiguration) validateSessions(ctx context.Context, c client.Client) error {
	if len(b.Spec.Sessions) == 0 {
		return fmt.Errorf("sessions can not be empty")
	}

	if err := b.validateSessionNameUnique(); err != nil {
		return err
	}

	for _, session := range b.Spec.Sessions {
		if err := b.validateRepositoryNameUnique(session); err != nil {
			return err
		}

		if err := b.validateSessionConfig(session); err != nil {
			return err
		}

		if err := b.validateAddonInfo(session); err != nil {
			return err
		}

		if err := b.validateRepositories(ctx, c, session); err != nil {
			return err
		}
	}

	if err := b.validateUniqueRepoDir(ctx, c); err != nil {
		return err
	}

	return nil
}

func (b *BackupConfiguration) validateSessionNameUnique() error {
	sessionMap := make(map[string]struct{})

	for _, session := range b.Spec.Sessions {
		if session.Name == "" {
			return fmt.Errorf("session name can not be empty")
		}

		if _, ok := sessionMap[session.Name]; ok {
			return fmt.Errorf("duplicate session name found: %q. Please choose a different session name", session.Name)
		}
		sessionMap[session.Name] = struct{}{}
	}
	return nil
}

func (b *BackupConfiguration) validateSessionConfig(session Session) error {
	if session.Scheduler == nil {
		return fmt.Errorf("scheduler is empty for session: %q. Please provide scheduler", session.Name)
	}

	return nil
}

func (b *BackupConfiguration) validateAddonInfo(session Session) error {
	if session.Addon == nil {
		return fmt.Errorf("addon info is empty for session: %q. Please provide addon info", session.Name)
	}

	if session.Addon.Name == "" {
		return fmt.Errorf("addon name is empty for session: %q. Please provide a valid addon name", session.Name)
	}

	if len(session.Addon.Tasks) == 0 {
		return fmt.Errorf("addon tasks are not provided for session: %q. Please provide valid tasks", session.Name)
	}

	for _, task := range session.Addon.Tasks {
		if task.Name == "" {
			return fmt.Errorf("addon task name is empty for session: %q. Please provide valid task name", session.Name)
		}
	}

	return nil
}

func (b *BackupConfiguration) validateRepositories(ctx context.Context, c client.Client, session Session) error {
	if len(session.Repositories) == 0 {
		return fmt.Errorf("no repository found for session: %q. Please provide atleast one repository", session.Name)
	}

	for _, repo := range session.Repositories {
		if repo.Backend != "" &&
			!b.backendMatched(repo) {
			return fmt.Errorf("backend %q for repository %q doesn't match with any of the given backends", repo.Backend, repo.Name)
		}

		if repo.Directory == "" {
			return fmt.Errorf("directory is not provided for repository: %q. Please provide a directory", repo.Name)
		}

		existingRepo, err := b.getRepository(ctx, c, repo.Name)
		if err != nil {
			if kerr.IsNotFound(err) {
				continue
			}
			return err
		}

		if !targetMatched(&existingRepo.Spec.AppRef, b.GetTargetRef()) {
			return fmt.Errorf("repository %q already exists in the cluster with a different target reference. Please, choose a different repository name", repo.Name)
		}

		if !storageRefMatched(b.GetStorageRef(repo.Backend), &existingRepo.Spec.StorageRef) {
			return fmt.Errorf("repository %q already exists in the cluster with a different storage reference. Please, choose a different repository name", repo.Name)
		}
	}

	return nil
}

func (b *BackupConfiguration) validateUniqueRepoDir(ctx context.Context, c client.Client) error {
	if err := b.validateRepoDirectories(); err != nil {
		return err
	}

	return b.validateExistingRepoDir(ctx, c)
}

func (b *BackupConfiguration) validateRepoDirectories() error {
	mapRepoDir := make(map[string]map[string]string)
	for _, session := range b.Spec.Sessions {
		for _, repo := range session.Repositories {
			if repoInfo, ok := mapRepoDir[repo.Directory]; ok && repoInfo[repo.Backend] != repo.Name {
				return fmt.Errorf("repository %q already using the directory %q. Please, choose a different directory for repository %q", repoInfo[repo.Backend], repo.Directory, repo.Name)
			}
			mapRepoDir[repo.Directory] = map[string]string{
				repo.Backend: repo.Name,
			}
		}
	}
	return nil
}

func (b *BackupConfiguration) validateExistingRepoDir(ctx context.Context, c client.Client) error {
	for _, session := range b.Spec.Sessions {
		for _, repo := range session.Repositories {
			if err := b.validateRepoDirExistence(ctx, c, repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *BackupConfiguration) validateRepoDirExistence(ctx context.Context, c client.Client, repo RepositoryInfo) error {
	existingRepos, err := b.getRepositories(ctx, c)
	if err != nil {
		return err
	}

	for _, existingRepo := range existingRepos {
		if existingRepo.Name == repo.Name &&
			existingRepo.Namespace == b.Namespace {
			continue
		}

		if storageRefMatched(b.GetStorageRef(repo.Backend), &existingRepo.Spec.StorageRef) &&
			existingRepo.Spec.Path == repo.Directory {
			return fmt.Errorf("repository %q already exists in the cluster with the same directory. Please, choose a different directory for repository %q", existingRepo.Name, repo.Name)
		}
	}
	return nil
}

func storageRefMatched(b1, b2 *kmapi.ObjectReference) bool {
	return b1.Name == b2.Name && b1.Namespace == b2.Namespace
}

func targetMatched(t1, t2 *kmapi.TypedObjectReference) bool {
	return t1.APIGroup == t2.APIGroup &&
		t1.Kind == t2.Kind &&
		t1.Namespace == t2.Namespace &&
		t1.Name == t2.Name
}

func (b *BackupConfiguration) validateRepositoryNameUnique(session Session) error {
	repoMap := make(map[string]struct{})
	for _, repo := range session.Repositories {
		if _, ok := repoMap[repo.Name]; ok {
			return fmt.Errorf("duplicate repository name found: %q. Please choose a different repository name", repo.Name)
		}
		repoMap[repo.Name] = struct{}{}
	}
	return nil
}

func (b *BackupConfiguration) backendMatched(repo RepositoryInfo) bool {
	for _, b := range b.Spec.Backends {
		if b.Name == repo.Backend {
			return true
		}
	}
	return false
}

func (b *BackupConfiguration) getRepository(ctx context.Context, c client.Client, name string) (*storageapi.Repository, error) {
	repo := &storageapi.Repository{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: b.Namespace,
		},
	}

	if err := c.Get(ctx, client.ObjectKeyFromObject(repo), repo); err != nil {
		return nil, err
	}
	return repo, nil
}

func (b *BackupConfiguration) getRepositories(ctx context.Context, c client.Client) ([]storageapi.Repository, error) {
	repos := &storageapi.RepositoryList{}

	if err := c.List(ctx, repos); err != nil {
		return nil, err
	}
	return repos.Items, nil
}

func (b *BackupConfiguration) validateBackendsAgainstUsagePolicy(ctx context.Context, c client.Client) error {
	for _, backend := range b.Spec.Backends {
		ns := &core.Namespace{ObjectMeta: v1.ObjectMeta{Name: b.Namespace}}
		if err := c.Get(ctx, client.ObjectKeyFromObject(ns), ns); err != nil {
			return err
		}

		if err := b.validateStorageUsagePolicy(ctx, c, backend.StorageRef, ns); err != nil {
			return err
		}

		if err := b.validateRetentionPolicyUsagePolicy(ctx, c, backend.RetentionPolicy, ns); err != nil {
			return err
		}
	}
	return nil
}

func (b *BackupConfiguration) validateStorageUsagePolicy(ctx context.Context, c client.Client, ref *kmapi.ObjectReference, ns *core.Namespace) error {
	bs, err := b.getBackupStorage(ctx, c, ref)
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		}
		return err
	}

	if !bs.UsageAllowed(ns) {
		return fmt.Errorf("namespace %q is not allowed to refer BackupStorage %s/%s. Please, check the `usagePolicy` of the BackupStorage", b.Namespace, bs.Name, bs.Namespace)
	}
	return nil
}

func (b *BackupConfiguration) getBackupStorage(ctx context.Context, c client.Client, ref *kmapi.ObjectReference) (*storageapi.BackupStorage, error) {
	bs := &storageapi.BackupStorage{
		ObjectMeta: v1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
	}

	if bs.Namespace == "" {
		bs.Namespace = b.Namespace
	}

	if err := c.Get(ctx, client.ObjectKeyFromObject(bs), bs); err != nil {
		return nil, err
	}
	return bs, nil
}

func (b *BackupConfiguration) validateRetentionPolicyUsagePolicy(ctx context.Context, c client.Client, ref *kmapi.ObjectReference, ns *core.Namespace) error {
	rp, err := b.getRetentionPolicy(ctx, c, ref)
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		}
		return err
	}

	if !rp.UsageAllowed(ns) {
		return fmt.Errorf("namespace %q is not allowed to refer RetentionPolicy %s/%s. Please, check the `usagePolicy` of the RetentionPolicy", b.Namespace, rp.Name, rp.Namespace)
	}
	return nil
}

func (b *BackupConfiguration) getRetentionPolicy(ctx context.Context, c client.Client, ref *kmapi.ObjectReference) (*storageapi.RetentionPolicy, error) {
	rp := &storageapi.RetentionPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
	}

	if rp.Namespace == "" {
		rp.Namespace = b.Namespace
	}

	if err := c.Get(ctx, client.ObjectKeyFromObject(rp), rp); err != nil {
		return nil, err
	}
	return rp, nil
}

func (b *BackupConfiguration) validateHookTemplatesAgainstUsagePolicy(ctx context.Context, c client.Client) error {
	hookTemplates := b.getHookTemplates()
	for _, ht := range hookTemplates {
		err := c.Get(ctx, client.ObjectKeyFromObject(&ht), &ht)
		if err != nil {
			if kerr.IsNotFound(err) {
				continue
			}
			return err
		}

		ns := &core.Namespace{ObjectMeta: v1.ObjectMeta{Name: b.Namespace}}
		if err := c.Get(ctx, client.ObjectKeyFromObject(ns), ns); err != nil {
			return err
		}

		if !ht.UsageAllowed(ns) {
			return fmt.Errorf("namespace %q is not allowed to refer HookTemplate %s/%s. Please, check the `usagePolicy` of the HookTemplate", b.Namespace, ht.Name, ht.Namespace)
		}
	}
	return nil
}

func (b *BackupConfiguration) getHookTemplates() []HookTemplate {
	var hookTemplates []HookTemplate
	for _, session := range b.Spec.Sessions {
		if session.Hooks != nil {
			hookTemplates = append(hookTemplates, b.getHookTemplatesFromHookInfo(session.Hooks.PreBackup)...)
			hookTemplates = append(hookTemplates, b.getHookTemplatesFromHookInfo(session.Hooks.PostBackup)...)
		}
	}
	return hookTemplates
}

func (b *BackupConfiguration) getHookTemplatesFromHookInfo(hooks []HookInfo) []HookTemplate {
	var hookTemplates []HookTemplate
	for _, hook := range hooks {
		if hook.HookTemplate != nil {
			hookTemplates = append(hookTemplates, HookTemplate{
				ObjectMeta: v1.ObjectMeta{
					Name:      hook.HookTemplate.Name,
					Namespace: hook.HookTemplate.Namespace,
				},
			})
		}
	}
	return hookTemplates
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (b *BackupConfiguration) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	backupconfigurationlog.Info("validate update", apis.KeyName, b.Name)
	c := apis.GetRuntimeClient()
	if err := b.validateBackends(); err != nil {
		return nil, err
	}

	if err := b.validateSessions(context.Background(), c); err != nil {
		return nil, err
	}

	if err := b.validateBackendsAgainstUsagePolicy(context.Background(), c); err != nil {
		return nil, err
	}

	return nil, b.validateHookTemplatesAgainstUsagePolicy(context.Background(), c)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (b *BackupConfiguration) ValidateDelete() (admission.Warnings, error) {
	backupconfigurationlog.Info("validate delete", apis.KeyName, b.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
