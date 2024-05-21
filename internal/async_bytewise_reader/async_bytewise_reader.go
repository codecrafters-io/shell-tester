package async_bytewise_reader

import (
	"errors"
	"io"
	"time"
)

var ErrNoData = errors.New("no data available")

// Inspired by https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/
type AsyncBytewiseReader struct {
	// data is used to send data between the reader goroutine and ReadByte calls
	data chan byte

	// err is used to store an error occurred during reading.
	// The error will only be returned on the next ReadByte call.
	err error

	// reader is the underlying reader that this wrapper will read from
	reader io.Reader
}

func New(reader io.Reader) *AsyncBytewiseReader {
	bytewiseReader := &AsyncBytewiseReader{
		reader: reader,
		data:   make(chan byte),
	}

	// This goroutine will keep reading until an error or EOF
	go bytewiseReader.start()

	return bytewiseReader
}

// ReadByte is the only function that this package exposes. It either reads a byte or returns ErrNoData.
func (r *AsyncBytewiseReader) ReadByte() (byte, error) {
	select {
	// The timeout is super low here, we're just trying to check if a byte is immediately available
	case <-time.After(1 * time.Millisecond):
		return 0, ErrNoData
	case readByte, ok := <-r.data:
		if !ok {
			return 0, r.err
		}

		return readByte, nil
	}
}

// Keeps reading forever until an error or EOF
func (r *AsyncBytewiseReader) start() {
	for {
		buf := make([]byte, 1024)
		n, err := r.reader.Read(buf)

		if n > 0 {
			for _, b := range buf[:n] {
				r.data <- b
			}
		}

		if err != nil {
			r.err = err
			close(r.data)
			return
		}
	}
}
