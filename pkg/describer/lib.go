package describer

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tabbedString(f func(io.Writer) error) (string, error) {
	out := new(tabwriter.Writer)
	buf := &bytes.Buffer{}
	out.Init(buf, 0, 8, 1, '\t', 0)

	err := f(out)
	if err != nil {
		return "", err
	}

	out.Flush()
	str := string(buf.String())
	return str, nil
}

func printLabelsMultiline(out io.Writer, title string, labels map[string]string) {
	printLabelsMultilineWithIndent(out, "", title, "\t", labels)
}

func printLabelsMultilineWithIndent(out io.Writer, initialIndent, title, innerIndent string, labels map[string]string) {

	fmt.Fprintf(out, "%s%s:%s", initialIndent, title, innerIndent)

	if labels == nil || len(labels) == 0 {
		fmt.Fprintln(out, "<none>")
		return
	}

	// to print labels in the sorted order
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		if i != 0 {
			fmt.Fprint(out, initialIndent)
			fmt.Fprint(out, innerIndent)
		}
		value := labels[key]
		if value != "" {
			fmt.Fprintf(out, "%s=%s\n", key, labels[key])
		} else {
			fmt.Fprintf(out, "%s\n", key)
		}
		i++
	}
}

func timeToString(t *metav1.Time) string {
	if t == nil {
		return ""
	}

	return t.Format(time.RFC1123Z)
}

func getAccessModesAsString(modes []core.PersistentVolumeAccessMode) string {
	modes = removeDuplicateAccessModes(modes)
	modesStr := []string{}
	if containsAccessMode(modes, core.ReadWriteOnce) {
		modesStr = append(modesStr, "RWO")
	}
	if containsAccessMode(modes, core.ReadOnlyMany) {
		modesStr = append(modesStr, "ROX")
	}
	if containsAccessMode(modes, core.ReadWriteMany) {
		modesStr = append(modesStr, "RWX")
	}
	return strings.Join(modesStr, ",")
}

func removeDuplicateAccessModes(modes []core.PersistentVolumeAccessMode) []core.PersistentVolumeAccessMode {
	accessModes := []core.PersistentVolumeAccessMode{}
	for _, m := range modes {
		if !containsAccessMode(accessModes, m) {
			accessModes = append(accessModes, m)
		}
	}
	return accessModes
}

func containsAccessMode(modes []core.PersistentVolumeAccessMode, mode core.PersistentVolumeAccessMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}
