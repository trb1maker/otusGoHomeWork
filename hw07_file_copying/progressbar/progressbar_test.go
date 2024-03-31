package progressbar

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintProgress(t *testing.T) {
	tests := []struct {
		percent int64
		want    string
	}{
		{0, "\r[                                                                                                    ] 0%"},
		{10, "\r[XXXXXXXXXX                                                                                          ] 10%"},
		{100, "\r[XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX] 100%"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("percent=%d", tt.percent), func(t *testing.T) {
			got := printProgress(tt.percent)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestProgressBar(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	pb := &ProgressBar{w: buf, limit: 50}
	pb.Write([]byte{1, 2, 3, 4, 5})
	require.Equal(t, "\r[XXXXXXXXXX                                                                                          ] 10%", buf.String())
}
