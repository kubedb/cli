package pkg

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/cli/pkg/audit/type"
	"github.com/k8sdb/cli/pkg/util"
	pgaudit "github.com/k8sdb/postgres/pkg/audit/compare"
	"github.com/spf13/cobra"
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

	diff := make([]string, 0)

	switch reportData1.Kind {
	case tapi.ResourceKindPostgres:
		index := cmdutil.GetFlagString(cmd, "index")
		diff = pgaudit.CompareSummaryReport(
			reportData1.SummaryReport.Postgres,
			reportData2.SummaryReport.Postgres,
			index,
		)
	}

	result := make([]string, 0)
	if len(diff) == 0 {
		result = append(result, "No change\n")
	} else {
		result = append(result, fmt.Sprintf("Total mismatch: %d\n", len(diff)))
		result = append(result, addFileNames(reportFile1, reportFile2))
		result = append(result, fmt.Sprintf("Diff:"))
		result = append(result, diff...)
		result = append(result, "")
	}

	outputDirectory := cmdutil.GetFlagString(cmd, "output")
	fileName := fmt.Sprintf("result-%v.txt", time.Now().UTC().Format("20060102-150405"))
	path := filepath.Join(outputDirectory, fileName)

	resultLines := strings.Join(result, "\n")
	if err := util.WriteJson(path, []byte(resultLines)); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf(`Compare result has been stored to '%v'`, path))

	show := cmdutil.GetFlagBool(cmd, "show")
	if show {
		fmt.Println()
		fmt.Println(resultLines)
	}

	return nil
}

func addFileNames(reportFile1, reportFile2 string) string {
	return fmt.Sprintf("%v --> %v\n", reportFile1, reportFile2)
}
