package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/golang/glog"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/cli/pkg/cmd/editor"
	"github.com/k8sdb/cli/pkg/cmd/encoder"
	"github.com/k8sdb/cli/pkg/cmd/printer"
	"github.com/k8sdb/cli/pkg/cmd/util"
	"github.com/k8sdb/cli/pkg/kube"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
	"k8s.io/kubernetes/pkg/util/yaml"
)

var (
	editLong = templates.LongDesc(`
		Edit a resource from the default editor.

		The edit command allows you to directly edit any API resource you can retrieve via the
		command line tools. It will open the editor defined by your KUBEDB_EDITOR, or EDITOR
		environment variables, or fall back to 'nano'`)

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

func RunEdit(f cmdutil.Factory, out, errOut io.Writer, cmd *cobra.Command, args []string) error {
	return runEdit(f, out, errOut, cmd, args)
}

func runEdit(f cmdutil.Factory, out, errOut io.Writer, cmd *cobra.Command, args []string) error {
	o, err := printer.NewEditPrinter(cmd)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		usageString := "Required resource not specified."
		return cmdutil.UsageError(cmd, usageString)
	}

	resources := strings.Split(args[0], ",")
	for i, r := range resources {
		items := strings.Split(r, "/")
		kind, err := util.GetSupportedResource(items[0])
		if err != nil {
			return err
		}

		if kind == tapi.ResourceKindSnapshot {
			return fmt.Errorf(`resource type "%v" doesn't support edit operation`, items[0])
		}

		items[0] = kind
		resources[i] = strings.Join(items, "/")
	}
	args[0] = strings.Join(resources, ",")

	mapper, resourceMapper, r, _, err := getMapperAndResult(f, cmd, args)
	if err != nil {
		return err
	}

	normalEditInfos, err := r.Infos()
	if err != nil {
		return err
	}

	var (
		edit = editor.NewDefaultEditor()
	)

	editFn := func(info *resource.Info, err error) error {
		var (
			results  = editResults{}
			original = []byte{}
			edited   = []byte{}
			file     string
		)

		containsError := false
		infos := normalEditInfos
		for {
			originalObj := infos[0].Object
			objToEdit := originalObj

			buf := &bytes.Buffer{}
			var w io.Writer = buf

			if o.AddHeader {
				results.header.writeTo(w)
			}

			if !containsError {
				if err := o.Printer.PrintObj(objToEdit, w); err != nil {
					return preservedFile(err, results.file, errOut)
				}
				original = buf.Bytes()
			} else {
				buf.Write(manualStrip(edited))
			}

			// launch the editor
			editedDiff := edited
			edited, file, err = edit.LaunchTempFile(fmt.Sprintf("%s-edit-", filepath.Base(os.Args[0])), o.Ext, buf)
			if err != nil {
				return preservedFile(err, results.file, errOut)
			}
			if containsError {
				if bytes.Equal(stripComments(editedDiff), stripComments(edited)) {
					return preservedFile(fmt.Errorf("%s", "Edit cancelled, no valid changes were saved."), file, errOut)
				}
			}

			// cleanup any file from the previous pass
			if len(results.file) > 0 {
				os.Remove(results.file)
			}
			glog.V(4).Infof("User edited:\n%s", string(edited))

			// Compare content without comments
			if bytes.Equal(stripComments(original), stripComments(edited)) {
				os.Remove(file)
				fmt.Fprintln(errOut, "Edit cancelled, no changes made.")
				return nil
			}

			results = editResults{
				file: file,
			}

			// parse the edited file
			updates, err := resourceMapper.InfoForData(stripComments(edited), "edited-file")
			if err != nil {
				// syntax error
				containsError = true
				results.header.reasons = append(results.header.reasons, editReason{head: fmt.Sprintf("The edited file had a syntax error: %v", err)})
				continue
			}

			containsError = false

			err = visitToPatch(f, originalObj, updates, mapper, resourceMapper, out, unversioned.GroupVersion{}, &results)
			if err != nil {
				return preservedFile(err, results.file, errOut)
			}

			if results.notfound > 0 {
				fmt.Fprintf(errOut, "The edits you made on deleted resources have been saved to %q\n", file)
				return cmdutil.ErrExit
			}

			if len(results.edit) == 0 {
				if results.notfound == 0 {
					os.Remove(file)
				} else {
					fmt.Fprintf(out, "The edits you made on deleted resources have been saved to %q\n", file)
				}
				return nil
			}

			if len(results.header.reasons) > 0 {
				containsError = true
			}
		}
	}

	return editFn(nil, nil)
}

func visitToPatch(
	f cmdutil.Factory,
	originalObj runtime.Object,
	updates *resource.Info,
	mapper meta.RESTMapper,
	resourceMapper *resource.Mapper,
	out io.Writer,
	defaultVersion unversioned.GroupVersion,
	results *editResults,
) error {
	client, err := f.ClientSet()
	if err != nil {
		return err
	}

	restClonfig, err := f.ClientConfig()
	if err != nil {
		return err
	}

	extClient := clientset.NewExtensionsForConfigOrDie(restClonfig)

	patchVisitor := resource.NewFlattenListVisitor(updates, resourceMapper)
	err = patchVisitor.Visit(func(info *resource.Info, incomingErr error) error {

		currOriginalObj, err := util.GetStructuredObject(originalObj)
		if err != nil {
			return err
		}

		originalSerialization, err := runtime.Encode(clientset.ExtendedCodec, currOriginalObj)
		if err != nil {
			return err
		}

		editedSerialization, err := encoder.Encode(info.Object)
		if err != nil {
			return err
		}

		editedSerialization = stripComments(editedSerialization)

		originalJS, err := yaml.ToJSON(originalSerialization)
		if err != nil {
			return err
		}
		editedJS, err := yaml.ToJSON(editedSerialization)
		if err != nil {
			return err
		}

		if reflect.DeepEqual(originalJS, editedJS) {
			// no edit, so just skip it.
			cmdutil.PrintSuccess(mapper, false, out, info.Mapping.Resource, info.Name, false, "skipped")
			return nil
		}

		kind := currOriginalObj.GetObjectKind().GroupVersionKind().Kind
		preconditions := util.GetPreconditionFunc(kind)
		patch, err := strategicpatch.CreateTwoWayMergePatch(originalJS, editedJS, currOriginalObj, preconditions...)
		if err != nil {
			if strategicpatch.IsPreconditionFailed(err) {
				return preconditionFailedError()
			}
			return err
		}

		resourceExists, err := util.CheckResourceExists(client, kind, info.Name, info.Namespace)
		if err != nil {
			return err
		}
		if resourceExists {
			conditionalPreconditions := util.GetConditionalPreconditionFunc(kind)
			err := util.CheckConditionalPrecondition(patch, conditionalPreconditions...)
			if err != nil {
				if util.IsPreconditionFailed(err) {
					return conditionalPreconditionFailedError(kind)
				}
				return err
			}
		}

		results.version = defaultVersion
		h := resource.NewHelper(extClient.RESTClient(), info.Mapping)
		patched, err := extClient.RESTClient().Patch(kapi.MergePatchType).
			NamespaceIfScoped(info.Namespace, h.NamespaceScoped).
			Resource(h.Resource).
			Name(info.Name).
			Body(patch).
			Do().
			Get()

		if err != nil {
			fmt.Fprintln(out, results.addError(err, info))
			return nil
		}

		info.Refresh(patched, true)
		cmdutil.PrintSuccess(mapper, false, out, info.Mapping.Resource, info.Name, false, "edited")
		return nil
	})
	return err
}

func getMapperAndResult(f cmdutil.Factory, cmd *cobra.Command, args []string) (meta.RESTMapper, *resource.Mapper, *resource.Result, string, error) {
	cmdNamespace, enforceNamespace := util.GetNamespace(cmd)
	var mapper meta.RESTMapper
	var typer runtime.ObjectTyper
	mapper, typer, err := f.UnstructuredObject()
	if err != nil {
		return nil, nil, nil, "", err
	}

	resourceMapper := &resource.Mapper{
		ObjectTyper:  typer,
		RESTMapper:   mapper,
		ClientMapper: resource.ClientMapperFunc(f.UnstructuredClientForMapping),
		Decoder:      runtime.UnstructuredJSONScheme,
	}

	b := resource.NewBuilder(mapper, typer, resource.ClientMapperFunc(f.UnstructuredClientForMapping), runtime.UnstructuredJSONScheme).
		ResourceTypeOrNameArgs(false, args...).
		RequireObject(true).
		Latest()

	r := b.NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, &resource.FilenameOptions{}).
		ContinueOnError().
		Flatten().
		Do()

	err = r.Err()
	if err != nil {
		return nil, nil, nil, "", err
	}
	return mapper, resourceMapper, r, cmdNamespace, err
}

type editReason struct {
	head  string
	other []string
}

type editHeader struct {
	reasons []editReason
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

func (h *editHeader) flush() {
	h.reasons = []editReason{}
}

type editPrinterOptions struct {
	printer   kubectl.ResourcePrinter
	ext       string
	addHeader bool
}

// editResults capture the result of an update
type editResults struct {
	header   editHeader
	notfound int
	edit     []*resource.Info
	file     string

	version unversioned.GroupVersion
}

func (r *editResults) addError(err error, info *resource.Info) string {
	switch {
	case k8serr.IsInvalid(err):
		r.edit = append(r.edit, info)
		reason := editReason{
			head: fmt.Sprintf("%s %q was not valid", info.Mapping.Resource, info.Name),
		}
		if err, ok := err.(k8serr.APIStatus); ok {
			if details := err.Status().Details; details != nil {
				for _, cause := range details.Causes {
					reason.other = append(reason.other, fmt.Sprintf("%s: %s", cause.Field, cause.Message))
				}
			}
		}
		r.header.reasons = append(r.header.reasons, reason)
		return fmt.Sprintf("error: %s %q is invalid", info.Mapping.Resource, info.Name)
	case k8serr.IsNotFound(err):
		r.notfound++
		return fmt.Sprintf("error: %s %q could not be found on the server", info.Mapping.Resource, info.Name)
	default:
		return fmt.Sprintf("error: %s %q could not be patched: %v", info.Mapping.Resource, info.Name, err)
	}
}

func preservedFile(err error, path string, out io.Writer) error {
	if len(path) > 0 {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			fmt.Fprintf(out, "A copy of your changes has been stored to %q\n", path)
		}
	}
	return err
}

func stripComments(file []byte) []byte {
	stripped := file
	stripped, err := yaml.ToJSON(stripped)
	if err != nil {
		stripped = manualStrip(file)
	}
	return stripped
}

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

func preconditionFailedError() error {
	return errors.New(`At least one of the following was changed:
	apiVersion
	kind
	name
	namespace
	status`)
}

func conditionalPreconditionFailedError(kind string) error {
	str := util.PreconditionSpecField[kind]
	strList := strings.Join(str, "\n\t")
	return fmt.Errorf(`At least one of the following was changed:
	%v`, strList)
}
