package roles

import (
	"fmt"
	"io"

	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/docker"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	rbac "k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

const ServiceAccountName = docker.OperatorName

var policyRuleOperator = []rbac.PolicyRule{
	{
		APIGroups: []string{extensions.GroupName},
		Resources: []string{"thirdpartyresources"},
		Verbs:     []string{"get", "create"},
	},
	{
		APIGroups: []string{rbac.GroupName},
		Resources: []string{"roles"},
		Verbs:     []string{"get", "create", "delete"},
	},
	{
		APIGroups: []string{rbac.GroupName},
		Resources: []string{"rolebindings"},
		Verbs:     []string{"create", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"serviceaccounts"},
		Verbs:     []string{"create", "delete"},
	},
	{
		APIGroups: []string{apps.GroupName},
		Resources: []string{"statefulsets"},
		Verbs:     []string{"get", "create", "update", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"services", "secrets"},
		Verbs:     []string{"get", "create", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"endpoints"},
		Verbs:     []string{"get"},
	},
	{
		APIGroups: []string{batch.GroupName},
		Resources: []string{"jobs"},
		Verbs:     []string{"get", "create", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"pods"},
		Verbs:     []string{"get", "list", "delete", "deletecollection"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"persistentvolumeclaims"},
		Verbs:     []string{"list", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"events"},
		Verbs:     []string{"create"},
	},
	{
		APIGroups: []string{tapi.GroupName},
		Resources: []string{rbac.ResourceAll},
		Verbs:     []string{rbac.VerbAll},
	},
	{
		APIGroups: []string{"monitoring.coreos.com"},
		Resources: []string{"servicemonitors"},
		Verbs:     []string{"get", "create", "update"},
	},
}

func EnsureRBACStuff(client kubernetes.Interface, namespace string, out io.Writer) error {
	name := ServiceAccountName
	// Ensure ClusterRoles for operator
	clusterRoleOperator, err := client.RbacV1beta1().ClusterRoles().Get(name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		// Create new one
		role := &rbac.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Rules: policyRuleOperator,
		}
		if _, err := client.RbacV1beta1().ClusterRoles().Create(role); err != nil {
			return err
		}
		fmt.Fprintln(out, "Successfully created cluster role.")
	} else {
		// Update existing one
		clusterRoleOperator.Rules = policyRuleOperator
		if _, err := client.RbacV1beta1().ClusterRoles().Update(clusterRoleOperator); err != nil {
			return err
		}
		fmt.Fprintln(out, "Successfully updated cluster role.")
	}

	// Ensure ServiceAccounts
	if _, err := client.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{}); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		sa := &apiv1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
		if _, err := client.CoreV1().ServiceAccounts(namespace).Create(sa); err != nil {
			return err
		}
		fmt.Fprintln(out, "Successfully created service account.")
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
	roleBinding, err := client.RbacV1beta1().ClusterRoleBindings().Get(name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}

		roleBinding := &rbac.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			RoleRef:  roleBindingRef,
			Subjects: roleBindingSubjects,
		}

		if _, err := client.RbacV1beta1().ClusterRoleBindings().Create(roleBinding); err != nil {
			return err
		}
		fmt.Fprintln(out, "Successfully created cluster role bindings.")
	} else {
		roleBinding.RoleRef = roleBindingRef
		roleBinding.Subjects = roleBindingSubjects
		if _, err := client.RbacV1beta1().ClusterRoleBindings().Update(roleBinding); err != nil {
			return err
		}
		fmt.Fprintln(out, "Successfully updated cluster role bindings.")
	}

	return nil
}
