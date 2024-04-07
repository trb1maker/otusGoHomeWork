package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	return runCmd(cmd, env, os.Stdout, os.Stderr)
}

func runCmd(cmd []string, env Environment, output io.Writer, errors io.Writer) (returnCode int) {
	envs := os.Environ()
	for name, value := range env {
		envs = append(envs, name+"="+value.Value)
	}
	app := &exec.Cmd{
		Path:   cmd[0],
		Args:   cmd,
		Env:    envs,
		Stdin:  os.Stdin,
		Stdout: output,
		Stderr: errors,
	}
	if err := app.Run(); err != nil {
		log.Printf("failed to run command: %v", err)
	}
	return app.ProcessState.ExitCode()
}
