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

func (t *TruncatedBuffer) ReadBuffer(shouldStopReadingBuffer func(string, []byte) error, expectedValue string) ([]byte, error) {
	return t.ReadBufferWithTimeout(10*time.Millisecond, shouldStopReadingBuffer, expectedValue)
}

func (t *TruncatedBuffer) ReadBufferWithTimeout(timeout time.Duration, shouldStopReadingBuffer func(string, []byte) error, expectedValue string) ([]byte, error) {
	data, err := t.readUntil(shouldStopReadingBuffer, timeout, expectedValue)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TruncatedBuffer) readUntil(condition func(string, []byte) error, timeout time.Duration, expectedValue string) ([]byte, error) {
	deadline := time.Now().Add(timeout)

	for !time.Now().After(deadline) {
		time.Sleep(1 * time.Millisecond) // Let's give some time for the buffer to fill up

		bytesData := t.buffer.Bytes()

		truncatedData := bytesData[t.offset:]
		t.lastReadValueLength = len(truncatedData)

		if errToBool(condition(expectedValue, bytesData)) {
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

func errToBool(err error) bool {
	return err == nil
}
