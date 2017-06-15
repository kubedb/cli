package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/k8sdb/apimachinery/pkg/docker"
	"github.com/k8sdb/cli/pkg/cmd/util"
	"github.com/k8sdb/cli/pkg/kube"
	"github.com/spf13/cobra"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	init_long = templates.LongDesc(`
		Create or upgrade unified operator for kubedb databases.`)

	init_example = templates.Examples(`
		# Create operator with version canary.
		kubedb init --version=0.1.0

		# Upgrade operator to use another version.
		kubedb init --version=0.1.0 --upgrade`)
)

func NewCmdInit(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Create or upgrade unified operator",
		Long:    init_long,
		Example: init_example,
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunInit(f, cmd, out, errOut))
		},
	}

	util.AddInitFlags(cmd)
	return cmd
}

func RunInit(f cmdutil.Factory, cmd *cobra.Command, out, errOut io.Writer) error {
	upgrade := cmdutil.GetFlagBool(cmd, "upgrade")
	namespace := cmdutil.GetFlagString(cmd, "namespace")
	version := cmdutil.GetFlagString(cmd, "version")

	f.RESTClient()
	client, err := f.ClientSet()
	if err != nil {
		return err
	}

	if upgrade {
		deployment, err := getOperatorDeployment(client, namespace)
		if err != nil {
			if kerr.IsNotFound(err) {
				message := fmt.Sprintf("Operator deployment \"%v\" not found.\n\n"+
					"Create operator using following commnad:\n"+
					"kubedb init --version=%v --namespace=%v", docker.OperatorName, version, namespace)

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

		if err := util.CheckDockerImageVersion(docker.ImageOperator, version); err != nil {
			fmt.Fprintln(errOut, fmt.Sprintf(`Operator image %v:%v not found.`, docker.ImageOperator, version))
			return nil
		}

		deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%v:%v", docker.ImageOperator, version)

		if err := updateOperatorDeployment(client, deployment); err != nil {
			return err
		}

		fmt.Fprintln(out, "Successfully upgraded operator deployment.")
	} else {
		if err := util.CheckDockerImageVersion(docker.ImageOperator, version); err != nil {
			fmt.Fprintln(errOut, fmt.Sprintf(`Operator image %v:%v not found.`, docker.ImageOperator, version))
			return nil
		}

		if err := createOperatorDeployment(client, namespace, version); err != nil {
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
	}

	return nil
}

func getOperatorDeployment(client *clientset.Clientset, namespace string) (*extensions.Deployment, error) {
	return client.ExtensionsClient.Deployments(namespace).Get(docker.OperatorName)
}

var operatorLabel = map[string]string{
	"app": docker.OperatorName,
}

func createOperatorDeployment(client *clientset.Clientset, namespace, version string) error {
	deployment := &extensions.Deployment{
		ObjectMeta: apiv1.ObjectMeta{
			Name:      docker.OperatorName,
			Namespace: namespace,
		},
		Spec: extensions.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: operatorLabel,
			},
			Replicas: 1,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: apiv1.ObjectMeta{
					Labels: operatorLabel,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  docker.OperatorContainer,
							Image: fmt.Sprintf("%v:%v", docker.ImageOperator, version),
							Args: []string{
								"run",
								fmt.Sprintf("--address=:%v", docker.OperatorPortNumber),
								"--v=3",
							},
							Env: []apiv1.EnvVar{
								{
									Name: "OPERATOR_NAMESPACE",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          docker.OperatorPortName,
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: docker.OperatorPortNumber,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.ExtensionsClient.Deployments(namespace).Create(deployment)
	return err
}

func createOperatorService(client *clientset.Clientset, namespace string) error {
	svc := &apiv1.Service{
		ObjectMeta: apiv1.ObjectMeta{
			Name:      docker.OperatorName,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Ports: []apiv1.ServicePort{
				{
					Name:       docker.OperatorPortName,
					Port:       docker.OperatorPortNumber,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromString(docker.OperatorPortName),
				},
			},
			Selector: operatorLabel,
		},
	}

	_, err := client.Core().Services(namespace).Create(svc)
	return err
}

func updateOperatorDeployment(client *clientset.Clientset, deployment *extensions.Deployment) error {
	_, err := client.ExtensionsClient.Deployments(deployment.Namespace).Update(deployment)
	return err
}
