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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"k8s.io/klog/v2"
)

const (
	ResticCMD = "restic"
)

type Snapshot struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	Tree     string    `json:"tree"`
	Paths    []string  `json:"paths"`
	Hostname string    `json:"hostname"`
	Username string    `json:"username"`
	UID      int       `json:"uid"`
	Gid      int       `json:"gid"`
	Tags     []string  `json:"tags"`
}

type backupParams struct {
	path     string
	host     string
	tags     []string
	excludes []string
	args     []string
}

type restoreParams struct {
	path        string
	host        string
	snapshotId  string
	destination string
	excludes    []string
	includes    []string
	args        []string
}

type keyParams struct {
	id   string
	user string
	host string
	file string
}

func (w *ResticWrapper) listSnapshots(repository string, snapshotIDs []string) ([]Snapshot, error) {
	result := make([]Snapshot, 0)
	args := w.appendCacheDirFlag([]any{"snapshots", "--json", "--quiet", "--no-lock"})
	b := w.getMatchedBackend(repository)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = append(args, b.Envs)
	for _, id := range snapshotIDs {
		args = append(args, id)
	}
	out, err := w.run(Command{Name: ResticCMD, Args: args})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(out, &result)
	return result, err
}

func (w *ResticWrapper) tryDeleteSnapshots(repository string, snapshotIDs []string) ([]byte, error) {
	args := w.appendCacheDirFlag([]any{"forget", "--quiet", "--prune"})
	b := w.getMatchedBackend(repository)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = append(args, b.Envs)
	for _, id := range snapshotIDs {
		args = append(args, id)
	}
	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) deleteSnapshots(repository string, snapshotIDs []string) ([]byte, error) {
	out, err := w.tryDeleteSnapshots(repository, snapshotIDs)
	if err == nil || !strings.Contains(err.Error(), "unlock") {
		return out, err
	}
	// repo is locked, so unlock first
	klog.Warningln("repo found locked, so unlocking before pruning, err:", err.Error())
	if out, err = w.unlock(repository); err != nil {
		return out, err
	}
	return w.tryDeleteSnapshots(repository, snapshotIDs)
}

func (w *ResticWrapper) repositoryExist(repository string) bool {
	klog.Infoln("Checking whether the backend repository exist or not....")
	b := w.getMatchedBackend(repository)
	args := w.appendCacheDirFlag([]any{"snapshots", "--json", "--no-lock"})
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = append(args, b.Envs)
	if _, err := w.run(Command{Name: ResticCMD, Args: args}); err == nil {
		return true
	}
	return false
}

func (w *ResticWrapper) getMatchedBackend(repository string) *Backend {
	// Use index for O(1) lookup if available
	if b := w.Config.GetBackend(repository); b != nil {
		return b
	}
	// Return an empty backend to avoid nil pointer dereference
	return &Backend{
		StorageConfig: &StorageConfig{},
		Envs:          make(map[string]string),
	}
}

func (w *ResticWrapper) initRepository(repository string) error {
	klog.Infoln("Initializing new restic repository for repository:", repository)
	b := w.getMatchedBackend(repository)
	if err := b.createLocalDir(); err != nil {
		return err
	}
	args := w.appendCacheDirFlag([]any{"init"})
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = append(args, b.Envs)
	_, err := w.run(Command{Name: ResticCMD, Args: args})
	return err
}

func (w *ResticWrapper) backup(params backupParams) ([]byte, error) {
	klog.Infoln("Backing up target data")
	var commands []Command
	commonArgs := []any{"backup", params.path, "--quiet", "--json"}
	if params.host != "" {
		commonArgs = append(commonArgs, "--host")
		commonArgs = append(commonArgs, params.host)
	}
	// add tags if any
	for _, tag := range params.tags {
		commonArgs = append(commonArgs, "--tag")
		commonArgs = append(commonArgs, tag)
	}
	// add exclude patterns if there is any
	for _, exclude := range params.excludes {
		commonArgs = append(commonArgs, "--exclude")
		commonArgs = append(commonArgs, exclude)
	}
	// add additional arguments passed by user to the backup process
	for i := range params.args {
		commonArgs = append(commonArgs, params.args[i])
	}
	commonArgs = w.appendCacheDirFlag(commonArgs)
	commonArgs = w.appendCleanupCacheFlag(commonArgs)

	for _, b := range w.Config.Backends {
		args := make([]any, len(commonArgs))
		copy(args, commonArgs)
		args = b.appendCaCertFlag(args)
		args = b.appendInsecureTLSFlag(args)
		args = b.appendMaxConnectionsFlag(args)
		args = append(args, b.Envs)
		commands = append(commands, Command{Name: ResticCMD, Args: args})
	}
	return w.run(commands...)
}

func (w *ResticWrapper) backupFromStdin(options BackupOptions) ([]byte, error) {
	klog.Infoln("Backing up stdin data")

	// first add StdinPipeCommands, then add restic command
	commands := options.StdinPipeCommands

	commonArgs := []any{"backup", "--stdin", "--quiet", "--json"}
	commonArgs = options.appendHost(commonArgs)
	commonArgs = options.appendStdinFileName(commonArgs)
	commonArgs = w.appendCacheDirFlag(commonArgs)
	commonArgs = w.appendCleanupCacheFlag(commonArgs)

	for _, b := range w.Config.Backends {
		args := make([]any, len(commonArgs))
		copy(args, commonArgs)
		args = b.appendCaCertFlag(args)
		args = b.appendInsecureTLSFlag(args)
		args = b.appendMaxConnectionsFlag(args)
		args = append(args, b.Envs)
		command := Command{Name: ResticCMD, Args: args}
		commands = append(commands, command)
	}
	return w.run(commands...)
}

func (w *ResticWrapper) restore(repository string, params restoreParams) ([]byte, error) {
	klog.Infoln("Restoring backed up data")

	args := []any{"restore"}
	if params.snapshotId != "" {
		args = append(args, params.snapshotId)
	} else {
		args = append(args, "latest")
	}
	if params.path != "" {
		args = append(args, "--path")
		args = append(args, params.path) // source-path specified in restic fileGroup
	}
	if params.host != "" {
		args = append(args, "--host")
		args = append(args, params.host)
	}

	if params.destination == "" {
		params.destination = "/" // restore in absolute path
	}
	args = append(args, "--target", params.destination)

	// add include patterns if there is any
	for _, include := range params.includes {
		args = append(args, "--include")
		args = append(args, include)
	}
	// add exclude patterns if there is any
	for _, exclude := range params.excludes {
		args = append(args, "--exclude")
		args = append(args, exclude)
	}
	// add additional arguments passed by user to the restore process
	for i := range params.args {
		args = append(args, params.args[i])
	}
	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = append(args, b.Envs)
	command := Command{Name: ResticCMD, Args: args}

	return w.run(command)
}

func (w *ResticWrapper) DumpOnce(repository string, dumpOptions DumpOptions) ([]byte, error) {
	klog.Infoln("Dumping backed up data")

	args := []any{"dump", "--quiet"}
	if dumpOptions.Snapshot != "" {
		args = append(args, dumpOptions.Snapshot)
	} else {
		args = append(args, "latest")
	}
	if dumpOptions.FileName != "" {
		args = append(args, dumpOptions.FileName)
	} else {
		args = append(args, "stdin")
	}
	if dumpOptions.SourceHost != "" {
		args = append(args, "--host")
		args = append(args, dumpOptions.SourceHost)
	}
	if dumpOptions.Path != "" {
		args = append(args, "--path")
		args = append(args, dumpOptions.Path)
	}
	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	command := Command{Name: ResticCMD, Args: args}

	// first add restic command, then add StdoutPipeCommands
	commands := []Command{command}
	commands = append(commands, dumpOptions.StdoutPipeCommands...)
	return w.run(commands...)
}

func (w *ResticWrapper) check(repository string) ([]byte, error) {
	klog.Infoln("Checking integrity of repository")
	args := w.appendCacheDirFlag([]any{"check", "--no-lock"})
	b := w.getMatchedBackend(repository)
	args = b.appendCaCertFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)
	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) stats(repository string, snapshotID string) ([]byte, error) {
	klog.Infoln("Reading repository status")
	args := w.appendCacheDirFlag([]any{"stats"})
	if snapshotID != "" {
		args = append(args, snapshotID)
	}
	b := w.getMatchedBackend(repository)
	args = b.appendMaxConnectionsFlag(args)
	args = append(args, "--quiet", "--json", "--mode", "raw-data", "--no-lock")
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) unlock(repository string) ([]byte, error) {
	klog.Infoln("Unlocking restic repository")
	args := w.appendCacheDirFlag([]any{"unlock", "--remove-all"})
	b := w.getMatchedBackend(repository)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)
	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) unlockStale(repository string) ([]byte, error) {
	klog.Infoln("Removing stale locks from restic repository")
	args := w.appendCacheDirFlag([]any{"unlock"})
	b := w.getMatchedBackend(repository)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)
	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) appendCacheDirFlag(args []any) []any {
	if w.Config.EnableCache {
		cacheDir := filepath.Join(w.Config.ScratchDir, resticCacheDir)
		return append(args, "--cache-dir", cacheDir)
	}
	return append(args, "--no-cache")
}

func (opt *BackupOptions) appendStdinFileName(args []any) []any {
	if opt.StdinFileName != "" {
		args = append(args, "--stdin-filename")
		args = append(args, opt.StdinFileName)
	}
	return args
}

func (opt *BackupOptions) appendHost(args []any) []any {
	if opt.Host != "" {
		args = append(args, "--host")
		args = append(args, opt.Host)
	}
	return args
}

func (w *ResticWrapper) appendCleanupCacheFlag(args []any) []any {
	if w.Config.EnableCache {
		return append(args, "--cleanup-cache")
	}
	return args
}

func (w *ResticWrapper) run(commands ...Command) ([]byte, error) {
	// write std errors into os.Stderr and buffer
	errBuff, err := circbuf.NewBuffer(256)
	if err != nil {
		return nil, err
	}

	newSh := *w.sh // Create a new shell instance to avoid pollution from existing environment variables.
	newSh.Stderr = io.MultiWriter(os.Stderr, errBuff)
	if w.Config.Timeout != nil {
		newSh.SetTimeout(w.Config.Timeout.Duration)
	}

	isLeafCommandRequired := isLeafCommandNecessary(commands...)
	for _, cmd := range commands {
		var useLeafCommand bool
		if isLeafCommandRequired && cmd.Name == ResticCMD {
			useLeafCommand = true
		}

		cmd, err = w.applyNiceSettingsIfCommandMatches(cmd, ResticCMD)
		if err != nil {
			return nil, err
		}
		if useLeafCommand {
			newSh.LeafCommand(cmd.Name, cmd.Args...)
		} else {
			newSh.Command(cmd.Name, cmd.Args...)
		}
	}

	out, err := newSh.Output()
	if err != nil {
		return nil, formatError(err, errBuff.String())
	}
	klog.Infoln("sh-output:", string(out))
	return out, nil
}

func (w *ResticWrapper) applyNiceSettingsIfCommandMatches(command Command, matchingCommands ...string) (Command, error) {
	var err error
	if slices.Contains(matchingCommands, command.Name) {
		// First apply standard settings, then apply I/O priority settings
		command, err = w.applyNiceSettings(command)
		if err != nil {
			return Command{}, err
		}
		command, err = w.applyIONiceSettings(command)
		if err != nil {
			return Command{}, err
		}
	}
	return command, nil
}

// return last line of std error as error reason
func formatError(err error, stdErr string) error {
	parts := strings.Split(strings.TrimSuffix(stdErr, "\n"), "\n")
	if len(parts) > 1 {
		if strings.Contains(parts[1], "signal terminated") {
			return errors.New(strings.Join(append([]string{"deadline exceeded or signal terminated"}, parts[2:]...), " "))
		}
		return errors.New(strings.Join(parts[1:], " "))
	}
	return err
}

func (w *ResticWrapper) applyIONiceSettings(oldCommand Command) (Command, error) {
	if w.Config.IONice == nil {
		return oldCommand, nil
	}

	// detect "ionice" installation path
	IONiceCMD, err := exec.LookPath("ionice")
	if err != nil {
		return Command{}, err
	}
	newCommand := Command{
		Name: IONiceCMD,
	}
	if w.Config.IONice.Class != nil {
		newCommand.Args = append(newCommand.Args, "-c", fmt.Sprint(*w.Config.IONice.Class))
	}
	if w.Config.IONice.ClassData != nil {
		newCommand.Args = append(newCommand.Args, "-n", fmt.Sprint(*w.Config.IONice.ClassData))
	}
	// TODO: should we use "-t" option with ionice ?
	// newCommand.Args = append(newCommand.Args, "-t")

	// append oldCommand as args of newCommand
	newCommand.Args = append(newCommand.Args, oldCommand.Name)
	newCommand.Args = append(newCommand.Args, oldCommand.Args...)
	return newCommand, nil
}

func (w *ResticWrapper) applyNiceSettings(oldCommand Command) (Command, error) {
	if w.Config.Nice == nil {
		return oldCommand, nil
	}

	// detect "nice" installation path
	NiceCMD, err := exec.LookPath("nice")
	if err != nil {
		return Command{}, err
	}
	newCommand := Command{
		Name: NiceCMD,
	}
	if w.Config.Nice.Adjustment != nil {
		newCommand.Args = append(newCommand.Args, "-n", fmt.Sprint(*w.Config.Nice.Adjustment))
	}

	// append oldCommand as args of newCommand
	newCommand.Args = append(newCommand.Args, oldCommand.Name)
	newCommand.Args = append(newCommand.Args, oldCommand.Args...)
	return newCommand, nil
}

func (w *ResticWrapper) addKey(repository string, params keyParams) ([]byte, error) {
	klog.Infoln("Adding new key to restic repository")

	args := []any{"key", "add", "--no-lock"}
	if params.host != "" {
		args = append(args, "--host", params.host)
	}

	if params.user != "" {
		args = append(args, "--user", params.user)
	}

	if params.file != "" {
		args = append(args, "--new-password-file", params.file)
	}

	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) listKey(repository string) ([]byte, error) {
	klog.Infoln("Listing restic keys")

	args := []any{"key", "list", "--no-lock"}

	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) listLocks(repository string) ([]byte, error) {
	klog.Infoln("Listing restic locks")

	args := []any{"list", "locks", "--no-lock"}

	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) lockStats(repository, lockID string) ([]byte, error) {
	klog.Infoln("Getting stats of restic lock")

	args := []any{"cat", "lock", lockID, "--no-lock"}

	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) updateKey(repository string, params keyParams) ([]byte, error) {
	klog.Infoln("Updating restic key")

	args := []any{"key", "passwd", "--no-lock"}

	if params.file != "" {
		args = append(args, "--new-password-file", params.file)
	}

	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func (w *ResticWrapper) removeKey(repository string, params keyParams) ([]byte, error) {
	klog.Infoln("Removing restic key")

	args := []any{"key", "remove", params.id, "--no-lock"}

	b := w.getMatchedBackend(repository)
	args = w.appendCacheDirFlag(args)
	args = b.appendMaxConnectionsFlag(args)
	args = b.appendCaCertFlag(args)
	args = b.appendInsecureTLSFlag(args)
	args = append(args, b.Envs)

	return w.run(Command{Name: ResticCMD, Args: args})
}

func isLeafCommandNecessary(commands ...Command) bool {
	var resticCommandCount int
	for _, command := range commands {
		if command.Name == ResticCMD {
			resticCommandCount++
		}
	}

	return resticCommandCount > 1
}
