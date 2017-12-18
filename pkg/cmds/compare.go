package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	tapi "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/cli/pkg/util"
	"github.com/spf13/cobra"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewCmdCompare(out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare summary reports",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(compareReports(cmd, out, cmdErr, args))
		},
	}
	util.AddCompareFlags(cmd)
	return cmd
}

func compareReports(cmd *cobra.Command, out, errOut io.Writer, args []string) error {
	if len(args) != 2 {
		fmt.Fprint(errOut, "You must provide two summary report to compare.")
		usageString := "Summary reports not provided."
		return cmdutil.UsageErrorf(cmd, usageString)
	}

	reportFile1 := args[0]
	reportFile2 := args[1]

	var report1 *tapi.Report
	if err := util.ReadFileAs(reportFile1, &report1); err != nil {
		return err
	}

	var report2 *tapi.Report
	if err := util.ReadFileAs(reportFile2, &report2); err != nil {
		return err
	}

	if report1.Kind != report2.Kind {
		return fmt.Errorf("summary reports are not for same type database. (%s, %s)", report1.Kind, report2.Kind)
	}

	report1Bytes, err := json.Marshal(report1)
	if err != nil {
		return err
	}

	report2Bytes, err := json.Marshal(report2)
	if err != nil {
		return err
	}

	// Then, compare them
	differ := diff.New()
	d, err := differ.Compare(report1Bytes, report2Bytes)
	if err != nil {
		return err
	}

	var aJson map[string]interface{}
	if err := json.Unmarshal(report1Bytes, &aJson); err != nil {
		return err
	}

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       false,
	}

	format := formatter.NewAsciiFormatter(aJson, config)
	diffString, err := format.Format(d)
	if err != nil {
		return err
	}

	outputDirectory := cmdutil.GetFlagString(cmd, "output")
	fileName := fmt.Sprintf("result-%v.txt", time.Now().UTC().Format("20060102-150405"))
	path := filepath.Join(outputDirectory, fileName)

	if err := util.WriteJson(path, []byte(diffString)); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf(`Comparison result has been stored in '%v'.`, path))

	show := cmdutil.GetFlagBool(cmd, "show")
	if show {
		config := formatter.AsciiFormatterConfig{
			ShowArrayIndex: true,
			Coloring:       true,
		}
		format := formatter.NewAsciiFormatter(aJson, config)
		diffString, err := format.Format(d)
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println(string(diffString))
	}

	return nil
}
