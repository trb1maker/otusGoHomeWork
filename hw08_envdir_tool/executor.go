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
	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
		}
		// По идее здесь бы добавить переменные в окружение, то тогда нет смысла
		// передавать их в командную RunCmd в качестве аргумента.
	}

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
