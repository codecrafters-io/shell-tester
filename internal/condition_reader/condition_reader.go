package condition_reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/async_bytewise_reader"
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
	bytewiseReader *async_bytewise_reader.AsyncBytewiseReader
}

func NewConditionReader(reader io.Reader) ConditionReader {
	return ConditionReader{
		bytewiseReader: async_bytewise_reader.New(bufio.NewReader(reader)),
	}
}

func (t *ConditionReader) ReadUntilCondition(condition func([]byte) bool) ([]byte, error) {
	return t.ReadUntilConditionOrTimeout(condition, 2000*time.Millisecond)
}

func (t *ConditionReader) ReadUntilConditionOrTimeout(condition func([]byte) bool, timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)
	var readBytes []byte

	for !time.Now().After(deadline) {
		readByte, err := t.bytewiseReader.ReadByte()
		if err != nil {
			if errors.Is(err, async_bytewise_reader.ErrNoData) {
				debugLog("condition_reader: No data available")

				// Since no data was available, let's avoid a busy loop
				time.Sleep(2 * time.Millisecond)

				continue
			} else {
				debugLog("condition_reader: Error while reading: %v", err)
				return readBytes, err
			}
		}

		debugLog("condition_reader: readByte: %q", string(readByte))
		readBytes = append(readBytes, readByte)

		// If the condition is met, we can return early. Else the loop runs again
		if condition(readBytes) {
			return readBytes, nil
		}
	}

	return readBytes, ErrConditionNotMet
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