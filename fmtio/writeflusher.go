package fmtio

import "io"

// WriteFlusher allows a buffered writer to be flushed.
type WriteFlusher interface {
	io.Writer
	Flush() error
}
