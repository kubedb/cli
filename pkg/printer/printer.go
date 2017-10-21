package printer

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/api"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/printers"
)

// ref: k8s.io/kubernetes/pkg/kubectl/resource_printer.go

func NewPrinter(cmd *cobra.Command) (printers.ResourcePrinter, error) {
	humanReadablePrinter := NewHumanReadablePrinter(PrintOptions{
		WithNamespace: cmdutil.GetFlagBool(cmd, "all-namespaces"),
		Wide:          cmdutil.GetWideFlag(cmd),
		ShowAll:       cmdutil.GetFlagBool(cmd, "show-all"),
		ShowLabels:    cmdutil.GetFlagBool(cmd, "show-labels"),
	})

	format := cmdutil.GetFlagString(cmd, "output")

	switch format {
	case "json":
		return &printers.JSONPrinter{}, nil
	case "yaml":
		return &printers.YAMLPrinter{}, nil
	case "name":
		return &printers.NamePrinter{
			Typer:    scheme.Scheme,
			Decoders: []runtime.Decoder{api.Codecs.UniversalDecoder()},
			Mapper:   api.Registry.RESTMapper(api.Registry.EnabledVersions()...),
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
	Printer   printers.ResourcePrinter
	Ext       string
	AddHeader bool
}

func NewEditPrinter(cmd *cobra.Command) (*editPrinterOptions, error) {
	switch format := cmdutil.GetFlagString(cmd, "output"); format {
	case "json":
		return &editPrinterOptions{
			Printer:   &printers.JSONPrinter{},
			Ext:       ".json",
			AddHeader: true,
		}, nil
	// If flag -o is not specified, use yaml as default
	case "yaml", "":
		return &editPrinterOptions{
			Printer:   &printers.YAMLPrinter{},
			Ext:       ".yaml",
			AddHeader: true,
		}, nil
	default:
		return nil, cmdutil.UsageError(cmd, "The flag 'output' must be one of yaml|json")
	}
}
