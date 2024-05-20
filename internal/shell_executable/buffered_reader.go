package shell_executable

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

type FileBuffer struct {
	descriptor   *os.File
	completeData []byte
}

func NewFileBuffer(descriptor *os.File) FileBuffer {
	return FileBuffer{
		descriptor:   descriptor,
		completeData: make([]byte, 0),
	}
}

func (b FileBuffer) feedStdin(command []byte) error {
	_, err := b.descriptor.Write(command)
	// b.logger.Debugf("Wrote %d bytes to stdin", n)
	if err != nil {
		return err
	}
	return nil
}

func (b FileBuffer) FeedStdin(command []byte) error {
	commandWithEnter := append(command, []byte("\n")...)
	return b.feedStdin(commandWithEnter)
}

func (t *FileBuffer) ReadBuffer(shouldStopReadingBuffer func(string, []byte) error, expectedValue string) ([]byte, error) {
	return t.ReadBufferWithTimeout(10*time.Millisecond, shouldStopReadingBuffer, expectedValue)
}

func (t *FileBuffer) ReadBufferWithTimeout(timeout time.Duration, shouldStopReadingBuffer func(string, []byte) error, expectedValue string) ([]byte, error) {
	data, err := t.readUntil(shouldStopReadingBuffer, timeout, expectedValue)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *FileBuffer) readUntil(condition func(string, []byte) error, timeout time.Duration, expectedValue string) ([]byte, error) {
	deadline := time.Now().Add(timeout)

	for !time.Now().After(deadline) {
		time.Sleep(5 * time.Millisecond) // Let's give some time for the buffer to fill up

		fullBuf := []byte{}
		buf := make([]byte, 1024)
		n, err := t.descriptor.Read(buf)
		// fmt.Println("n: ", n, "err: ", err, "buf: ", buf[:n], string(buf[:n]))
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fullBuf = append(fullBuf, buf[:n]...)
		// cleanedFullBuf := StripANSI(fullBuf)
		cleanedFullBuf := fullBuf

		if errToBool(condition(expectedValue, cleanedFullBuf)) {
			t.completeData = append(t.completeData, cleanedFullBuf...)
			return cleanedFullBuf, nil
		} else {
			time.Sleep(2 * time.Millisecond) // Let's wait a bit before trying again
		}
	}
	return nil, fmt.Errorf("timeout while waiting for condition")
}

func RemoveControlSequence(data []byte) []byte {
	PROMPT_START := '$'

	for startIdx, r := range string(data) {
		// Standard escape codes are prefixed with Escape (27)
		if r == 27 {
			// remove from here upto PROMPT_START
			for endIdx, r2 := range string(data[startIdx:]) {
				if r2 == PROMPT_START {
					// Remove from start_idx to end_idx-1
					data = append(data[:startIdx], data[endIdx:]...)
					break
				}
			}
		}
	}

	return data
}

func RemoveControlSequences(data []byte) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

	line := re.ReplaceAll(data, []byte(""))
	return string(line)
}

func StripANSI(data []byte) []byte {
	// https://github.com/acarl005/stripansi/blob/master/stripansi.go
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	var re = regexp.MustCompile(ansi)

	return re.ReplaceAll(data, []byte(""))
}

func errToBool(err error) bool {
	return err == nil
}
