package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"errors"
	acio "github.com/appscode/go/io"
	"github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/kubedb/pkg/cmd/decoder"
	"github.com/k8sdb/kubedb/pkg/cmd/editor"
	"github.com/k8sdb/kubedb/pkg/cmd/util"
	"github.com/k8sdb/kubedb/pkg/kube"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
	utilerrors "k8s.io/kubernetes/pkg/util/errors"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
	"k8s.io/kubernetes/pkg/util/yaml"
)

var (
	editLong = templates.LongDesc(`
		Edit a resource from the default editor.

		The edit command allows you to directly edit any API resource you can retrieve via the
		command line tools. It will open the editor defined by your KUBEDB_EDITOR, or EDITOR
		environment variables, or fall back to 'vi' for Linux.
		You can edit multiple objects, although changes are applied one at a time.`)

	editExample = templates.Examples(`
		# Edit the elastic named 'elasticsearch-demo':
		kubedb edit es/elasticsearch-demo

		# Use an alternative editor
		KUBEDB_EDITOR="nano" kubedb edit es/elasticsearch-demo`)
)

func NewCmdEdit(out, errOut io.Writer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "edit (RESOURCE/NAME)",
		Short:   "Edit a resource on the server",
		Long:    editLong,
		Example: fmt.Sprintf(editExample),
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunEdit(f, out, errOut, cmd, args))
		},
	}
	return cmd
}

func RunEdit(f cmdutil.Factory, out, cmdErr io.Writer, cmd *cobra.Command, args []string) error {

	cmdNamespace, _, err := f.DefaultNamespace()
	if err != nil {
		return err
	}

	mapper, typer, err := f.UnstructuredObject()
	if err != nil {
		return err
	}

	resources := strings.Split(args[0], ",")
	for i, r := range resources {
		items := strings.Split(r, "/")
		kind, err := util.GetSupportedResourceKind(items[0])
		if err != nil {
			return err
		}
		items[0] = kind
		resources[i] = strings.Join(items, "/")
	}
	args[0] = strings.Join(resources, ",")

	r := resource.NewBuilder(
		mapper,
		typer,
		resource.ClientMapperFunc(f.UnstructuredClientForMapping),
		runtime.UnstructuredJSONScheme).
		ResourceTypeOrNameArgs(true, args...).
		Latest().
		NamespaceParam(cmdNamespace).
		DefaultNamespace().
		ContinueOnError().
		Flatten().
		Do()

	infos, err := r.Infos()
	if err != nil {
		return err
	}

	rPrinter := &kubectl.YAMLPrinter{}
	if err != nil {
		return err
	}

	restClonfig, _ := f.ClientConfig()
	extClient := clientset.NewExtensionsForConfigOrDie(restClonfig)

	editorName := editor.GetEditor()

	editFn := func(info *resource.Info) error {
		for {

			buf := &bytes.Buffer{}
			var w io.Writer = buf

			if err := rPrinter.PrintObj(info.Object, w); err != nil {
				return err
			}

			ok, path := editor.WriteTempFile(info, buf)
			if !ok {
				return errors.New("Fail to write temp file")
			}

			if err := editor.OpenEditor(editorName, path); err != nil {
				return fmt.Errorf("Editor %s not working. Change editor by setting EDITOR environment variable", editorName)
			}

			editedData, err := acio.ReadFile(path)
			if err != nil {
				return err
			}

			editedRuntimeObject, err := decoder.Decode(info.GetObjectKind().GroupVersionKind().Kind, []byte(editedData))
			if err != nil {
				return err
			}

			originalSerialization, err := runtime.Encode(clientset.ExtendedCodec, info.Object)
			if err != nil {
				return err
			}

			editedSerialization, err := runtime.Encode(clientset.ExtendedCodec, editedRuntimeObject)
			if err != nil {
				return err
			}

			originalJS, err := yaml.ToJSON(originalSerialization)
			if err != nil {
				return err
			}
			editedJS, err := yaml.ToJSON(editedSerialization)
			if err != nil {
				return err
			}

			preconditions := []strategicpatch.PreconditionFunc{
				strategicpatch.RequireKeyUnchanged("apiVersion"),
				//strategicpatch.RequireKeyUnchanged("kind"),
				strategicpatch.RequireMetadataKeyUnchanged("name"),
			}

			patch, err := strategicpatch.CreateTwoWayMergePatch(originalJS, editedJS, editedRuntimeObject, preconditions...)
			if err != nil {
				if strategicpatch.IsPreconditionFailed(err) {
					return fmt.Errorf("%s", "At least one of apiVersion, kind and name was changed")
				}
				return err
			}

			h := resource.NewHelper(extClient.RESTClient(), info.Mapping)

			patched, err := extClient.RESTClient().Patch(kapi.MergePatchType).
				NamespaceIfScoped(info.Namespace, h.NamespaceScoped).
				Resource(h.Resource).
				Name(info.Name).
				Body(patch).
				Do().
				Get()

			if err != nil {

				return err
			}

			info.Refresh(patched, true)
			cmdutil.PrintSuccess(mapper, false, out, info.Mapping.Resource, info.Name, false, "edited")
		}
		return nil
	}

	allErrs := []error{}

	for _, info := range infos {
		err := editFn(info)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return utilerrors.NewAggregate(allErrs)
}
