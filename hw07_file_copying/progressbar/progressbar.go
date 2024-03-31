package progressbar

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ProgressBar struct {
	offset int64
	limit  int64
	w      io.Writer
}

func printProgress(percent int64) string {
	return fmt.Sprintf(
		"\r[%s%s] %d%%",
		strings.Repeat("|", int(percent)),
		strings.Repeat(" ", 100-int(percent)),
		percent)
}

func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.offset += int64(n)
	percent := pb.offset * 100 / pb.limit
	fmt.Fprint(pb.w, printProgress(percent))
	return
}

func NewProgressBar(limit int64) *ProgressBar {
	return &ProgressBar{
		offset: 0,
		limit:  limit,
		w:      os.Stdout,
	}
}
