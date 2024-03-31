package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSourceIsDir           = errors.New("source is directory")
	ErrLimitIsNegative       = errors.New("limit is negative")
	ErrOffsetIsNegative      = errors.New("offset is negative")
)

func Copy(from string, to string, offset int64, limit int64) error {
	if offset < 0 {
		return ErrOffsetIsNegative
	}
	if limit < 0 {
		return ErrLimitIsNegative
	}

	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return ErrSourceIsDir
	}

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > stat.Size()-offset {
		limit = stat.Size() - offset
	}

	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	return copyWithOffset(src, dst, offset, limit)
}

func copyWithOffset(r io.ReaderAt, w io.Writer, offset int64, limit int64) error {
	pb := progressbar.DefaultBytes(limit)
	r = io.NewSectionReader(r, offset, limit)
	_, err := io.Copy(io.MultiWriter(w, pb), r.(io.Reader))
	return err
}
