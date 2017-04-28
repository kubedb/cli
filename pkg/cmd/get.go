package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/k8sdb/kubedb/pkg/cmd/printer"
	"github.com/k8sdb/kubedb/pkg/cmd/util"
	"github.com/k8sdb/kubedb/pkg/kube"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/kubectl"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
	utilerrors "k8s.io/kubernetes/pkg/util/errors"
)

type GetOptions struct {
	resource.FilenameOptions

	IgnoreNotFound bool
	Raw            string
}

func NewCmdGet(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {
			client := kube.GetKubeCmd(cmd)
			cmdutil.CheckErr(util.CheckSupportedResources(args))
			cmdutil.CheckErr(RunGet(client, out, errOut, cmd, args))
		},
	}

	util.AddContextFlag(cmd)
	util.AddGetFlags(cmd)
	return cmd
}

const (
	valid_resources = `Valid resource types include:

    * all
    * elastic
    * postgres
    * databasesnapshot
    * deleteddatabase
    `
)

func RunGet(client *kube.Client, out, errOut io.Writer, cmd *cobra.Command, args []string) error {

	allNamespaces := cmdutil.GetFlagBool(cmd, "all-namespaces")

	mapper, typer, err := client.UnstructuredObject()
	if err != nil {
		return err
	}
	cmdNamespace, _, err := client.DefaultNamespace()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		fmt.Fprint(errOut, "You must specify the type of resource to get. ", valid_resources)
		usageString := "Required resource not specified."
		return cmdutil.UsageError(cmd, usageString)
	}

	var printAll bool = false
	resourceList := args[0]
	resources := strings.Split(resourceList, ",")
	for _, a := range resources {
		if a == "all" {
			printAll = true
			supported, err := util.GetAllSupportedResources(client)
			if err != nil {
				return err
			}
			newArgs := make([]string, 0)
			newArgs = append(newArgs, supported)
			args = append(newArgs, args[1:]...)
			break
		}
	}

	argsHasNames, err := resource.HasNames(args)
	if err != nil {
		return err
	}
	if argsHasNames {
		cmd.Flag("show-all").Value.Set("true")
	}

	r := resource.NewBuilder(
		mapper,
		typer,
		resource.ClientMapperFunc(client.UnstructuredClientForMapping),
		runtime.UnstructuredJSONScheme,
	).NamespaceParam(cmdNamespace).DefaultNamespace().AllNamespaces(allNamespaces).
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
	if cmdutil.MustPrintWithKinds(objs, infos, nil, printAll) {
		showKind = true
	}

	var lastMapping *meta.RESTMapping

	w := kubectl.GetNewTabWriter(out)
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
