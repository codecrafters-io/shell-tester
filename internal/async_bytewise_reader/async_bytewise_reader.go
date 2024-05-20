package async_bytewise_reader

import (
	"errors"
	"io"
	"time"
)

var ErrNoData = errors.New("no data available")

// Inspired by https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/
type AsyncBytewiseReader struct {
	data   chan byte
	err    error
	reader io.Reader
}

func New(r io.Reader) *AsyncBytewiseReader {
	c := &AsyncBytewiseReader{
		reader: r,
		data:   make(chan byte),
	}

	// This goroutine will keep reading until an error or EOF
	go c.start()

	return c
}

func (c *AsyncBytewiseReader) ReadByteWithTimeout(timeout time.Duration) (byte, error) {
	select {
	case <-time.After(timeout):
		return 0, ErrNoData
	case readByte, ok := <-c.data:
		if !ok {
			return 0, c.err
		}

		return readByte, nil
	}
}

// Keeps reading forever until an error or EOF
func (c *AsyncBytewiseReader) start() {
	for {
		buf := make([]byte, 1024)
		n, err := c.reader.Read(buf)

		if n > 0 {
			for _, b := range buf[:n] {
				c.data <- b
			}
		}

		if err != nil {
			c.err = err
			close(c.data)
			return
		}
	}
}
