package cmds

import (
	"fmt"
	"io"

	"github.com/appscode/kutil"
	core_util "github.com/appscode/kutil/core/v1"
	rbac_util "github.com/appscode/kutil/rbac/v1beta1"
	"github.com/kubedb/apimachinery/apis/kubedb"
	apps "k8s.io/api/apps/v1beta1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	storage "k8s.io/api/storage/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const ServiceAccountName = operatorName

var policyRuleOperator = []rbac.PolicyRule{
	{
		APIGroups: []string{apiextensions.GroupName},
		Resources: []string{"customresourcedefinitions"},
		Verbs:     []string{"create", "delete", "get", "list"},
	},
	{
		APIGroups: []string{rbac.GroupName},
		Resources: []string{"rolebindings", "roles"},
		Verbs:     []string{"create", "delete", "get", "patch"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"services"},
		Verbs:     []string{"create", "delete", "get", "patch"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"secrets", "serviceaccounts"},
		Verbs:     []string{"create", "delete", "get", "patch"},
	},
	{
		APIGroups: []string{apps.GroupName},
		Resources: []string{"deployments", "statefulsets"},
		Verbs:     []string{"create", "delete", "get", "patch", "update"},
	},
	{
		APIGroups: []string{batch.GroupName},
		Resources: []string{"jobs"},
		Verbs:     []string{"create", "delete", "get", "list", "watch"},
	},
	{
		APIGroups: []string{storage.GroupName},
		Resources: []string{"storageclasses"},
		Verbs:     []string{"get"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"pods"},
		Verbs:     []string{"deletecollection", "get", "list", "patch", "watch"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"persistentvolumeclaims"},
		Verbs:     []string{"create", "delete", "get", "list", "patch", "watch"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"configmaps"},
		Verbs:     []string{"create", "delete", "get", "update"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"events"},
		Verbs:     []string{"create"},
	},
	{
		APIGroups: []string{core.GroupName},
		Resources: []string{"nodes"},
		Verbs:     []string{"get", "list", "watch"},
	},
	{
		APIGroups: []string{kubedb.GroupName},
		Resources: []string{rbac.ResourceAll},
		Verbs:     []string{rbac.VerbAll},
	},
	{
		APIGroups: []string{"monitoring.coreos.com"},
		Resources: []string{"servicemonitors"},
		Verbs:     []string{"create", "delete", "get", "list", "update"},
	},
}

func EnsureRBACStuff(client kubernetes.Interface, namespace string, out io.Writer) error {

	name := ServiceAccountName

	// Ensure ClusterRoles for operator
	cr, vt1, err := rbac_util.CreateOrPatchClusterRole(
		client,
		metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		func(in *rbac.ClusterRole) *rbac.ClusterRole {
			in.Labels = core_util.UpsertMap(in.Labels, operatorLabel)
			in.Rules = policyRuleOperator
			return in
		},
	)
	if err != nil {
		return err
	}
	if vt1 != kutil.VerbUnchanged {
		fmt.Fprintln(out, fmt.Sprintf(`ClusterRole "%s" successfully %v`, cr.Name, vt1))
	}

	// Ensure ServiceAccounts
	sa, vt2, err := core_util.CreateOrPatchServiceAccount(
		client,
		metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		func(in *core.ServiceAccount) *core.ServiceAccount {
			in.Labels = core_util.UpsertMap(in.Labels, operatorLabel)
			return in
		},
	)
	if err != nil {
		return err
	}
	if vt2 != kutil.VerbUnchanged {
		fmt.Fprintln(out, fmt.Sprintf(`ServiceAccount "%s" successfully %v`, sa.Name, vt2))
	}

	var roleBindingRef = rbac.RoleRef{
		APIGroup: rbac.GroupName,
		Kind:     "ClusterRole",
		Name:     name,
	}
	var roleBindingSubjects = []rbac.Subject{
		{
			Kind:      rbac.ServiceAccountKind,
			Name:      name,
			Namespace: namespace,
		},
	}

	// Ensure ClusterRoleBindings
	crb, vt3, err := rbac_util.CreateOrPatchClusterRoleBinding(
		client,
		metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		func(in *rbac.ClusterRoleBinding) *rbac.ClusterRoleBinding {
			in.Labels = core_util.UpsertMap(in.Labels, operatorLabel)
			in.RoleRef = roleBindingRef
			in.Subjects = roleBindingSubjects

			return in
		},
	)
	if err != nil {
		return err
	}
	if vt3 != kutil.VerbUnchanged {
		fmt.Fprintln(out, fmt.Sprintf(`ClusterRoleBinding "%s" successfully %v`, crb.Name, vt3))
	}
	return nil
}
