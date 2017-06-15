package printer

import (
	"fmt"

	"github.com/spf13/cobra"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/kubernetes/pkg/kubectl"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// ref: k8s.io/kubernetes/pkg/kubectl/resource_printer.go

func NewPrinter(cmd *cobra.Command) (kubectl.ResourcePrinter, error) {
	humanReadablePrinter := NewHumanReadablePrinter(PrintOptions{
		WithNamespace: cmdutil.GetFlagBool(cmd, "all-namespaces"),
		Wide:          cmdutil.GetWideFlag(cmd),
		ShowAll:       cmdutil.GetFlagBool(cmd, "show-all"),
		ShowLabels:    cmdutil.GetFlagBool(cmd, "show-labels"),
	})

	format := cmdutil.GetFlagString(cmd, "output")

	switch format {
	case "json":
		return &kubectl.JSONPrinter{}, nil
	case "yaml":
		return &kubectl.YAMLPrinter{}, nil
	case "name":
		return &kubectl.NamePrinter{
			Typer:   apiv1.Scheme,
			Decoder: apiv1.Codecs.UniversalDecoder(),
		}, nil
	case "wide":
		fallthrough
	case "":
		return humanReadablePrinter, nil
	default:
		return nil, fmt.Errorf("output format %q not recognized", format)
	}
}

type editPrinterOptions struct {
	Printer   kubectl.ResourcePrinter
	Ext       string
	AddHeader bool
}

func NewEditPrinter(cmd *cobra.Command) (*editPrinterOptions, error) {
	switch format := cmdutil.GetFlagString(cmd, "output"); format {
	case "json":
		return &editPrinterOptions{
			Printer:   &kubectl.JSONPrinter{},
			Ext:       ".json",
			AddHeader: true,
		}, nil
	// If flag -o is not specified, use yaml as default
	case "yaml", "":
		return &editPrinterOptions{
			Printer:   &kubectl.YAMLPrinter{},
			Ext:       ".yaml",
			AddHeader: true,
		}, nil
	default:
		return nil, cmdutil.UsageError(cmd, "The flag 'output' must be one of yaml|json")
	}
}
