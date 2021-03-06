package util

import (
	"io"

	"github.com/leeola/motley/fmtio"
)

type nopFlusher struct {
	io.Writer
}

func NopFlusher(w io.Writer) fmtio.WriteFlusher {
	return &nopFlusher{Writer: w}
}

func (nopFlusher) Flush() error {
	return nil
}
