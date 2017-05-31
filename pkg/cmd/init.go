package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/k8sdb/apimachinery/pkg/docker"
	"github.com/k8sdb/kubedb/pkg/cmd/util"
	"github.com/k8sdb/kubedb/pkg/kube"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	kext "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/util/intstr"
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

const (
	operatorName       = "kubedb-operator"
	operatorContainer  = "operator"
	operatorPortName   = "web"
	operatorPortNumber = 8080
)

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
			if k8serr.IsNotFound(err) {
				message := fmt.Sprintf("Operator deployment \"%v\" not found.\n\n"+
					"Create operator using following commnad:\n"+
					"kubedb init --version=%v --namespace=%v", operatorName, version, namespace)

				fmt.Fprintln(errOut, message)
				return nil
			}

			return err
		}

		containers := deployment.Spec.Template.Spec.Containers
		if len(containers) == 0 {
			fmt.Fprintln(errOut, fmt.Sprintf(`Invalid operator deployment "%v"`, operatorName))
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
			if k8serr.IsAlreadyExists(err) {
				fmt.Fprintln(errOut, "Operator deployment already exists.")
			} else {
				return err
			}
		} else {
			fmt.Fprintln(out, "Successfully created operator deployment.")
		}

		if err := createOperatorService(client, namespace); err != nil {
			if k8serr.IsAlreadyExists(err) {
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

func getOperatorDeployment(client *internalclientset.Clientset, namespace string) (*kext.Deployment, error) {
	return client.ExtensionsClient.Deployments(namespace).Get(operatorName)
}

var operatorLabel = map[string]string{
	"app": operatorName,
}

func createOperatorDeployment(client *internalclientset.Clientset, namespace, version string) error {
	deployment := &kext.Deployment{
		ObjectMeta: kapi.ObjectMeta{
			Name:      operatorName,
			Namespace: namespace,
		},
		Spec: kext.DeploymentSpec{
			Selector: &unversioned.LabelSelector{
				MatchLabels: operatorLabel,
			},
			Replicas: 1,
			Template: kapi.PodTemplateSpec{
				ObjectMeta: kapi.ObjectMeta{
					Labels: operatorLabel,
				},
				Spec: kapi.PodSpec{
					Containers: []kapi.Container{
						{
							Name:  operatorContainer,
							Image: fmt.Sprintf("%v:%v", docker.ImageOperator, version),
							Args: []string{
								"run",
								fmt.Sprintf("--address=:%v", operatorPortNumber),
								"--v=3",
							},
							Env: []kapi.EnvVar{
								{
									Name: "OPERATOR_NAMESPACE",
									ValueFrom: &kapi.EnvVarSource{
										FieldRef: &kapi.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
							},
							Ports: []kapi.ContainerPort{
								{
									Name:          operatorPortName,
									Protocol:      kapi.ProtocolTCP,
									ContainerPort: operatorPortNumber,
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

func createOperatorService(client *internalclientset.Clientset, namespace string) error {
	svc := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      operatorName,
			Namespace: namespace,
		},
		Spec: kapi.ServiceSpec{
			Type: kapi.ServiceTypeClusterIP,
			Ports: []kapi.ServicePort{
				{
					Name:       operatorPortName,
					Port:       operatorPortNumber,
					Protocol:   kapi.ProtocolTCP,
					TargetPort: intstr.FromString(operatorPortName),
				},
			},
			Selector: operatorLabel,
		},
	}

	_, err := client.Core().Services(namespace).Create(svc)
	return err
}

func updateOperatorDeployment(client *internalclientset.Clientset, deployment *kext.Deployment) error {
	_, err := client.ExtensionsClient.Deployments(deployment.Namespace).Update(deployment)
	return err
}
