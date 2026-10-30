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
	"fmt"
	"io"
)

const (
	levelError = "ERROR"
	levelWarn  = "WARN"
	levelInfo  = "OK"
)

type finding struct {
	Level string
	Msg   string
}

// findings are the results of the pre-flight checks. An ERROR blocks the
// adoption; a WARN needs a look before the handover.
type findings []finding

func (f *findings) errorf(format string, a ...any) {
	*f = append(*f, finding{levelError, fmt.Sprintf(format, a...)})
}

func (f *findings) warnf(format string, a ...any) {
	*f = append(*f, finding{levelWarn, fmt.Sprintf(format, a...)})
}

func (f *findings) infof(format string, a ...any) {
	*f = append(*f, finding{levelInfo, fmt.Sprintf(format, a...)})
}

func (f findings) hasErrors() bool {
	for _, x := range f {
		if x.Level == levelError {
			return true
		}
	}
	return false
}

func (f findings) print(w io.Writer) {
	for _, x := range f {
		_, _ = fmt.Fprintf(w, "  [%-5s] %s\n", x.Level, x.Msg)
	}
}
