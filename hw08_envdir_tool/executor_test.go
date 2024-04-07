package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	output := new(bytes.Buffer)
	errors := new(bytes.Buffer)
	envs := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
	}
	cmd := []string{"/bin/bash", "-c", `echo -e "$BAR,$EMPTY,$HELLO,${FOO}"`}
	want := "bar,,\"hello\",   foo\nwith new line\n"
	resultCode := runCmd(cmd, envs, output, errors)

	if errors.String() != "" {
		t.Errorf("unexpected error: %s", errors.String())
	}

	require.Equal(t, want, output.String())
	require.Equal(t, 0, resultCode)
}
