package shell_executable

import (
	"fmt"
	"time"
)

func (b *ShellExecutable) ReadBuffer(selector string) ([]byte, error) {
	// selector can be either "stdout" or "stderr"
	// XXX: We might have to read both buffers, and return the one holding the data
	// Unsure if we can rely on users to sending the data to the correct pipe
	return b.ReadBufferWithTimeout(10*time.Millisecond, selector)
}

// Use it like this:
// buffer, err := b.ReadBufferCustom(50*time.Millisecond, "stderr", func(b []byte) bool { return len(b) > 50*(i+1) })
func (b *ShellExecutable) ReadBufferCustom(timeout time.Duration, selector string, shouldStopReadingBuffer func([]byte) bool) ([]byte, error) {
	data, err := b.readUntil(shouldStopReadingBuffer, timeout, selector)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (b *ShellExecutable) ReadBufferWithTimeout(timeout time.Duration, selector string) ([]byte, error) {
	shouldStopReadingBuffer := func(buf []byte) bool {
		if len(buf) < 2 {
			return false
		}
		// After completing the current command, the shell would move on to the next line with the prompt
		// XXX : What about users without this functionality ?
		return string(buf[len(buf)-2:]) == "$ "
	}

	data, err := b.readUntil(shouldStopReadingBuffer, timeout, selector)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (b *ShellExecutable) readUntil(condition func([]byte) bool, timeout time.Duration, selector string) ([]byte, error) {
	deadline := time.Now().Add(timeout)

	for !time.Now().After(deadline) {
		time.Sleep(1 * time.Millisecond) // Let's give some time for the buffer to fill up

		var readData []byte
		switch selector {
		case "stdout":
			readData = b.executable.StdoutBuffer.Bytes()
		case "stderr":
			readData = b.executable.StderrBuffer.Bytes()
		}

		if condition(readData) {
			return readData, nil
		} else {
			time.Sleep(1 * time.Millisecond) // Let's wait a bit before trying again
		}
	}
	return nil, fmt.Errorf("timeout while waiting for condition")
}
