package shell_executable

import (
	"errors"
	"io"
	"os"
	"syscall"
)

const ptyReadBufferSize = 4096

// ptyRelay continuously drains bytes from the PTY and writes them into the virtual terminal.
//
// This serves two purposes:
//  1. It keeps the PTY kernel buffer drained so the child process never blocks on write.
//  2. It keeps the virtual terminal up-to-date so assertion polling always sees current screen state.
//
// The relay runs in a dedicated goroutine. When the PTY signals EOF or EIO (process exited),
// the relay stores the terminal error and closes relayExited, which callers can select on.
type ptyRelay struct {
	ptyFile         *os.File
	virtualTerminal io.Writer

	// relayExited is closed by the goroutine when it exits (either on process exit or error).
	relayExited chan bool

	// terminalErr holds the error that caused the relay to stop.
	// It is written exactly once before relayExited is closed, so readers
	// that receive from relayExited can safely read terminalErr without a mutex.
	terminalErr error
}

func newPtyRelay(ptyFile *os.File, virtualTerminal io.Writer) *ptyRelay {
	return &ptyRelay{
		ptyFile:         ptyFile,
		virtualTerminal: virtualTerminal,
		relayExited:     make(chan bool),
	}
}

// start launches the relay goroutine. Call this once after creating the relay.
func (r *ptyRelay) start() {
	go r.run()
}

func (r *ptyRelay) run() {
	buf := make([]byte, ptyReadBufferSize)

	for {
		bytesRead, readErr := r.ptyFile.Read(buf)

		if bytesRead > 0 {
			// Write to the virtual terminal before handling the error, because on some
			// platforms (Linux) the final data arrives in the same call that returns EIO.
			r.virtualTerminal.Write(buf[:bytesRead]) //nolint:errcheck
		}

		if readErr != nil {
			r.terminalErr = readErr
			close(r.relayExited)
			return
		}
	}
}

// processExited returns true if the relay has detected that the child process exited.
func (r *ptyRelay) processExited() bool {
	select {
	case <-r.relayExited:
		return true
	default:
		return false
	}
}

// isPtyTerminalError returns true if the error signals that the PTY's child process has exited.
// Linux returns EIO; macOS returns EOF.
func isPtyTerminalError(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EIO)
}
