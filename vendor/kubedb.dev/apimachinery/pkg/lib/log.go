/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"io"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initZapLogger(w io.Writer, vLevel zap.AtomicLevel) *zap.Logger {
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(enc, zapcore.AddSync(w), vLevel)
	return zap.New(core, zap.AddCaller())
}

func NewLogger(verbosity int) logr.Logger {
	vLevel := convertVerbosityToZapLevel(verbosity)
	return zapr.NewLogger(initZapLogger(os.Stdout, vLevel))
}

func convertVerbosityToZapLevel(verbosity int) zap.AtomicLevel {
	if verbosity < 0 {
		verbosity = 0
	}
	if verbosity >= 4 {
		verbosity = 4
	}

	var lvl zapcore.Level
	switch verbosity {
	case 0:
		lvl = zapcore.DPanicLevel
	case 1:
		lvl = zapcore.ErrorLevel
	case 2:
		lvl = zapcore.WarnLevel
	case 4:
		lvl = zapcore.DebugLevel
	default:
		lvl = zapcore.InfoLevel
	}
	return zap.NewAtomicLevelAt(lvl)
}
