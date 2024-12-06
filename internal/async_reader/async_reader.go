package async_reader

import (
	"errors"
	"io"
	"time"
)

var ErrNoData = errors.New("no data available")

// AsyncReader : Inspired by https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/
// We don't require a BytewiseReader anymore, but we still require a cancellable AsyncReader
type AsyncReader struct {
	// data is used to send data between the reader goroutine and ReadByte calls
	data chan []byte

	// err is used to store an error occurred during reading.
	// The error will only be returned on the next ReadByte call.
	err error

	// reader is the underlying reader that this wrapper will read from
	reader io.Reader

	// unreadBuffer is used to store unprocessed bytes
	unreadBuffer []byte
}

func New(reader io.Reader) *AsyncReader {
	asyncReader := &AsyncReader{
		reader: reader,
		data:   make(chan []byte),
	}

	// This goroutine will keep reading until an error or EOF
	go asyncReader.start()

	return asyncReader
}

// ReadBytes is the only function that this package exposes. It either reads a byte or returns ErrNoData.
func (r *AsyncReader) Read() ([]byte, error) {
	select {
	// We're checking whether a byte is immediately available, so the timeout can be super low
	case <-time.After(1 * time.Millisecond):
		return nil, ErrNoData
	case readBytes, ok := <-r.data:
		if !ok {
			return nil, r.err
		}

		if len(r.unreadBuffer) > 0 {
			readBytes = append(r.unreadBuffer, readBytes...)
			r.unreadBuffer = []byte{}
		}

		return readBytes, nil
	}
}

// Keeps reading forever until an error or EOF
func (r *AsyncReader) start() {
	for {
		buf := make([]byte, 1024)
		n, err := r.reader.Read(buf)

		if n > 0 {
			r.data <- buf[:n]
		}

		if err != nil {
			r.err = err
			close(r.data)
			return
		}
	}
}
