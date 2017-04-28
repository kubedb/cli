package printer

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/client/clientset"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
)

const (
	tabwriterMinWidth = 10
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = 0
)

type handlerEntry struct {
	printFunc reflect.Value
	args      []reflect.Value
}

type PrintOptions struct {
	WithNamespace bool
	WithKind      bool
	Wide          bool
	ShowAll       bool
	ShowLabels    bool
	Kind          string
}

type HumanReadablePrinter struct {
	handlerMap   map[reflect.Type]*handlerEntry
	options      PrintOptions
	lastType     reflect.Type
	hiddenObjNum int
}

func NewHumanReadablePrinter(options PrintOptions) *HumanReadablePrinter {
	printer := &HumanReadablePrinter{
		handlerMap: make(map[reflect.Type]*handlerEntry),
		options:    options,
	}
	printer.addDefaultHandlers()
	return printer
}

func shortHumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < -1 {
		return fmt.Sprintf("<invalid>")
	} else if seconds < 0 {
		return fmt.Sprintf("0s")
	} else if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if minutes := int(d.Minutes()); minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	} else if hours := int(d.Hours()); hours < 24 {
		return fmt.Sprintf("%dh", hours)
	} else if hours < 24*364 {
		return fmt.Sprintf("%dd", hours/24)
	}
	return fmt.Sprintf("%dy", int(d.Hours()/24/365))
}

func (h *HumanReadablePrinter) addDefaultHandlers() {
	h.Handler(h.printElasticList)
	h.Handler(h.printElastic)
	h.Handler(h.printPostgresList)
	h.Handler(h.printPostgres)
}

func (h *HumanReadablePrinter) Handler(printFunc interface{}) error {
	printFuncValue := reflect.ValueOf(printFunc)
	if err := h.validatePrintHandlerFunc(printFuncValue); err != nil {
		glog.Errorf("Unable to add print handler: %v", err)
		return err
	}

	objType := printFuncValue.Type().In(0)

	h.handlerMap[objType] = &handlerEntry{
		printFunc: printFuncValue,
	}
	return nil
}

func (h *HumanReadablePrinter) validatePrintHandlerFunc(printFunc reflect.Value) error {
	if printFunc.Kind() != reflect.Func {
		return fmt.Errorf("invalid print handler. %#v is not a function", printFunc)
	}
	funcType := printFunc.Type()
	if funcType.NumIn() != 3 || funcType.NumOut() != 1 {
		return fmt.Errorf("invalid print handler." +
			"Must accept 3 parameters and return 1 value.")
	}
	if funcType.In(1) != reflect.TypeOf((*io.Writer)(nil)).Elem() ||
		funcType.In(2) != reflect.TypeOf((*PrintOptions)(nil)).Elem() ||
		funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("invalid print handler. The expected signature is: "+
			"func handler(obj %v, w io.Writer, options PrintOptions) error", funcType.In(0))
	}
	return nil
}

func getColumns(options PrintOptions, t reflect.Type) []string {
	columns := make([]string, 0)

	if options.WithNamespace {
		columns = append(columns, "NAMESPACE")
	}

	columns = append(columns, "NAME", "STATUS")

	columns = append(columns, formatWideHeaders(options.Wide, t)...)

	columns = append(columns, "AGE")

	if options.ShowLabels {
		columns = append(columns, "LABELS")
	}

	return columns
}

func (h *HumanReadablePrinter) printElastic(item *tapi.Elastic, w io.Writer, options PrintOptions) error {
	name := formatResourceName(options.Kind, item.Name, options.WithKind)

	namespace := item.Namespace

	if options.WithNamespace {
		if _, err := fmt.Fprintf(w, "%s\t", namespace); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "%s\t%s\t", name, item.Status.DatabaseStatus); err != nil {
		return err
	}

	if options.Wide {
		if _, err := fmt.Fprintf(w, "%s\t", item.Spec.Version); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "%s", translateTimestamp(item.CreationTimestamp)); err != nil {
		return err
	}

	_, err := fmt.Fprint(w, appendAllLabels(options.ShowLabels, item.Labels))

	return err
}

func (h *HumanReadablePrinter) printElasticList(podList *tapi.ElasticList, w io.Writer, options PrintOptions) error {
	for _, pod := range podList.Items {
		if err := h.printElastic(pod, w, options); err != nil {
			return err
		}
	}
	return nil
}

func (h *HumanReadablePrinter) printPostgres(item *tapi.Postgres, w io.Writer, options PrintOptions) error {
	name := formatResourceName(options.Kind, item.Name, options.WithKind)

	namespace := item.Namespace

	if options.WithNamespace {
		if _, err := fmt.Fprintf(w, "%s\t", namespace); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "%s\t%s\t", name, item.Status.DatabaseStatus); err != nil {
		return err
	}

	if options.Wide {
		if _, err := fmt.Fprintf(w, "%s\t", item.Spec.Version); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "%s", translateTimestamp(item.CreationTimestamp)); err != nil {
		return err
	}

	_, err := fmt.Fprint(w, appendAllLabels(options.ShowLabels, item.Labels))

	return err
}

func (h *HumanReadablePrinter) printPostgresList(podList *tapi.PostgresList, w io.Writer, options PrintOptions) error {
	for _, pod := range podList.Items {
		if err := h.printPostgres(pod, w, options); err != nil {
			return err
		}
	}
	return nil
}

func decode(kind string, data []byte) (runtime.Object, error) {
	switch kind {
	case tapi.ResourceKindElastic:
		var elastic *tapi.Elastic
		if err := yaml.Unmarshal(data, &elastic); err != nil {
			return nil, err
		}
		return elastic, nil
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(data, &postgres); err != nil {
			return nil, err
		}
		return postgres, nil
	}

	return nil, fmt.Errorf(`Invalid kind: "%v"`, kind)
}

func (h *HumanReadablePrinter) PrintObj(obj runtime.Object, output io.Writer) error {
	w, found := output.(*tabwriter.Writer)
	if !found {
		w = GetNewTabWriter(output)
		defer w.Flush()
	}

	kind := obj.GetObjectKind().GroupVersionKind().Kind

	switch obj.(type) {
	case *runtime.UnstructuredList, *runtime.Unstructured, *runtime.Unknown:
		if objBytes, err := runtime.Encode(clientset.ExtendedCodec, obj); err == nil {
			if decodedObj, err := decode(kind, objBytes); err == nil {
				obj = decodedObj
			}
		}
	}

	t := reflect.TypeOf(obj)
	if handler := h.handlerMap[t]; handler != nil {
		if t != h.lastType {
			headers := getColumns(h.options, t)
			if h.lastType != nil {
				printNewline(w)
			}
			h.printHeader(headers, w)
			h.lastType = t
		}
		args := []reflect.Value{reflect.ValueOf(obj), reflect.ValueOf(w), reflect.ValueOf(h.options)}
		resultValue := handler.printFunc.Call(args)[0]
		if resultValue.IsNil() {
			return nil
		}
		return resultValue.Interface().(error)
	}

	return fmt.Errorf(`kubedb doesn't support: "%v"`, kind)
}

func (h *HumanReadablePrinter) HandledResources() []string {
	return []string{}
}

func (h *HumanReadablePrinter) AfterPrint(io.Writer, string) error {
	return nil
}

func (h *HumanReadablePrinter) printHeader(columnNames []string, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s\n", strings.Join(columnNames, "\t")); err != nil {
		return err
	}
	return nil
}

func printNewline(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
	return nil
}

func formatResourceName(kind, name string, withKind bool) string {
	if !withKind || kind == "" {
		return name
	}

	return kind + "/" + name
}

func (h *HumanReadablePrinter) GetResourceKind() string {
	return h.options.Kind
}

func (h *HumanReadablePrinter) EnsurePrintWithKind(kind string) {
	h.options.WithKind = true
	h.options.Kind = kind
}

func appendAllLabels(showLabels bool, itemLabels map[string]string) string {
	var buffer bytes.Buffer

	if showLabels {
		buffer.WriteString(fmt.Sprint("\t"))
		buffer.WriteString(labels.FormatLabels(itemLabels))
	}
	buffer.WriteString("\n")

	return buffer.String()
}

func translateTimestamp(timestamp unversioned.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}
	return shortHumanDuration(time.Now().Sub(timestamp.Time))
}

func GetNewTabWriter(output io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(
		output,
		tabwriterMinWidth,
		tabwriterWidth,
		tabwriterPadding,
		tabwriterPadChar,
		tabwriterFlags,
	)
}

// headers for -o wide
func formatWideHeaders(wide bool, t reflect.Type) []string {
	if wide {
		if t.String() == "*api.Elastic" || t.String() == "*api.ElasticList" {
			return []string{"VERSION"}
		}
		if t.String() == "*api.Postgres" || t.String() == "*api.PostgresList" {
			return []string{"VERSION"}
		}
	}
	return nil
}
