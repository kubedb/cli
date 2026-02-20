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

package restic

import (
	"math"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/utils/ptr"
)

// RunBackup takes backup, cleanup old snapshots, check repository integrity etc.
// It extracts valuable information from respective restic command it runs and return them for further use.
func (w *ResticWrapper) RunBackup(backupOption BackupOptions) ([]BackupOutput, error) {
	// Start clock to measure total session duration
	startTime := time.Now()
	var backupOutput []BackupOutput

	// Run backup
	hostStats, err := w.runBackup(backupOption)
	for _, hostStat := range hostStats {
		if err != nil {
			hostStat.Phase = HostBackupFailed
			hostStat.Error = err.Error()
		} else {
			hostStat.Phase = HostBackupSucceeded
			hostStat.Duration = time.Since(startTime).String()
			st := metav1.Time{Time: startTime}
			et := metav1.Now()
			hostStat.StartTime = &st
			hostStat.EndTime = &et
		}
		backupOutput = append(backupOutput, BackupOutput{
			Stats: []HostBackupStats{hostStat},
		})
	}

	return backupOutput, err
}

func (w *ResticWrapper) runBackup(backupOption BackupOptions) ([]HostBackupStats, error) {
	var hostStats []HostBackupStats

	// fmt.Println("shell: ",w)
	// Backup from stdin
	if len(backupOption.StdinPipeCommands) != 0 {
		out, err := w.backupFromStdin(backupOption)
		if err != nil {
			return hostStats, err
		}
		// Extract information from the output of backup command
		snapStats, err := extractBackupInfo(out, backupOption.StdinFileName)
		if err != nil {
			return hostStats, err
		}
		for _, snap := range snapStats {
			hostStat := HostBackupStats{
				Hostname:  backupOption.Host,
				Snapshots: []SnapshotStats{snap},
			}
			hostStats = append(hostStats, hostStat)
		}
		return hostStats, nil
	}

	// Backup all target paths
	hostStats = make([]HostBackupStats, len(w.Config.Backends))
	for idx := range hostStats {
		hostStats[idx] = HostBackupStats{
			Hostname: backupOption.Host,
		}
	}
	for _, path := range backupOption.BackupPaths {
		startTime := time.Now()
		params := backupParams{
			path:     path,
			host:     backupOption.Host,
			excludes: backupOption.Exclude,
			args:     backupOption.Args,
		}
		out, err := w.backup(params)
		if err != nil {
			return hostStats, err
		}
		// Extract information from the output of backup command
		stats, err := extractBackupInfo(out, path)
		if err != nil {
			return hostStats, err
		}
		for idx, stat := range stats {
			hostStats[idx] = upsertSnapshotStats(hostStats[idx], stat)
		}
		w.updateElapsedTimeout(startTime)
	}

	return hostStats, nil
}

func (w *ResticWrapper) updateElapsedTimeout(startTime time.Time) {
	if w.Config.Timeout != nil {
		w.Config.Lock()
		defer w.Config.Unlock()

		// Calculate the new timeout by subtracting the elapsed duration
		elapsed := time.Since(startTime)
		newTimeout := time.Duration(math.Max(float64(w.Config.Timeout.Duration-elapsed), float64(0)))
		w.Config.Timeout = &metav1.Duration{
			Duration: newTimeout,
		}
	}
}

// RunParallelBackup runs multiple backup in parallel.
// Host must be different for each backup.
func (w *ResticWrapper) RunParallelBackup(backupOptions []BackupOptions, maxConcurrency int) ([]BackupOutput, error) {
	// WaitGroup to wait until all go routine finishes
	wg := sync.WaitGroup{}
	// concurrencyLimiter channel is used to limit maximum number simultaneous go routine
	concurrencyLimiter := make(chan bool, maxConcurrency)
	defer close(concurrencyLimiter)

	var (
		backupErrs []error
		mu         sync.Mutex // use lock to avoid racing condition
	)

	multipleOutputs := make([]BackupOutput, len(w.Config.Backends))

	for i := range backupOptions {
		// try to send message in concurrencyLimiter channel.
		// if maximum allowed concurrent backup is already running, program control will get stuck here.
		concurrencyLimiter <- true

		// starting new go routine. add it to WaitGroup
		wg.Add(1)

		go func(opt BackupOptions, startTime time.Time) {
			// when this go routine completes it task, release a slot from the concurrencyLimiter channel
			// so that another go routine can start. Also, tell the WaitGroup that it is done with its task.
			defer func() {
				<-concurrencyLimiter
				wg.Done()
			}()

			// sh field in ResticWrapper is a pointer. we must not use same w in multiple go routine.
			// otherwise they might enter in racing condition.
			nw := w.Copy()
			hostStats, err := nw.runBackup(opt)
			for idx, hostStat := range hostStats {
				if err != nil {
					hostStat.Phase = HostBackupFailed
					hostStat.Error = err.Error()
					mu.Lock()
					backupErrs = append(backupErrs, err)
					mu.Unlock()
				} else {
					hostStat.Phase = HostBackupSucceeded
					hostStat.Duration = time.Since(startTime).String()
					st := metav1.Time{Time: startTime}
					et := metav1.Now()
					hostStat.StartTime = &st
					hostStat.EndTime = &et
				}
				mu.Lock()
				multipleOutputs[idx].upsertHostBackupStats(hostStat)
				// add hostStats to backupOutput. use lock to avoid racing condition.
				mu.Unlock()
			}
		}(backupOptions[i], time.Now())
	}

	// wait for all the go routines to complete
	wg.Wait()

	return multipleOutputs, errors.NewAggregate(backupErrs)
}

func upsertSnapshotStats(hostStats HostBackupStats, snapStats SnapshotStats) HostBackupStats {
	for i, s := range hostStats.Snapshots {
		// if there is already an entry for this snapshot, then update it
		if s.Name == snapStats.Name {
			hostStats.Snapshots[i] = snapStats
			return hostStats
		}
	}
	// no entry for this snapshot. add a new entry
	hostStats.Snapshots = append(hostStats.Snapshots, snapStats)
	return hostStats
}

func (backupOutput *BackupOutput) upsertHostBackupStats(hostStats HostBackupStats) {
	// check if an entry already exist for this host in backupOutput. If exist then update it.
	for i, v := range backupOutput.Stats {
		if v.Hostname == hostStats.Hostname {
			backupOutput.Stats[i] = hostStats
			return
		}
	}

	// no entry for this host. add a new entry
	backupOutput.Stats = append(backupOutput.Stats, hostStats)
}

func (w *ResticWrapper) RepositoryAlreadyExist(repository string) bool {
	return w.repositoryExist(repository)
}

func (w *ResticWrapper) InitializeRepository(repository string) error {
	return w.initRepository(repository)
}

func (w *ResticWrapper) VerifyRepositoryIntegrity(repository string) (*RepositoryStats, error) {
	// Check repository integrity
	out, err := w.check(repository)
	if err != nil {
		return nil, err
	}
	// Extract information from output of "check" command
	integrity := extractCheckInfo(out)
	// Read repository statics after cleanup
	out, err = w.stats(repository, "")
	if err != nil {
		return nil, err
	}
	// Extract information from output of "stats" command
	repoSize, err := extractStatsInfo(out)
	if err != nil {
		return nil, err
	}
	return &RepositoryStats{Integrity: ptr.To(integrity), Size: repoSize}, nil
}
