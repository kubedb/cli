package pkg

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	tapi "github.com/k8sdb/apimachinery/api"
	auditpg "github.com/k8sdb/cli/pkg/audit/postgres"
	"github.com/k8sdb/cli/pkg/audit/type"
	"github.com/k8sdb/cli/pkg/util"
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
	if len(args) == 0 {
		fmt.Fprint(errOut, "You must specify the type of resource. ", valid_resources_for_report)
		usageString := "Required resource not specified."
		return cmdutil.UsageError(cmd, usageString)
	}

	if len(strings.Split(args[0], ",")) > 1 {
		return errors.New("audit doesn't support multiple resource")
	}

	if len(args) > 1 {
		return errors.New("audit doesn't support resource name")
	}

	kubedbType := args[0]
	if len(strings.Split(kubedbType, "/")) > 1 {
		return errors.New("audit doesn't support resource name")
	}
	kubedbType, err := util.GetResourceType(kubedbType)
	if err != nil {
		return err
	}

	original := cmdutil.GetFlagString(cmd, "original")
	if original == "" {
		return errors.New("Summary report file of original database is not provided")
	}

	duplicate := cmdutil.GetFlagString(cmd, "duplicate")
	if duplicate == "" {
		return errors.New("Summary report file of duplicate database is not provided")
	}

	diff := make([]string, 0)

	switch kubedbType {
	case tapi.ResourceTypePostgres:
		var originalData *types.Summary
		if err := util.ReadFileAs(original, &originalData); err != nil {
			return err
		}
		serializedOriginalData := auditpg.ToSerializePostgresReport(originalData)

		var duplicateData *types.Summary
		if err := util.ReadFileAs(duplicate, &duplicateData); err != nil {
			return err
		}
		serializedDuplicateData := auditpg.ToSerializePostgresReport(duplicateData)

		diff = auditpg.DiffSerializePostgresReport(serializedOriginalData, serializedDuplicateData)
	}

	outputDirectory := cmdutil.GetFlagString(cmd, "output")
	fileName := fmt.Sprintf("result-%v.txt", time.Now().UTC().Format("20060102-150405"))
	path := filepath.Join(outputDirectory, fileName)

	diffLines := strings.Join(diff, "\n")
	if err := util.WriteJson(path, []byte(diffLines)); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf(`Compare result has been stored to '%v'`, path))
	return nil
}
