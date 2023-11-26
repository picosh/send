package utils

import (
	"io"
)

type ReadAndReaderAt interface {
	io.ReaderAt
	io.Reader
}

type ReaderAtCloser interface {
	io.ReaderAt
	io.ReadCloser
}

func NopReaderAtCloser(r ReadAndReaderAt) ReaderAtCloser {
	return nopReaderAtCloser{r}
}

type nopReaderAtCloser struct {
	ReadAndReaderAt
}

func (nopReaderAtCloser) Close() error { return nil }
