package condition_reader

import (
	"bufio"
	"errors"
	"os"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/async_bytewise_reader"
)

var conditionFailedError = errors.New("condition failed")

type ConditionReader struct {
	bytewiseReader *async_bytewise_reader.AsyncBytewiseReader
}

func NewConditionReader(descriptor *os.File) ConditionReader {
	return ConditionReader{
		bytewiseReader: async_bytewise_reader.New(bufio.NewReader(descriptor)),
	}
}

func (t *ConditionReader) ReadUntilTimeout(timeout time.Duration) ([]byte, error) {
	alwaysFalseCondition := func([]byte) bool {
		return false
	}

	data, err := t.ReadUntilConditionWithTimeout(alwaysFalseCondition, timeout)

	// We expect that the condition is never met, so let's return nil as the error
	if err == conditionFailedError {
		return data, nil
	}

	return data, err
}

func (t *ConditionReader) ReadUntilCondition(condition func([]byte) bool) ([]byte, error) {
	return t.ReadUntilConditionWithTimeout(condition, 2000*time.Millisecond)
}

func (t *ConditionReader) ReadUntilConditionWithTimeout(condition func([]byte) bool, timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)
	readBytes := []byte{}

	for !time.Now().After(deadline) {
		readByte, err := t.bytewiseReader.ReadByteWithTimeout(2 * time.Millisecond)
		if err != nil {
			if err == async_bytewise_reader.ErrNoData {
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
		if condition(readBytes) {
			return readBytes, nil
		}
	}

	return readBytes, conditionFailedError
}
