package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	acio "github.com/appscode/go/io"
	"github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/kubedb/pkg/cmd/decoder"
	"github.com/k8sdb/kubedb/pkg/cmd/encoder"
	"github.com/k8sdb/kubedb/pkg/cmd/printer"
	"github.com/k8sdb/kubedb/pkg/cmd/util"
	"github.com/k8sdb/kubedb/pkg/kube"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/editor"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
	utilerrors "k8s.io/kubernetes/pkg/util/errors"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
	"k8s.io/kubernetes/pkg/util/yaml"
	"os"
	"path/filepath"
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

	util.AddEditFlags(cmd)
	return cmd
}

func RunEdit(f cmdutil.Factory, out, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	o, err := printer.NewEditPrinter(cmd)
	if err != nil {
		return err
	}

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

	var (
		edit = editor.NewDefaultEditor(f.EditorEnvs())
	)

	restClonfig, _ := f.ClientConfig()
	extClient := clientset.NewExtensionsForConfigOrDie(restClonfig)

	editFn := func(info *resource.Info) error {

		var (
			results = editResults{}
			//originalByte = []byte{}
			edited = []byte{}
			//file     string
		)

		containsError := false
		original, err := util.GetStructuredObject(info.Object)
		if err != nil {
			return err
		}
		for {

			buf := &bytes.Buffer{}
			var w io.Writer = buf
			var editedData, path string

			results.header.writeTo(w)

			if !containsError {
				if err := rPrinter.PrintObj(original, w); err != nil {
					return preservedFile(err, results.file, cmdErr)
				}
				/*var ok bool
				if ok, path = editor.WriteTempFile(info, buf); !ok {
					return errors.New("Fail to write temp file")
				}*/

			}

			/*if err := editor.OpenEditor(editorName, path); err != nil {
				return fmt.Errorf("Editor %s not working. Change editor by setting EDITOR environment variable", editorName)
			}*/

			editedDiff := edited
			edited, file, err := edit.LaunchTempFile(fmt.Sprintf("%s-edit-", filepath.Base(os.Args[0])), o.Ext, buf)
			if err != nil {
				return preservedFile(err, results.file, cmdErr)
			}

			if editMode == NormalEditMode || containsError {
				if bytes.Equal(stripComments(editedDiff), stripComments(edited)) {
					// Ugly hack right here. We will hit this either (1) when we try to
					// save the same changes we tried to save in the previous iteration
					// which means our changes are invalid or (2) when we exit the second
					// time. The second case is more usual so we can probably live with it.
					// TODO: A less hacky fix would be welcome :)
					return preservedFile(fmt.Errorf("%s", "Edit cancelled, no valid changes were saved."), file, errOut)
				}
			}

			if editedData, err = acio.ReadFile(path); err != nil {
				return err
			}

			editedRuntimeObject, err := decoder.Decode(info.GetObjectKind().GroupVersionKind().Kind, []byte(editedData))
			if err != nil {
				return err
			}

			originalSerialization, err := runtime.Encode(clientset.ExtendedCodec, original)
			if err != nil {
				return err
			}

			editedSerialization, err := encoder.Encode(editedRuntimeObject)
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

			fmt.Println("--+ ", string(originalJS))
			fmt.Println()
			fmt.Println("--+ ", string(editedJS))

			// Compare content without comments
			if bytes.Equal(stripComments(originalJS), stripComments(editedJS)) {
				os.Remove(path)
				fmt.Fprintln(cmdErr, "Edit cancelled, no changes made.")
				return nil
			}

			preconditions := []strategicpatch.PreconditionFunc{
				strategicpatch.RequireKeyUnchanged("apiVersion"),
				strategicpatch.RequireKeyUnchanged("kind"),
				strategicpatch.RequireMetadataKeyUnchanged("name"),
			}

			patch, err := strategicpatch.CreateTwoWayMergePatch(originalJS, editedJS, original, preconditions...)
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
				continue
			}

			info.Refresh(patched, true)
			cmdutil.PrintSuccess(mapper, false, out, info.Mapping.Resource, info.Name, false, "edited")
			return nil
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

func stripComments(file []byte) []byte {
	stripped := file
	stripped, err := yaml.ToJSON(stripped)
	if err != nil {
		stripped = manualStrip(file)
	}
	return stripped
}

// manualStrip is used for dropping comments from a YAML file
func manualStrip(file []byte) []byte {
	stripped := []byte{}
	lines := bytes.Split(file, []byte("\n"))
	for i, line := range lines {
		if bytes.HasPrefix(bytes.TrimSpace(line), []byte("#")) {
			continue
		}
		stripped = append(stripped, line...)
		if i < len(lines)-1 {
			stripped = append(stripped, '\n')
		}
	}
	return stripped
}

type editReason struct {
	head  string
	other []string
}

// editHeader includes a list of reasons the edit must be retried
type editHeader struct {
	reasons []editReason
}

// editResults capture the result of an update
type editResults struct {
	header    editHeader
	retryable int
	notfound  int
	edit      []*resource.Info
	file      string

	version unversioned.GroupVersion
}

// writeTo outputs the current header information into a stream
func (h *editHeader) writeTo(w io.Writer) error {
	fmt.Fprint(w, `# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
`)
	for _, r := range h.reasons {
		if len(r.other) > 0 {
			fmt.Fprintf(w, "# %s:\n", r.head)
		} else {
			fmt.Fprintf(w, "# %s\n", r.head)
		}
		for _, o := range r.other {
			fmt.Fprintf(w, "# * %s\n", o)
		}
		fmt.Fprintln(w, "#")
	}
	return nil
}

func preservedFile(err error, path string, out io.Writer) error {
	if len(path) > 0 {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			fmt.Fprintf(out, "A copy of your changes has been stored to %q\n", path)
		}
	}
	return err
}
