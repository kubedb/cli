package cmds

import (
	"fmt"
	"io"
	"strings"

	"github.com/k8sdb/cli/pkg/kube"
	"github.com/k8sdb/cli/pkg/printer"
	"github.com/k8sdb/cli/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/printers"
)

// ref: k8s.io/kubernetes/pkg/kubectl/cmd/get.go

var (
	get_long = templates.LongDesc(`
		Display one or many resources.

		` + valid_resources)

	get_example = templates.Examples(`
		# List all elastic in ps output format.
		kubedb get elastics

		# List all elastic in ps output format with more information (such as version).
		kubedb get elastics -o wide

		# List a single postgres with specified NAME in ps output format.
		kubedb get postgres database

		# List a single snapshot in JSON output format.
		kubedb get -o json snapshot snapshot-xyz

		# List all postgreses and elastics together in ps output format.
		kubedb get postgreses,elastics

		# List one or more resources by their type and names.
		kubedb get elastic/es-db postgres/pg-db`)
)

func NewCmdGet(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Display one or many resources",
		Long:    get_long,
		Example: get_example,
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunGet(f, cmd, out, errOut, args))
		},
	}

	util.AddGetFlags(cmd)
	return cmd
}

const (
	valid_resources = `Valid resource types include:

    * all
    * elastic
    * postgres
    * mysql
    * mongodb
    * redis
    * memcached
    * snapshot
    * dormantdatabase
    `
)

func RunGet(f cmdutil.Factory, cmd *cobra.Command, out, errOut io.Writer, args []string) error {
	selector := cmdutil.GetFlagString(cmd, "selector")
	cmdNamespace, enforceNamespace := util.GetNamespace(cmd)
	allNamespaces := cmdutil.GetFlagBool(cmd, "all-namespaces")

	categoryExpander := f.CategoryExpander()
	mapper, typer, err := f.UnstructuredObject()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		fmt.Fprint(errOut, "You must specify the type of resource to get. ", valid_resources)
		usageString := "Required resource not specified."
		return cmdutil.UsageErrorf(cmd, usageString)
	}

	var printAll bool = false
	resources := strings.Split(args[0], ",")
	for i, r := range resources {
		if r == "all" {
			printAll = true
		} else {
			items := strings.Split(r, "/")
			kind, err := util.GetSupportedResource(items[0])
			if err != nil {
				return err
			}
			items[0] = kind
			resources[i] = strings.Join(items, "/")
		}
	}

	if printAll {
		if resources, err = util.GetAllSupportedResources(f); err != nil {
			return err
		}
	}

	args[0] = strings.Join(resources, ",")

	argsHasNames, err := resource.HasNames(args)
	if err != nil {
		return err
	}
	if argsHasNames {
		cmd.Flag("show-all").Value.Set("true")
	}

	r := resource.NewBuilder(mapper, categoryExpander, typer, resource.ClientMapperFunc(f.UnstructuredClientForMapping), unstructured.UnstructuredJSONScheme).
		NamespaceParam(cmdNamespace).DefaultNamespace().AllNamespaces(allNamespaces).
		FilenameParam(enforceNamespace, &resource.FilenameOptions{}).
		SelectorParam(selector).
		ResourceTypeOrNameArgs(true, args...).
		ContinueOnError().
		Latest().
		Flatten().
		Do()

	err = r.Err()
	if err != nil {
		return err
	}

	allErrs := []error{}
	infos, err := r.Infos()
	if err != nil {
		allErrs = append(allErrs, err)
	}
	if len(infos) == 0 && len(allErrs) == 0 {
		outputEmptyListWarning(errOut)
	}

	objs := make([]runtime.Object, len(infos))
	for ix := range infos {
		objs[ix] = infos[ix].Object
	}

	rPrinter, err := printer.NewPrinter(cmd)
	if err != nil {
		return err
	}

	showKind := cmdutil.GetFlagBool(cmd, "show-kind")
	if printAll {
		showKind = true
	} else {
		if cmdutil.MustPrintWithKinds(objs, infos, nil) {
			showKind = true
		}
	}

	var lastMapping *meta.RESTMapping

	w := printers.GetNewTabWriter(out)
	for ix := range objs {
		var mapping *meta.RESTMapping
		var original runtime.Object
		mapping = infos[ix].Mapping
		original = infos[ix].Object

		if resourcePrinter, found := rPrinter.(*printer.HumanReadablePrinter); found {
			if lastMapping == nil || mapping.Resource != lastMapping.Resource {
				lastMapping = mapping
			}
			var resourceName string
			if mapping != nil {
				resourceName = lastMapping.Resource

				if alias, ok := util.ResourceShortFormFor(mapping.Resource); ok {
					resourceName = alias
				} else if resourceName == "" {
					resourceName = "none"
				}
			} else {
				resourceName = "none"
			}

			if showKind {
				resourcePrinter.EnsurePrintWithKind(resourceName)
			}

			if err := rPrinter.PrintObj(original, w); err != nil {
				allErrs = append(allErrs, err)
			}
			continue
		}

		if err := rPrinter.PrintObj(original, w); err != nil {
			allErrs = append(allErrs, err)
			continue
		}
	}
	w.Flush()
	return utilerrors.NewAggregate(allErrs)
}

func outputEmptyListWarning(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s\n", "No resources found.")
	return err
}
