package roles

import (
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/docker"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	rbac "k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

var policyRuleOperator = []rbac.PolicyRule{
	{
		APIGroups: []string{extensions.GroupName},
		Resources: []string{"thirdpartyresources"},
		Verbs:     []string{"get", "create"},
	},
	{
		APIGroups: []string{tapi.GroupName},
		Resources: []string{rbac.ResourceAll},
		Verbs:     []string{"get", "list", "watch", "create", "update", "delete"},
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
		APIGroups: []string{batch.GroupName},
		Resources: []string{"jobs"},
		Verbs:     []string{"get", "create", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"events"},
		Verbs:     []string{"create"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"pods"},
		Verbs:     []string{"get", "list", "delete"},
	},
	{
		APIGroups: []string{apiv1.GroupName},
		Resources: []string{"persistentvolumeclaims"},
		Verbs:     []string{"list", "delete"},
	},
	{
		APIGroups: []string{"monitoring.coreos.com"},
		Resources: []string{"servicemonitors"},
		Verbs:     []string{"get", "create", "update"},
	},
}

func EnsureRBACStuff(client kubernetes.Interface, namespace string) error {
	operatorName := docker.OperatorName

	if _, err := client.RbacV1beta1().ClusterRoles().Create(operatorName, metav1.GetOptions{}); err != nil {
		return err
	}


	role := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: docker.OperatorName,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{extensions.GroupName},
				Resources: []string{"thirdpartyresources"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{tapi.GroupName},
				Resources: []string{rbac.ResourceAll},
				Verbs:     []string{"get", "list", "watch", "create", "update", "delete"},
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
				APIGroups: []string{batch.GroupName},
				Resources: []string{"jobs"},
				Verbs:     []string{"get", "create", "delete"},
			},
			{
				APIGroups: []string{apiv1.GroupName},
				Resources: []string{"events"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{apiv1.GroupName},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "delete"},
			},
			{
				APIGroups: []string{apiv1.GroupName},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"list", "delete"},
			},
			{
				APIGroups: []string{"monitoring.coreos.com"},
				Resources: []string{"servicemonitors"},
				Verbs:     []string{"get", "create", "update"},
			},
		},
	}

	if _, err := client.RbacV1beta1().ClusterRoles().Create(role); err != nil {
		return err
	}

	sa := &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccount,
			Namespace: namespace,
		},
	}
	if _, err := client.CoreV1().ServiceAccounts(namespace).Create(sa); err != nil {
		return err
	}

	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.OperatorName,
			Namespace: namespace,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     docker.OperatorName,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      serviceAccount,
				Namespace: namespace,
			},
		},
	}
	if _, err := client.RbacV1beta1().ClusterRoleBindings().Create(roleBinding); err != nil {
		return err
	}

	return nil
}
