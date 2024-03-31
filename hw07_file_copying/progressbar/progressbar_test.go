package progressbar

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	rx = regexp.MustCompile(`\[X{10}\s{90}\]\s10%`)
)

func TestPrintProgress(t *testing.T) {
	printed := printProgress(10)

	require.True(t, rx.MatchString(printed))
}

func TestProgressBar(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	pb := &ProgressBar{w: buf, limit: 50}
	pb.Write([]byte{1, 2, 3, 4, 5})
	require.True(t, rx.MatchString(buf.String()))
}
