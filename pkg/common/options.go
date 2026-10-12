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

package common

// OptionFunc tunes how a per-engine Opts value is built.
type OptionFunc func(*optionConfig)

type optionConfig struct {
	skipReadinessCheck bool
}

// SkipReadinessCheck drops the `status.phase == Ready` gate on the database CR.
// Commands that only read the CR to render manifests do not need a healthy
// database, and refusing to run against a degraded source is actively unhelpful
// for disaster-recovery workflows, where the source being unhealthy is the
// whole reason the command is being run.
func SkipReadinessCheck() OptionFunc {
	return func(cfg *optionConfig) {
		cfg.skipReadinessCheck = true
	}
}

func buildOptionConfig(opts []OptionFunc) optionConfig {
	var cfg optionConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
