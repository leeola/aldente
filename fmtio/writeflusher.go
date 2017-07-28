package fmtio

import "io"

// WriteFlusher allows a buffered writer to be flushed.
type WriteFlusher interface {
	io.Writer
	Flush() error
}

type nopFlusher io.Writer

func NopFlusher(w io.Writer) WriteFlusher {
	return nopFlusher(w)
}

func (f nopFlusher) Flush() error {
	return nil
}
