package cmds

import (
	"errors"
	"fmt"
	"io"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	"github.com/kubedb/cli/pkg/kube"
	"github.com/kubedb/cli/pkg/util"
	"github.com/kubedb/cli/pkg/validator"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

// ref: k8s.io/kubernetes/pkg/kubectl/cmd/create.go

var (
	createLong = templates.LongDesc(`
		Create a resource by filename or stdin.

		JSON and YAML formats are accepted.`)

	createExample = templates.Examples(`
		# Create a elasticsearch using the data in elastic.json.
		kubedb create -f ./elastic.json

		# Create a elasticsearch based on the JSON passed into stdin.
		cat elastic.json | kubedb create -f -`)
)

func NewCmdCreate(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &resource.FilenameOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a resource by filename or stdin",
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			if cmdutil.IsFilenameSliceEmpty(options.Filenames) {
				defaultRunFunc := cmdutil.DefaultSubCommandRun(errOut)
				defaultRunFunc(cmd, args)
				return
			}
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunCreate(f, cmd, out, options))
		},
	}

	util.AddCreateFlags(cmd, options)
	return cmd
}

func RunCreate(f cmdutil.Factory, cmd *cobra.Command, out io.Writer, options *resource.FilenameOptions) error {
	cmdNamespace, enforceNamespace := util.GetNamespace(cmd)

	r := f.NewBuilder().Unstructured().Schema(util.Validator()).
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, options).
		Flatten().
		Do()

	err := r.Err()
	if err != nil {
		return err
	}

	config, err := f.ClientConfig()
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	extClient, err := cs.NewForConfig(config)
	if err != nil {
		return err
	}

	infoList := make([]*resource.Info, 0)
	err = r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		kind := info.Object.GetObjectKind().GroupVersionKind().Kind
		if err := util.CheckSupportedResource(kind); err != nil {
			return err
		}

		if kind == api.ResourceKindDormantDatabase {
			return fmt.Errorf(`resource type "%v" doesn't support create operation`, kind)
		}

		fmt.Println(fmt.Sprintf(`validating "%v"`, info.Source))
		if err := validator.Validate(client, extClient, info); err != nil {
			return cmdutil.AddSourceToErr("validating", info.Source, err)
		}

		infoList = append(infoList, info)
		return nil
	})
	if err != nil {
		return err
	}

	count := 0
	for _, info := range infoList {
		if err := createAndRefresh(info); err != nil {
			return cmdutil.AddSourceToErr("creating", info.Source, err)
		}
		count++
		cmdutil.PrintSuccess(false, out, info.Object, false, "created")
	}

	if count == 0 {
		return errors.New("no objects passed to create")
	}
	return nil
}

func createAndRefresh(info *resource.Info) error {
	obj, err := resource.NewHelper(info.Client, info.Mapping).Create(info.Namespace, true, info.Object)
	if err != nil {
		return err
	}
	info.Refresh(obj, true)
	return nil
}
