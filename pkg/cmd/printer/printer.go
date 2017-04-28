package printer

import (
	"fmt"

	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubectl"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

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
			Typer:   kapi.Scheme,
			Decoder: kapi.Codecs.UniversalDecoder(),
		}, nil
	case "wide":
		fallthrough
	case "":
		return humanReadablePrinter, nil
	default:
		return nil, fmt.Errorf("output format %q not recognized", format)
	}
}
