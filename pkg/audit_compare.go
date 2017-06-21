package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/k8sdb/cli/pkg/audit/type"
	"github.com/k8sdb/cli/pkg/util"
	"github.com/spf13/cobra"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewCmdAuditCompare(out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit_compare",
		Short: "Compare audit report",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(compareReport(cmd, out, cmdErr, args))
		},
	}
	util.AddAuditCompareFlags(cmd)
	return cmd
}

func compareReport(cmd *cobra.Command, out, errOut io.Writer, args []string) error {
	if len(args) != 2 {
		fmt.Fprint(errOut, "You must provide two summary report to compare.")
		usageString := "Summary reports not provided."
		return cmdutil.UsageError(cmd, usageString)
	}

	reportFile1 := args[0]
	reportFile2 := args[1]

	var reportData1 *types.Summary
	if err := util.ReadFileAs(reportFile1, &reportData1); err != nil {
		return err
	}

	var reportData2 *types.Summary
	if err := util.ReadFileAs(reportFile2, &reportData2); err != nil {
		return err
	}

	if reportData1.Kind != reportData2.Kind {
		return errors.New("Unable to compare these two summary. Kind mismatch")
	}

	reportData1Byte, err := json.Marshal(reportData1)
	if err != nil {
		return err
	}

	reportData2Byte, err := json.Marshal(reportData2)
	if err != nil {
		return err
	}

	// Then, compare them
	differ := diff.New()
	d, err := differ.Compare(reportData1Byte, reportData2Byte)
	if err != nil {
		return err
	}

	var aJson map[string]interface{}
	if err := json.Unmarshal(reportData1Byte, &aJson); err != nil {
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

	fmt.Println(fmt.Sprintf(`Compare result has been stored to '%v'`, path))

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
