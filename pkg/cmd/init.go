package cmd

import (
	"fmt"
	"io"
	"strings"

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
)

var (
	init_long = templates.LongDesc(`
		Create or upgrade unified operator for kubedb databases.`)

	init_example = templates.Examples(`
		# Create operator with version canary.
		kubedb init --version=canary

		# Upgrade operator to use another version.
		kubedb init --version=canary --upgrade`)
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
	operatorName      = "kubedb-operator"
	operatorImage     = "kubedb/operator"
	operatorContainer = "operator"
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

		if image != operatorImage {
			fmt.Fprintln(errOut, fmt.Sprintf(`Operator image mismatch. Can't upgrade to version "%v"`, version))
			return nil
		}

		if tag == version {
			fmt.Fprintln(out, "Operator deployment is already using this version.")
			return nil
		}

		if err := util.CheckDockerImageVersion(operatorImage, version); err != nil {
			fmt.Fprintln(errOut, fmt.Sprintf(`Operator image %v:%v not found.`, operatorImage, version))
			return nil
		}

		deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%v:%v", operatorImage, version)

		if err := updateOperatorDeployment(client, deployment); err != nil {
			return err
		}

		fmt.Fprintln(out, "Successfully upgraded operator deployment.")
	} else {
		var err error
		if _, err = getOperatorDeployment(client, namespace); err == nil {
			fmt.Fprintln(errOut, "Operator already exists.")
			return nil
		} else {
			if !k8serr.IsNotFound(err) {
				return err
			}
		}

		if err := util.CheckDockerImageVersion(operatorImage, version); err != nil {
			fmt.Fprintln(errOut, fmt.Sprintf(`Operator image %v:%v not found.`, operatorImage, version))
			return nil
		}

		if err := createOperatorDeployment(client, namespace, version); err != nil {
			return err
		}

		fmt.Fprintln(out, "Successfully created operator deployment.")
	}

	return nil
}

func getOperatorDeployment(client *internalclientset.Clientset, namespace string) (*kext.Deployment, error) {
	return client.ExtensionsClient.Deployments(namespace).Get(operatorName)
}

func createOperatorDeployment(client *internalclientset.Clientset, namespace, version string) error {
	label := map[string]string{
		"run": operatorName,
	}

	deployment := &kext.Deployment{
		ObjectMeta: kapi.ObjectMeta{
			Name:      operatorName,
			Namespace: namespace,
		},
		Spec: kext.DeploymentSpec{
			Selector: &unversioned.LabelSelector{
				MatchLabels: label,
			},
			Replicas: 1,
			Template: kapi.PodTemplateSpec{
				ObjectMeta: kapi.ObjectMeta{
					Labels: label,
				},
				Spec: kapi.PodSpec{
					Containers: []kapi.Container{
						{
							Name:  operatorContainer,
							Image: fmt.Sprintf("%v:%v", operatorImage, version),
							Args: []string{
								"run",
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
						},
					},
				},
			},
		},
	}

	_, err := client.ExtensionsClient.Deployments(namespace).Create(deployment)
	return err
}

func updateOperatorDeployment(client *internalclientset.Clientset, deployment *kext.Deployment) error {
	_, err := client.ExtensionsClient.Deployments(deployment.Namespace).Update(deployment)
	return err
}
