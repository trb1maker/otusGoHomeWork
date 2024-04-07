package main

import (
	"log/slog"
	"os"
)

func main() {
	envDir := os.Args[1]
	args := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		slog.Error("failed to read env", err)
		return
	}

	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
		}
		// По идее здесь бы добавить переменные в окружение, то тогда нет смысла
		// передавать их в командную RunCmd в качестве аргумента.
	}

	returnCode := RunCmd(args, env)
	os.Exit(returnCode)
}
