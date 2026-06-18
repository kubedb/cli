package sh

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var ErrExecTimeout = errors.New("execute timeout")

// unmarshal shell output to decode json
func (s *Session) UnmarshalJSON(data interface{}) (err error) {
	oldout, oldbuf := s.Stdout, s.enableOutputBuffer
	defer func() {
		s.Stdout, s.enableOutputBuffer = oldout, oldbuf
	}()
	bufrw := bytes.NewBuffer(nil)
	s.Stdout, s.enableOutputBuffer = bufrw, true
	err = s.Run()
	err = errors.Join(err, s.writeCmdOutputToStdOut())
	if err != nil {
		return
	}
	return json.NewDecoder(bufrw).Decode(data)
}

// unmarshal command output into xml
func (s *Session) UnmarshalXML(data interface{}) (err error) {
	oldout, oldbuf := s.Stdout, s.enableOutputBuffer
	defer func() {
		s.Stdout, s.enableOutputBuffer = oldout, oldbuf
	}()
	bufrw := bytes.NewBuffer(nil)
	s.Stdout, s.enableOutputBuffer = bufrw, true
	err = s.Run()
	err = errors.Join(err, s.writeCmdOutputToStdOut())
	if err != nil {
		return
	}
	return xml.NewDecoder(bufrw).Decode(data)
}

// start command
func (s *Session) Start() (err error) {
	s.started = true
	if s.ShowCMD {
		s.displayCommandChain()
	}

	if len(s.cmds) == 0 {
		return s.executeLeafCommands(nil)
	}
	return s.executeCommandChain(0, nil)
}

func (s *Session) executeCommandChain(index int, stdin *io.PipeReader) error {
	if index >= len(s.cmds) {
		return nil
	}
	pipeReaders, pipeWriters := createPipes(s.determinePipeCount(index))

	cmd := s.cmds[index]
	cmd.Stdin = s.selectCmdStdin(index, stdin)
	cmd.Stdout, cmd.Stderr = s.configureCmdOutput(index, pipeWriters)

	s.pipeWriters = append(s.pipeWriters, pipeWriters...)
	if err := cmd.Start(); err != nil {
		return err
	}

	if s.isLastCommand(index) && len(s.leafCmds) != 0 {
		return s.executeLeafCommands(pipeReaders)
	}
	return s.executeCommandChain(index+1, pipeReaders[0])
}

func (s *Session) selectCmdStdin(index int, stdin *io.PipeReader) io.Reader {
	if index == 0 {
		return s.Stdin
	}
	return stdin
}

func (s *Session) configureCmdOutput(index int, pipeWriters []*io.PipeWriter) (io.Writer, io.Writer) {
	if s.isLastCommand(index) && len(s.leafCmds) == 0 {
		if s.enableOutputBuffer {
			cmdOutput := &safeBuffer{}
			s.lastOutputBuffer = cmdOutput
			if s.enableErrsBuffer {
				return cmdOutput, cmdOutput
			}
			return cmdOutput, s.Stderr
		}
		return s.Stdout, s.Stderr
	}

	stdout := io.MultiWriter(pipeWritersToWriters(pipeWriters)...)
	var stderr io.Writer = os.Stderr
	if s.PipeStdErrors {
		stderr = s.Stderr
	}

	return stdout, stderr
}

func pipeWritersToWriters(pipeWriters []*io.PipeWriter) []io.Writer {
	var writers []io.Writer
	for _, writer := range pipeWriters {
		writers = append(writers, writer)
	}
	return writers
}

func (s *Session) executeLeafCommands(readers []*io.PipeReader) error {
	for idx, cmd := range s.leafCmds {
		cmd.Stdin = s.selectLeafCmdStdin(idx, readers)
		cmd.Stdout = s.selectLeafCmdStdout()
		cmd.Stderr = s.selectLeafCmdStderr()

		if err := cmd.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) selectLeafCmdStdin(index int, readers []*io.PipeReader) io.Reader {
	if readers != nil {
		return readers[index]
	}
	return s.Stdin
}

func (s *Session) selectLeafCmdStderr() io.Writer {
	if s.enableErrsBuffer {
		return s.selectLeafCmdStdout()
	}
	return s.Stderr
}

func (s *Session) selectLeafCmdStdout() io.Writer {
	if s.enableOutputBuffer {
		cmdOutput := &safeBuffer{}
		s.leafOutputBuffer = append(s.leafOutputBuffer, cmdOutput)
		return cmdOutput
	}
	return os.Stdout
}

func createPipes(count int) ([]*io.PipeReader, []*io.PipeWriter) {
	readers := make([]*io.PipeReader, count)
	writers := make([]*io.PipeWriter, count)

	for i := 0; i < count; i++ {
		r, w := io.Pipe()
		readers[i] = r
		writers[i] = w
	}

	return readers, writers
}

func (s *Session) determinePipeCount(index int) int {
	if s.isLastCommand(index) && len(s.leafCmds) != 0 {
		return len(s.leafCmds)
	}
	return 1
}

func (s *Session) isLastCommand(index int) bool {
	return index == len(s.cmds)-1
}

func (s *Session) displayCommandChain() {
	joinCmds := func(cmds []*exec.Cmd) []string {
		result := make([]string, len(cmds))
		for i, cmd := range cmds {
			result[i] = strings.Join(cmd.Args, " ")
		}
		return result
	}
	primaryCmds, leafCmds := joinCmds(s.cmds), joinCmds(s.leafCmds)

	totalCmd := strings.Join(primaryCmds, " | ")
	if len(leafCmds) > 0 {
		totalCmd += " | " + strings.Join(leafCmds, " , ")
	}

	s.writePrompt(totalCmd)
}

// Should be call after Start()
// only catch the last command error
func (s *Session) Wait() error {
	var pipeErr, lastErr error
	for idx, writter := range s.pipeWriters {
		if idx < len(s.cmds) {
			cmd := s.cmds[idx]
			if lastErr = cmd.Wait(); lastErr != nil {
				pipeErr = lastErr
			}
		}
		writter.Close()
	}
	var combineErrs []error
	for _, cmd := range s.leafCmds {
		if err := cmd.Wait(); err != nil {
			combineErrs = append(combineErrs, err)
		}
	}

	if s.PipeFail {
		return pipeErr
	}

	combineErrs = append([]error{pipeErr}, combineErrs...)
	return errors.Join(combineErrs...)
}

func (s *Session) Kill(sig os.Signal) {
	for _, cmd := range s.cmds {
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}
	}
}

func (s *Session) WaitTimeout(timeout time.Duration) (err error) {
	select {
	case <-time.After(timeout):
		s.Kill(syscall.SIGKILL)
		return ErrExecTimeout
	case err = <-Go(s.Wait):
		return err
	}
}

func Go(f func() error) chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- f()
	}()
	return ch
}

func (s *Session) Run() (err error) {
	s.resetOutputBuffer()
	if err = s.Start(); err != nil {
		return
	}
	if s.timeout != time.Duration(0) {
		return s.WaitTimeout(s.timeout)
	}
	return s.Wait()
}

func (s *Session) Output() (out []byte, err error) {
	oldout, oldbuf := s.Stdout, s.enableOutputBuffer
	defer func() {
		s.Stdout, s.enableOutputBuffer = oldout, oldbuf
	}()
	stdout := bytes.NewBuffer(nil)
	s.Stdout = stdout
	s.enableOutputBuffer = true
	err = s.Run()
	err = errors.Join(err, s.writeCmdOutputToStdOut())
	out = stdout.Bytes()
	return
}

func (s *Session) WriteStdout(f string) error {
	oldout, oldbuf := s.Stdout, s.enableOutputBuffer
	defer func() {
		s.Stdout, s.enableOutputBuffer = oldout, oldbuf
	}()

	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer out.Close()
	s.Stdout = out
	s.enableOutputBuffer = true
	err = s.Run()
	err = errors.Join(err, s.writeCmdOutputToStdOut())
	return err
}

func (s *Session) AppendStdout(f string) error {
	oldout, oldbuf := s.Stdout, s.enableOutputBuffer
	defer func() {
		s.Stdout, s.enableOutputBuffer = oldout, oldbuf
	}()

	out, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	s.Stdout = out
	s.enableOutputBuffer = true
	err = s.Run()
	err = errors.Join(err, s.writeCmdOutputToStdOut())
	return err
}

func (s *Session) CombinedOutput() (out []byte, err error) {
	oldout, olderr := s.Stdout, s.Stderr
	oldbuf, olderrbuf := s.enableOutputBuffer, s.enableErrsBuffer
	defer func() {
		s.Stdout, s.Stderr = oldout, olderr
		s.enableOutputBuffer, s.enableErrsBuffer = oldbuf, olderrbuf
	}()
	stdout := bytes.NewBuffer(nil)
	s.Stdout = stdout
	s.Stderr = stdout

	s.enableErrsBuffer = true
	s.enableOutputBuffer = true
	err = s.Run()
	err = errors.Join(err, s.writeCmdOutputToStdOut())
	out = stdout.Bytes()
	return
}

func (s *Session) writeCmdOutputToStdOut() error {
	var errs []error
	for _, buffer := range s.leafOutputBuffer {
		_, err := s.Stdout.Write(buffer.Bytes())
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(s.leafOutputBuffer) == 0 && s.lastOutputBuffer != nil {
		_, err := s.Stdout.Write(s.lastOutputBuffer.Bytes())
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (s *Session) resetOutputBuffer() {
	for _, buffer := range s.leafOutputBuffer {
		buffer.Reset()
	}
	if s.lastOutputBuffer != nil {
		s.lastOutputBuffer.Reset()
	}
}

// CurrentOutput returns a snapshot of command output at the given index.
// If leaf commands exist, index maps to the leaf command index.
// Otherwise, only index 0 is valid and returns the last command output.
func (s *Session) CurrentOutput(index int) ([]byte, error) {
	if len(s.leafOutputBuffer) > 0 {
		if index < 0 || index >= len(s.leafOutputBuffer) {
			return nil, fmt.Errorf("index %d out of range [0, %d)", index, len(s.leafOutputBuffer))
		}
		return s.leafOutputBuffer[index].Bytes(), nil
	}

	if s.lastOutputBuffer == nil {
		return nil, fmt.Errorf("no output buffer available; ensure the command was started with buffered output")
	}
	if index != 0 {
		return nil, fmt.Errorf("index %d out of range: only index 0 is valid for non-leaf command chains", index)
	}
	return s.lastOutputBuffer.Bytes(), nil
}

// CurrentLeafOutput is kept for backward compatibility. Use CurrentOutput instead.
func (s *Session) CurrentLeafOutput(index int) ([]byte, error) {
	return s.CurrentOutput(index)
}
