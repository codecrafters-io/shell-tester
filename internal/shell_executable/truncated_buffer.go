package shell_executable

import (
	"bytes"
	"fmt"
	"time"
)

type TruncatedBuffer struct {
	buffer *bytes.Buffer
	// XXX
	// For updating offset we need the lastReadValueLength,
	// We can't rely on current buffer length, it will increase in an async manner
	// We also don't want to pass the lastReadValue's length as a param, so tracking it here makes sense
	lastReadValueLength int
	offset              int
}

func NewTruncatedBuffer(buffer *bytes.Buffer) TruncatedBuffer {
	return TruncatedBuffer{
		buffer: buffer,
		offset: 0,
	}
}

func (t *TruncatedBuffer) updateOffsetToCurrentLength() {
	t.offset = t.buffer.Len()
}

func (t *TruncatedBuffer) ReadBuffer() ([]byte, error) {
	return t.ReadBufferWithTimeout(10 * time.Millisecond)
}

func (t *TruncatedBuffer) ReadBufferWithTimeout(timeout time.Duration) ([]byte, error) {
	shouldStopReadingBuffer := func(buf []byte) bool {
		if len(buf) < 2 {
			return false
		}
		// After completing the current command, the shell would move on to the next line with the prompt
		// XXX : What about users without this functionality ?
		return string(buf[len(buf)-2:]) == "$ "
	}

	data, err := t.readUntil(shouldStopReadingBuffer, timeout)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TruncatedBuffer) readUntil(condition func([]byte) bool, timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)

	for !time.Now().After(deadline) {
		time.Sleep(1 * time.Millisecond) // Let's give some time for the buffer to fill up

		bytesData := t.buffer.Bytes()

		truncatedData := bytesData[t.offset:]
		t.lastReadValueLength = len(truncatedData)
		if condition(bytesData) {
			return truncatedData, nil
		} else {
			time.Sleep(1 * time.Millisecond) // Let's wait a bit before trying again
		}
	}
	return nil, fmt.Errorf("timeout while waiting for condition")
}

func (t *TruncatedBuffer) UpdateOffsetToCurrentLength() {
	t.offset += t.lastReadValueLength
}

// Duplicated for now
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
