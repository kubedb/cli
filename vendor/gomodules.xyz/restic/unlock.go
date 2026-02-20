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
	"bytes"
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (w *ResticWrapper) UnlockRepository(repository string) error {
	_, err := w.unlock(repository)
	return err
}

// getLockIDs lists every lock ID currently held in the repository.
func (w *ResticWrapper) getLockIDs(repository string) ([]string, error) {
	w.sh.ShowCMD = true
	out, err := w.listLocks(repository)
	if err != nil {
		return nil, err
	}
	return extractLockIDs(bytes.NewReader(out))
}

// getLockStats returns the decoded JSON for a single lock.
func (w *ResticWrapper) getLockStats(repository, lockID string) (*LockStats, error) {
	w.sh.ShowCMD = true
	out, err := w.lockStats(repository, lockID)
	if err != nil {
		return nil, err
	}
	return extractLockStats(out)
}

// hasExclusiveLock checks if any exclusive lock exists in the repository.
// This should be called AFTER unlockStale() - any remaining exclusive locks are active.
func (w *ResticWrapper) hasExclusiveLock(repository string) (bool, string, error) {
	ids, err := w.getLockIDs(repository)
	if err != nil {
		return false, "", fmt.Errorf("failed to list locks: %w", err)
	}

	if len(ids) == 0 {
		return false, "", nil
	}

	// Check each lock to find exclusive locks
	for _, id := range ids {
		st, err := w.getLockStats(repository, id)
		if err != nil {
			klog.Warningf("Failed to inspect lock %s: %v", id, err)
			continue
		}

		if st.Exclusive {
			klog.Infof("Found exclusive lock: %s (hostname: %s)", id, st.Hostname)
			return true, st.Hostname, nil
		}
	}

	return false, "", nil
}

// EnsureNoExclusiveLock ensures the repository is ready for a new operation by:
// 1. Removing all stale locks (restic determines which locks are stale)
// 2. Checking if any exclusive locks remain (if they do, they're active)
// 3. Waiting for active exclusive locks to be released
// Reference: https://forum.restic.net/t/locks-being-created-and-not-cleared/4836/3
func (w *ResticWrapper) EnsureNoExclusiveLock(rClient client.Client, namespace string) error {
	klog.Infoln("Ensuring repository is ready for new operation...")

	for _, b := range w.Config.Backends {
		klog.Infof("Processing repository: %s", b.Repository)

		// Remove stale locks
		klog.Infof("Removing stale locks from repository: %s", b.Repository)
		_, err := w.unlockStale(b.Repository)
		if err != nil {
			klog.Warningf("Failed to remove stale locks (non-fatal): %v", err)
		}

		// Check if any exclusive locks remain
		// If they do, restic determined they're active (it would have removed them if stale)
		klog.Infof("Checking for exclusive locks in repository: %s", b.Repository)
		hasLock, podName, err := w.hasExclusiveLock(b.Repository)
		if err != nil {
			return fmt.Errorf("failed to check for exclusive locks in repository %s: %w", b.Repository, err)
		}

		if !hasLock {
			klog.Infof("No exclusive lock found. Repository %s is ready.", b.Repository)
			continue
		}

		// : Wait for the exclusive lock to be released
		// Periodically retry unlockStale() in case the process crashes during wait
		const lockWaitTimeout = 1 * time.Hour

		klog.Infof("Exclusive lock found (held by %s). Waiting up to %v for it to be released...", podName, lockWaitTimeout)
		err = wait.PollUntilContextTimeout(
			context.Background(),
			10*time.Second,
			lockWaitTimeout,
			true,
			func(ctx context.Context) (bool, error) {
				klog.Infof("Polling: checking if exclusive lock is released...")

				// Try to cleanup stale locks (in case process crashed)
				_, unlockErr := w.unlockStale(b.Repository)
				if unlockErr != nil {
					klog.Warningf("Failed to remove stale locks during polling: %v", unlockErr)
				}

				// Check if exclusive lock still exists
				hasLock, currentPodName, err := w.hasExclusiveLock(b.Repository)
				if err != nil {
					klog.Warningf("Error checking locks during polling: %v", err)
					return false, nil // Don't fail, retry
				}

				if !hasLock {
					klog.Infof("Exclusive lock released. Repository is ready.")
					return true, nil
				}

				// Lock still exists
				klog.Infof("Exclusive lock still held by %s. Waiting...", currentPodName)
				return false, nil
			},
		)
		if err != nil {
			return fmt.Errorf("timeout waiting for exclusive lock to be released in repository %s: %w", b.Repository, err)
		}

		klog.Infof("Repository %s is ready.", b.Repository)
	}

	klog.Infoln("All repositories are ready for new operations.")
	return nil
}

/*
Link: https://restic.readthedocs.io/en/v0.4.0/Design/#locks

Exclusive Locks
- Only one exclusive lock can run at a time.
- It blocks all non-exclusive locks.

Non-Exclusivity Locks
- Multiple non-exclusive locks can run at the same time.
- They do not block other non-exclusive locks.
- They do block exclusive locks (writers).
*/
