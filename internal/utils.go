package internal

import "github.com/codecrafters-io/tester-utils/executable"

type ExecutableResultDetailed struct {
	stdout          []byte
	stderr          []byte
	currentStdout   []byte
	currentStderr   []byte
	stdoutBytesRead int
	stderrBytesRead int
	exitCode        int
}

func NewDetailedResult(result executable.ExecutableResult) *ExecutableResultDetailed {
	return &ExecutableResultDetailed{
		stdout:          result.Stdout,
		stderr:          result.Stderr,
		currentStdout:   result.Stdout,
		currentStderr:   result.Stderr,
		stdoutBytesRead: 0,
		stderrBytesRead: 0,
		exitCode:        result.ExitCode,
	}
}

func (r *ExecutableResultDetailed) CurrentCommandStdOut(skipNextLinePrompt bool) []byte {
	endIndex := len(r.stdout)
	if skipNextLinePrompt {
		if endIndex >= 2 {
			endIndex -= 2
		}
	}
	data := r.stdout[r.stdoutBytesRead:endIndex]
	r.stdoutBytesRead = len(r.stdout)
	return data
}

func (r *ExecutableResultDetailed) CurrentCommandStdErr(skipNextLinePrompt bool) []byte {
	endIndex := len(r.stderr)
	if skipNextLinePrompt {
		if endIndex >= 2 {
			endIndex -= 2
		}
	}
	data := r.stderr[r.stderrBytesRead:endIndex]
	r.stderrBytesRead = len(r.stderr)
	return data
}

func (r *ExecutableResultDetailed) UpdateData(result executable.ExecutableResult) error {
	r.stdout = result.Stdout
	r.stderr = result.Stderr
	r.exitCode = result.ExitCode
	return nil
}
