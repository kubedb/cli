package printer

// ref: k8s.io/kubernetes/pkg/kubectl/resource_printer.go

// DescriberSettings holds display configuration for each object
// describer to control what is printed.
type DescriberSettings struct {
	ShowEvents   bool
	ShowWorkload bool
	ShowSecret   bool
}
