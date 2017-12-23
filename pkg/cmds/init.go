package cmds

import (
	"fmt"
	"io"
	"strings"

	"github.com/appscode/go/types"
	"github.com/kubedb/apimachinery/pkg/docker"
	"github.com/kubedb/cli/pkg/kube"
	"github.com/kubedb/cli/pkg/util"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	initLong = templates.LongDesc(`
		Install or upgrade KubeDB operator.`)

	initExample = templates.Examples(`
		# Install latest released operator.
		kubedb init

		# Upgrade operator to use another version.
		kubedb init --version=0.8.0 --upgrade`)
)

func NewCmdInit(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Create or upgrade KubeDB operator",
		Long:    initLong,
		Example: initExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunInit(cmd, out, errOut))
		},
	}

	util.AddInitFlags(cmd)
	return cmd
}

func RunInit(cmd *cobra.Command, out, errOut io.Writer) error {
	upgrade := cmdutil.GetFlagBool(cmd, "upgrade")

	if upgrade {
		return updateOperatorDeployment(cmd, out, errOut)
	} else {
		return createOperatorDeployment(cmd, out, errOut)
	}
	return nil
}

var operatorLabel = map[string]string{
	"app": "kubedb",
}

func createOperatorDeployment(cmd *cobra.Command, out, errOut io.Writer) error {

	namespace := cmdutil.GetFlagString(cmd, "operator-namespace")
	version := cmdutil.GetFlagString(cmd, "version")
	configureRBAC := cmdutil.GetFlagBool(cmd, "rbac")
	governingService := cmdutil.GetFlagString(cmd, "governing-service")
	exporterTag := cmdutil.GetFlagString(cmd, "exporter-tag")
	address := cmdutil.GetFlagString(cmd, "address")
	esdumpTag := cmdutil.GetFlagString(cmd, "elasticdump.tag")

	client, err := kube.NewKubeClient(cmd)
	if err != nil {
		return err
	}

	if err := docker.CheckDockerImageVersion(docker.ImageOperator, version); err != nil {
		fmt.Fprintln(errOut, fmt.Sprintf(`Operator image %v:%v not found.`, docker.ImageOperator, version))
		return nil
	}

	if configureRBAC {
		if err := EnsureRBACStuff(client, namespace, out); err != nil {
			return err
		}
	}

	deployment := &extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.OperatorName,
			Namespace: namespace,
			Labels:    operatorLabel,
		},
		Spec: extensions.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: operatorLabel,
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  docker.OperatorContainer,
							Image: fmt.Sprintf("%v:%v", docker.ImageOperator, version),
							Args: []string{
								"run",
								fmt.Sprintf("--governing-service=%v", governingService),
								fmt.Sprintf("--exporter-tag=%v", exporterTag),
								fmt.Sprintf("--address=%v", address),
								fmt.Sprintf("--rbac=%v", configureRBAC),
								fmt.Sprintf("--elasticdump.tag=%v", esdumpTag),
								"--v=3",
							},
							Env: []core.EnvVar{
								{
									Name: "OPERATOR_NAMESPACE",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
							},
							Ports: []core.ContainerPort{
								{
									Name:          docker.OperatorPortName,
									Protocol:      core.ProtocolTCP,
									ContainerPort: docker.OperatorPortNumber,
								},
							},
						},
					},
				},
			},
		},
	}

	if configureRBAC {
		deployment.Spec.Template.Spec.ServiceAccountName = ServiceAccountName
	}

	if _, err := client.ExtensionsV1beta1().Deployments(namespace).Create(deployment); err != nil {
		if kerr.IsAlreadyExists(err) {
			fmt.Fprintln(errOut, "Operator deployment already exists.")
		} else {
			return err
		}
	} else {
		fmt.Fprintln(out, "Successfully created operator deployment.")
	}

	if err := createOperatorService(client, namespace); err != nil {
		if kerr.IsAlreadyExists(err) {
			fmt.Fprintln(errOut, "Operator service already exists.")
		} else {
			return err
		}
	} else {
		fmt.Fprintln(out, "Successfully created operator service.")
	}

	return nil
}

func updateOperatorDeployment(cmd *cobra.Command, out, errOut io.Writer) error {

	namespace := cmdutil.GetFlagString(cmd, "operator-namespace")
	version := cmdutil.GetFlagString(cmd, "version")
	configureRBAC := cmdutil.GetFlagBool(cmd, "rbac")
	governingService := cmdutil.GetFlagString(cmd, "governing-service")
	exporterTag := cmdutil.GetFlagString(cmd, "exporter-tag")
	address := cmdutil.GetFlagString(cmd, "address")
	esdumpTag := cmdutil.GetFlagString(cmd, "elasticdump.tag")

	client, err := kube.NewKubeClient(cmd)
	if err != nil {
		return err
	}

	deployment, err := client.ExtensionsV1beta1().Deployments(namespace).Get(docker.OperatorName, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			message := fmt.Sprintf("Operator deployment \"%v\" not found.\n\n"+
				"Create operator using following commnad:\n"+
				"kubedb init --version=%v --operator-namespace=%v", docker.OperatorName, version, namespace)

			fmt.Fprintln(errOut, message)
			return nil
		}

		return err
	}

	containers := deployment.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		fmt.Fprintln(errOut, fmt.Sprintf(`Invalid operator deployment "%v"`, docker.OperatorName))
		return nil
	}

	items := strings.Split(containers[0].Image, ":")

	image := items[0]
	tag := items[1]

	if image != docker.ImageOperator {
		fmt.Fprintln(errOut, fmt.Sprintf(`Operator image mismatch. Can't upgrade to version "%v"`, version))
		return nil
	}

	if tag == version {
		fmt.Fprintln(out, "Operator deployment is already using this version.")
		return nil
	}

	if err := docker.CheckDockerImageVersion(docker.ImageOperator, version); err != nil {
		fmt.Fprintln(errOut, fmt.Sprintf(`Operator image %v:%v not found.`, docker.ImageOperator, version))
		return nil
	}

	deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%v:%v", docker.ImageOperator, version)

	if configureRBAC {
		if err := EnsureRBACStuff(client, namespace, out); err != nil {
			return err
		}
		deployment.Spec.Template.Spec.ServiceAccountName = ServiceAccountName
	} else {
		deployment.Spec.Template.Spec.ServiceAccountName = ""
	}

	deployment.Spec.Template.Spec.Containers[0].Args = []string{
		"run",
		fmt.Sprintf("--governing-service=%v", governingService),
		fmt.Sprintf("--exporter-tag=%v", exporterTag),
		fmt.Sprintf("--address=%v", address),
		fmt.Sprintf("--rbac=%v", configureRBAC),
		fmt.Sprintf("--elasticdump.tag=%v", esdumpTag),
		"--v=3",
	}

	if _, err := client.ExtensionsV1beta1().Deployments(deployment.Namespace).Update(deployment); err != nil {
		return err
	}

	fmt.Fprintln(out, "Successfully upgraded operator deployment.")

	return nil
}

func createOperatorService(client kubernetes.Interface, namespace string) error {
	svc := &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.OperatorName,
			Namespace: namespace,
			Labels:    operatorLabel,
		},
		Spec: core.ServiceSpec{
			Type: core.ServiceTypeClusterIP,
			Ports: []core.ServicePort{
				{
					Name:       docker.OperatorPortName,
					Port:       docker.OperatorPortNumber,
					Protocol:   core.ProtocolTCP,
					TargetPort: intstr.FromString(docker.OperatorPortName),
				},
			},
			Selector: operatorLabel,
		},
	}

	_, err := client.CoreV1().Services(namespace).Create(svc)
	return err
}
