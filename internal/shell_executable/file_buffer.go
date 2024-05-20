package shell_executable

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/async_buffered_reader"
)

type FileBuffer struct {
	descriptor     *os.File
	bufferedReader *async_buffered_reader.AsyncBufferedReader
}

func NewFileBuffer(descriptor *os.File) FileBuffer {
	return FileBuffer{
		descriptor:     descriptor,
		bufferedReader: async_buffered_reader.New(bufio.NewReader(descriptor)),
	}
}

func (t *FileBuffer) ReadBuffer(shouldStopReadingBuffer func([]byte) error) ([]byte, error) {
	return t.ReadBufferWithTimeout(2000*time.Millisecond, shouldStopReadingBuffer)
}

func (t *FileBuffer) ReadBufferWithTimeout(timeout time.Duration, shouldStopReadingBuffer func([]byte) error) ([]byte, error) {
	data, err := t.readUntil(shouldStopReadingBuffer, timeout)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (t *FileBuffer) readUntil(condition func([]byte) error, timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)
	readBytes := []byte{}

	for !time.Now().After(deadline) {
		readByte, err := t.bufferedReader.ReadByteWithTimeout(2 * time.Millisecond)
		if err != nil {
			if err == async_buffered_reader.ErrNoData {
				// fmt.Println("No data available")
				time.Sleep(2 * time.Millisecond) // Let's wait a bit before trying again
				continue
			} else {
				// fmt.Printf("Error while reading: %v\n", err)
				return readBytes, err
			}
		}

		// fmt.Printf("readByte: %q\n", string(readByte))
		readBytes = append(readBytes, readByte)

		// If the condition is met, return. Else the loop runs again
		if condition(readBytes) == nil {
			return readBytes, nil
		}
	}

	// TODO: Use a better error message here?
	return readBytes, fmt.Errorf("timeout while waiting for condition")
}

func StripANSI(data []byte) []byte {
	// https://github.com/acarl005/stripansi/blob/master/stripansi.go
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	var re = regexp.MustCompile(ansi)

	return re.ReplaceAll(data, []byte(""))
}
