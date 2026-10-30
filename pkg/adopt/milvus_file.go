/*
Copyright AppsCode Inc. and Contributors
Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package adopt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	adoptionFileAPIVersion = "cli.kubedb.com/v1alpha1"
	adoptionFileKind       = "MilvusAdoption"
	// adoptionFileName is the file "plan --output-dir" writes.
	adoptionFileName = "adopt.yaml"
)

// adoptionFile holds the options of an adoption, so they need not be typed as
// flags on every command. A flag given on the command line wins over the file.
type adoptionFile struct {
	APIVersion    string           `json:"apiVersion"`
	Kind          string           `json:"kind"`
	Release       string           `json:"release"`
	Namespace     string           `json:"namespace,omitempty"`
	Name          string           `json:"name,omitempty"`
	Version       string           `json:"version,omitempty"`
	Source        adoptionSource   `json:"source,omitempty"`
	Auth          adoptionAuth     `json:"auth,omitempty"`
	RestrictedSCC *bool            `json:"restrictedSCC,omitempty"`
	Handover      adoptionHandover `json:"handover,omitempty"`
	Rollback      adoptionRollback `json:"rollback,omitempty"`
}

type adoptionSource struct {
	Selector   string             `json:"selector,omitempty"`
	RoleLabel  string             `json:"roleLabel,omitempty"`
	Workloads  []adoptionWorkload `json:"workloads,omitempty"`
	ConfigFile string             `json:"configFile,omitempty"`
}

type adoptionWorkload struct {
	Role  string `json:"role"`
	Group string `json:"group,omitempty"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
}

type adoptionAuth struct {
	AuthSecret          string `json:"authSecret,omitempty"`
	RootPasswordFile    string `json:"rootPasswordFile,omitempty"`
	EtcdAuthSecret      string `json:"etcdAuthSecret,omitempty"`
	SkipCredentialCheck *bool  `json:"skipCredentialCheck,omitempty"`
}

type adoptionHandover struct {
	Method          string           `json:"method,omitempty"`
	QueryNodeDelay  *metav1.Duration `json:"querynodeDelay,omitempty"`
	Settle          *metav1.Duration `json:"settle,omitempty"`
	Timeout         *metav1.Duration `json:"timeout,omitempty"`
	ReadyTimeout    *metav1.Duration `json:"readyTimeout,omitempty"`
	SkipConfigCheck *bool            `json:"skipConfigCheck,omitempty"`
}

type adoptionRollback struct {
	Cleanup *bool `json:"cleanup,omitempty"`
}

// loadAdoptionFile reads and checks an adoption file. Unknown keys are errors,
// so a misspelt option is not silently ignored.
func loadAdoptionFile(path string) (*adoptionFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	af := &adoptionFile{}
	if err := yaml.UnmarshalStrict(raw, af); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if err := af.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return af, nil
}

func (af *adoptionFile) validate() error {
	if af.APIVersion != adoptionFileAPIVersion || af.Kind != adoptionFileKind {
		return fmt.Errorf("want apiVersion %s and kind %s, got %q and %q", adoptionFileAPIVersion, adoptionFileKind, af.APIVersion, af.Kind)
	}
	if af.Release == "" {
		return errors.New("release is required")
	}
	if m := af.Handover.Method; m != "" && m != methodOverlap && m != methodStopStart {
		return fmt.Errorf("handover.method must be %s or %s, got %q", methodOverlap, methodStopStart, m)
	}
	for i, w := range af.Source.Workloads {
		if !contains(milvusRoles, w.Role) {
			return fmt.Errorf("source.workloads[%d]: role %q is not one of %v", i, w.Role, milvusRoles)
		}
		if w.Kind != "Deployment" && w.Kind != "StatefulSet" {
			return fmt.Errorf("source.workloads[%d]: kind must be Deployment or StatefulSet, got %q", i, w.Kind)
		}
		if w.Name == "" {
			return fmt.Errorf("source.workloads[%d]: name is required", i)
		}
	}
	return nil
}

// workloadFlag renders a workload the way --workload takes it.
func (w adoptionWorkload) workloadFlag() string {
	role := w.Role
	if w.Group != "" {
		role += "/" + w.Group
	}
	return role + "=" + w.Kind + "/" + w.Name
}

type flagValue struct {
	flag   string
	values []string
}

// flagValues maps the file onto the command's flags.
func (af *adoptionFile) flagValues() []flagValue {
	var out []flagValue
	str := func(flag, v string) {
		if v != "" {
			out = append(out, flagValue{flag, []string{v}})
		}
	}
	boolean := func(flag string, v *bool) {
		if v != nil {
			out = append(out, flagValue{flag, []string{strconv.FormatBool(*v)}})
		}
	}
	duration := func(flag string, v *metav1.Duration) {
		if v != nil {
			out = append(out, flagValue{flag, []string{v.Duration.String()}})
		}
	}
	str("name", af.Name)
	str("version", af.Version)
	str("selector", af.Source.Selector)
	str("role-label", af.Source.RoleLabel)
	if len(af.Source.Workloads) > 0 {
		fv := flagValue{flag: "workload"}
		for _, w := range af.Source.Workloads {
			fv.values = append(fv.values, w.workloadFlag())
		}
		out = append(out, fv)
	}
	str("config-file", af.Source.ConfigFile)
	str("auth-secret", af.Auth.AuthSecret)
	str("root-password-file", af.Auth.RootPasswordFile)
	str("etcd-auth-secret", af.Auth.EtcdAuthSecret)
	boolean("skip-credential-check", af.Auth.SkipCredentialCheck)
	boolean("restricted-scc", af.RestrictedSCC)
	str("method", af.Handover.Method)
	duration("querynode-delay", af.Handover.QueryNodeDelay)
	duration("settle", af.Handover.Settle)
	duration("timeout", af.Handover.Timeout)
	duration("ready-timeout", af.Handover.ReadyTimeout)
	boolean("skip-config-check", af.Handover.SkipConfigCheck)
	boolean("cleanup", af.Rollback.Cleanup)
	return out
}

// applyFlags sets every flag of fs that the file gives and the command line
// did not. Options a command does not have are skipped.
func (af *adoptionFile) applyFlags(fs *pflag.FlagSet) error {
	for _, fv := range af.flagValues() {
		fl := fs.Lookup(fv.flag)
		if fl == nil || fl.Changed {
			continue
		}
		for _, v := range fv.values {
			if err := fs.Set(fv.flag, v); err != nil {
				return fmt.Errorf("--%s from the file: %w", fv.flag, err)
			}
		}
	}
	return nil
}

// fileOptions is the -f flag every adopt milvus command has.
type fileOptions struct {
	Filename string
	file     *adoptionFile
}

func addFileFlag(cmd *cobra.Command, o *fileOptions) {
	cmd.Flags().StringVarP(&o.Filename, "filename", "f", "", "File with the adoption options (kind MilvusAdoption; \"plan --output-dir\" writes one). Flags given on the command line win")
}

// resolve loads the file, fills the flags it gives, and returns the release:
// the RELEASE argument, or release from the file.
func (o *fileOptions) resolve(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 1 {
		return "", errors.New("give at most one RELEASE")
	}
	if o.Filename != "" {
		af, err := loadAdoptionFile(o.Filename)
		if err != nil {
			return "", err
		}
		if err := af.applyFlags(cmd.Flags()); err != nil {
			return "", err
		}
		o.file = af
	}
	switch {
	case len(args) == 1 && o.file != nil && o.file.Release != args[0]:
		return "", fmt.Errorf("RELEASE %q differs from release %q in %s", args[0], o.file.Release, o.Filename)
	case len(args) == 1:
		return args[0], nil
	case o.file != nil:
		return o.file.Release, nil
	}
	return "", errors.New("give the RELEASE (the Helm release, or the instance label value of the workloads), or a file with -f")
}

// namespace is the file's namespace, used when -n is not given.
func (o *fileOptions) namespace() string {
	if o.file == nil {
		return ""
	}
	return o.file.Namespace
}

// adoptionFileFor writes down what a plan found, as a file to review and pass
// to run, status and rollback with -f.
func adoptionFileFor(p *adoptPlan, o planOptions, ns string) ([]byte, error) {
	version, _ := p.Objects.Milvus.Object["spec"].(map[string]any)["version"].(string)
	af := adoptionFile{
		APIVersion: adoptionFileAPIVersion,
		Kind:       adoptionFileKind,
		Release:    o.Release,
		Namespace:  ns,
		Name:       p.State.Name,
		Version:    version,
		Source:     adoptionSource{ConfigFile: o.ConfigFile},
		Auth: adoptionAuth{
			AuthSecret:       o.AuthSecret,
			RootPasswordFile: o.RootPasswordFile,
			EtcdAuthSecret:   o.EtcdAuthSecret,
		},
		Handover: adoptionHandover{
			Method:         methodOverlap,
			QueryNodeDelay: &metav1.Duration{Duration: time.Minute},
			Settle:         &metav1.Duration{Duration: 20 * time.Second},
			Timeout:        &metav1.Duration{Duration: 15 * time.Minute},
			ReadyTimeout:   &metav1.Duration{Duration: 15 * time.Minute},
		},
		Rollback: adoptionRollback{Cleanup: ptrBool(false)},
	}
	if o.RestrictedSCC {
		af.RestrictedSCC = ptrBool(true)
	}
	for _, u := range p.State.Units {
		af.Source.Workloads = append(af.Source.Workloads, adoptionWorkload{Role: u.Role, Group: u.Group, Kind: u.Kind, Name: u.Name})
	}
	b, err := yaml.Marshal(af)
	if err != nil {
		return nil, err
	}
	header := "# Options of the adoption, as plan found them. Review, then:\n" +
		"#   kubectl dba adopt milvus run -f " + adoptionFileName + " --yes\n" +
		"# Raise handover.querynodeDelay for large collections (about the time one QueryNode takes to load).\n"
	return append([]byte(header), b...), nil
}

func ptrBool(b bool) *bool { return &b }
