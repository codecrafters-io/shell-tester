package condition_reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/async_reader"
)

type ConditionNotMetError struct {
	ExpectedCondition string
	Output            string
}

func (e *ConditionNotMetError) Error() string {
	return fmt.Sprintf("expected %s but condition was not met. Program output:\n%s", e.ExpectedCondition, e.Output)
}

// ConditionReader wraps an io.Reader and provides methods to read until a condition is met
type ConditionReader struct {
	asyncReader *async_reader.AsyncReader
	output      strings.Builder
}

func NewConditionReader(reader io.Reader) ConditionReader {
	return ConditionReader{
		asyncReader: async_reader.New(bufio.NewReader(reader)),
		output:      strings.Builder{},
	}
}

func (t *ConditionReader) ReadUntilConditionOrTimeout(condition func() bool, timeout time.Duration, expectedCondition string) error {
	deadline := time.Now().Add(timeout)

	if condition() {
		return nil
	}

	for !time.Now().After(deadline) {
		data, err := t.asyncReader.Read()
		if err != nil {
			if errors.Is(err, async_reader.ErrNoData) {
				// Since no data was available, let's avoid a busy loop
				time.Sleep(2 * time.Millisecond)
				continue
			} else {
				return err
			}
		}

		t.output.Write(data)

		if condition() {
			return nil
		}
	}

	return &ConditionNotMetError{
		ExpectedCondition: expectedCondition,
		Output:            t.output.String(),
	}
}
