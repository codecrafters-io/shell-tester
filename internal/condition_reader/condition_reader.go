package condition_reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/async_reader"
)

var ErrConditionNotMet = errors.New("condition not met")

var debugLogsAreEnabled = false

func debugLog(format string, args ...interface{}) {
	if debugLogsAreEnabled {
		fmt.Printf(format+"\n", args...)
	}
}

// ConditionReader wraps an io.Reader and provides methods to read until a condition is met
type ConditionReader struct {
	asyncReader          *async_reader.AsyncReader
	nonAccumulatedBuffer []byte
}

func NewConditionReader(reader io.Reader) ConditionReader {
	return ConditionReader{
		asyncReader: async_reader.New(bufio.NewReader(reader)),
	}
}

func (t *ConditionReader) ReadUntilCondition(condition func([]byte) bool) ([]byte, error) {
	return t.ReadUntilConditionOrTimeout(condition, 2000*time.Millisecond)
}

func (t *ConditionReader) ReadUntilConditionOrTimeout(condition func([]byte) bool, timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)
	var accumulatedReadBytes []byte

	if len(t.nonAccumulatedBuffer) > 0 {
		accumulatedReadBytes = append(accumulatedReadBytes, t.nonAccumulatedBuffer...)
		t.nonAccumulatedBuffer = []byte{}
	}

	for !time.Now().After(deadline) {
		readBytes, err := t.asyncReader.ReadBytes()
		if err != nil {
			if errors.Is(err, async_reader.ErrNoData) {
				debugLog("condition_reader: No data available")

				// Since no data was available, let's avoid a busy loop
				time.Sleep(2 * time.Millisecond)

				continue
			} else {
				debugLog("condition_reader: Error while reading: %v", err)
				return readBytes, err
			}
		}

		// There might be a situation where we read more than the string `S` that satisfies the condition.
		// For that reason, we'll accumulate byte by byte and make sure we don't overshoot the condition.
		debugLog("condition_reader: readBytes: %q", string(readBytes))

		for i, byte := range readBytes {
			accumulatedReadBytes = append(accumulatedReadBytes, byte)
			// If the condition is met, we can return early. Else the loop runs again
			if condition(accumulatedReadBytes) {
				// Of the complete string `S`, if S[:i] satisfies the conditon, we can't discard the bytes after `i`
				// Then our next line's readUntilCondition will miss that starting bytes and fail
				t.nonAccumulatedBuffer = readBytes[i+1:]

				return accumulatedReadBytes, nil
			}
		}
	}

	return accumulatedReadBytes, ErrConditionNotMet
}

func (t *ConditionReader) ReadUntilTimeout(timeout time.Duration) ([]byte, error) {
	alwaysFalseCondition := func([]byte) bool {
		return false
	}

	data, err := t.ReadUntilConditionOrTimeout(alwaysFalseCondition, timeout)

	// We expect that the condition is never met, so let's return nil as the error
	if errors.Is(err, ErrConditionNotMet) {
		return data, nil
	}

	return data, err
}
