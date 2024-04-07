package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
	myProgressBar "github.com/trb1maker/otus_golang_home_work/hw07_file_copying/progressbar"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSourceIsDir           = errors.New("source is directory")
	ErrDestinationIsDir      = errors.New("destination is directory")
	ErrEqualFiles            = errors.New("files are equal")
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

	// До открытия файла проверяю, что файл не является директорией
	// Настраиваю limit в зависимости от размера исходного файла
	srcStat, err := os.Stat(from)
	if err != nil {
		return err
	}
	if srcStat.IsDir() {
		return ErrSourceIsDir
	}

	if offset > srcStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > srcStat.Size()-offset {
		limit = srcStat.Size() - offset
	}

	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	// До открытия / создания файла проверяю, что файл не является директорией
	// и src и dst не ссылаются на один и тот же файл
	dstStat, err := os.Stat(to)
	if err == nil {
		if dstStat.IsDir() {
			return ErrDestinationIsDir
		}
		if os.SameFile(srcStat, dstStat) {
			return ErrEqualFiles
		}
	} else if os.IsExist(err) {
		// Если при проверке файла возникла ошибка, значит, dst либо не существует,
		// либо есть какие-то другие проблемы с доступом к нему. В первом случае
		// создаю файл, во втором - возвращаю ошибку
		return err
	}

	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	return copyWithOffset(src, dst, offset, limit)
}

func copyWithOffset(r io.ReaderAt, w io.Writer, offset int64, limit int64) error {
	var pb io.Writer
	if useMyProgressBar {
		pb = myProgressBar.NewProgressBar(limit)
	} else {
		pb = progressbar.DefaultBytes(limit)
	}
	r = io.NewSectionReader(r, offset, limit)
	_, err := io.Copy(io.MultiWriter(w, pb), r.(io.Reader))
	return err
}
