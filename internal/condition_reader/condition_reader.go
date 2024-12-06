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
	asyncReader *async_reader.AsyncReader
}

func NewConditionReader(reader io.Reader) ConditionReader {
	return ConditionReader{
		asyncReader: async_reader.New(bufio.NewReader(reader)),
	}
}

func (t *ConditionReader) ReadUntilCondition(condition func() bool) error {
	return t.ReadUntilConditionOrTimeout(condition, 2000*time.Millisecond)
}

func (t *ConditionReader) ReadUntilConditionOrTimeout(condition func() bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for !time.Now().After(deadline) {
		readBytes, err := t.asyncReader.Read()
		if err != nil {
			if errors.Is(err, async_reader.ErrNoData) {
				debugLog("condition_reader: No data available")

				// Since no data was available, let's avoid a busy loop
				time.Sleep(2 * time.Millisecond)

				continue
			} else {
				debugLog("condition_reader: Error while reading: %v", err)
				return err
			}
		}

		debugLog("condition_reader: readBytes: %q", string(readBytes))

		if condition() {
			return nil
		}
	}

	return ErrConditionNotMet
}

func (t *ConditionReader) ReadUntilTimeout(timeout time.Duration) error {
	alwaysFalseCondition := func() bool {
		return false
	}

	err := t.ReadUntilConditionOrTimeout(alwaysFalseCondition, timeout)

	// We expect that the condition is never met, so let's return nil as the error
	if errors.Is(err, ErrConditionNotMet) {
		return nil
	}

	return err
}
