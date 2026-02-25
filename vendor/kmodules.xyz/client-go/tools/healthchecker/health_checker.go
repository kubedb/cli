/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package healthchecker

import (
	"context"
	"sync"
	"time"

	kmapi "kmodules.xyz/client-go/api/v1"

	"k8s.io/klog/v2"
)

type HealthChecker struct {
	healthCheckerMap map[string]healthCheckerData
	mux              sync.Mutex
}

type healthCheckerData struct {
	cancel            context.CancelFunc
	ticker            *time.Ticker
	lastPeriodSeconds int32
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		healthCheckerMap: make(map[string]healthCheckerData),
		mux:              sync.Mutex{},
	}
}

// Start creates a health check go routine.
// Call this method after successful creation of all the replicas of a database.
func (hc *HealthChecker) Start(key string, healthCheckSpec kmapi.HealthCheckSpec, fn func(string, *HealthCard)) {
	if healthCheckSpec.PeriodSeconds == nil || healthCheckSpec.TimeoutSeconds == nil || healthCheckSpec.FailureThreshold == nil {
		klog.Errorf("spec.healthCheck values are nil, can't start or modify health check.")
		return
	}

	if *healthCheckSpec.PeriodSeconds <= 0 {
		klog.Errorf("spec.healthCheck.PeriodSeconds can't be less than 1, can't start or modify health check.")
		return
	}

	if !hc.keyExists(key) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		ticker := time.NewTicker(time.Duration(*healthCheckSpec.PeriodSeconds) * time.Second)
		healthCheckStore := newHealthCard(key, *healthCheckSpec.FailureThreshold)
		hc.set(key, healthCheckerData{
			cancel:            cancel,
			ticker:            ticker,
			lastPeriodSeconds: *healthCheckSpec.PeriodSeconds,
		})
		go func() {
			for {
				select {
				case <-ctx.Done():
					hc.delete(key)
					cancel()
					ticker.Stop()
					klog.Infoln("Health check stopped for key " + key)
					return
				case <-ticker.C:
					klog.V(5).Infoln("Health check running for key " + key)
					fn(key, healthCheckStore)
					klog.V(5).Infof("Debug client count = %d\n", healthCheckStore.GetClientCount())
				}
			}
		}()
	} else {
		data := hc.get(key)
		if data.lastPeriodSeconds != *healthCheckSpec.PeriodSeconds {
			data.ticker.Reset(time.Duration(*healthCheckSpec.PeriodSeconds) * time.Second)
			data.lastPeriodSeconds = *healthCheckSpec.PeriodSeconds
			hc.set(key, data)
		}
	}
}

// Stop stops a health check go routine.
// Call this method when the database is deleted or halted.
func (hc *HealthChecker) Stop(key string) {
	if hc.keyExists(key) {
		hc.get(key).cancel()
		hc.delete(key)
	}
}

func (hc *HealthChecker) keyExists(key string) bool {
	hc.mux.Lock()
	defer hc.mux.Unlock()
	_, ok := hc.healthCheckerMap[key]
	return ok
}

func (hc *HealthChecker) get(key string) healthCheckerData {
	hc.mux.Lock()
	defer hc.mux.Unlock()
	return hc.healthCheckerMap[key]
}

func (hc *HealthChecker) set(key string, data healthCheckerData) {
	hc.mux.Lock()
	defer hc.mux.Unlock()
	hc.healthCheckerMap[key] = data
}

func (hc *HealthChecker) delete(key string) {
	hc.mux.Lock()
	defer hc.mux.Unlock()
	delete(hc.healthCheckerMap, key)
}
