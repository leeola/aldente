package aldente

import (
	"io"
	"time"

	"github.com/leeola/errors"
)

// LogStore stores logs in basic key-value stores.
//
// It *is not expected to be threadsafe*. Parallel writes may overwrite
// previous data.
type LogStore interface {
	PageCount(key string) (int, error)
	Read(key string, page int) ([]LogLine, error)
	Write(key string, page int, lines []LogLine) error
}

type WriteFlusher interface {
	io.Writer
	Flush() error
}
type LogLines struct {
	Lines []LogLine `json:"logLines"`
}

type LogLine struct {
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`
}

// NewWriteCloser implements appending to "dumb" key-value stores.
//
// The key-value store does not need to support appending to use this
// function.
//
// First, it buffers incoming data until a line is encountered. For
// every line, a LogLine is created.
// LogLines are buffered until the configured pagesize is achieved,
// at which time it's written to the key-value store.
//
// This takes the responsibility of paginating off of the caller and
// provides a simple WriteCloser interface for it.
//
// IMPORTANT: Writes are buffered, so Flush() must be called when done
// to assure the log is written properly to the store.
func NewWriteAppender(s LogStore, key string, pageSize int) (WriteFlusher, error) {
	return nil, errors.New("not implemented")
}
