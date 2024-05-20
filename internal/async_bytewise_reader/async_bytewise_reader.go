package async_bytewise_reader

import (
	"errors"
	"io"
	"time"
)

var ErrNoData = errors.New("no data available")

// Inspired by https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/
type AsyncBytewiseReader struct {
	// data is used to send data between the reader goroutine and ReadByteWithTimeout calls
	data chan byte

	// err is used to store errors occurred during reading.
	// They're returned on the next ReadByteWithTimeout call
	err error

	// reader is the underlying reader that this wrapper will read from
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

// ReadByte is the only function that this package exposes. It either reads a byte or returns ErrNoData.
func (c *AsyncBytewiseReader) ReadByte() (byte, error) {
	select {
	// The timeout is super low here, we're just trying to check if a byte is immediately available
	case <-time.After(1 * time.Millisecond):
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
