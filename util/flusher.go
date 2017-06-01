package util

import (
	"io"

	"github.com/leeola/aldente/fmtio"
)

type nopFlusher struct {
	w io.Writer
}

func NopFlusher(w io.Writer) fmtio.WriteFlusher {
	return &nopFlusher{w: w}
}
