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

// ConditionReader wraps an io.Reader and provides methods to read until a condition is met
type ConditionReader struct {
	asyncReader *async_reader.AsyncReader
}

func NewConditionReader(reader io.Reader) ConditionReader {
	return ConditionReader{
		asyncReader: async_reader.New(bufio.NewReader(reader)),
	}
}

func (t *ConditionReader) ReadUntilConditionOrTimeout(condition func() bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	if condition() {
		return nil
	}

	for !time.Now().After(deadline) {
		bytes, err := t.asyncReader.Read()
		if len(bytes) > 0 {
			fmt.Printf("Read bytes: %q\n", string(bytes))
		}
		if err != nil {
			if errors.Is(err, async_reader.ErrNoData) {
				// Since no data was available, let's avoid a busy loop
				time.Sleep(2 * time.Millisecond)
				continue
			} else {
				return err
			}
		}

		if condition() {
			return nil
		}
	}

	return ErrConditionNotMet
}
