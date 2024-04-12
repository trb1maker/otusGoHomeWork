package main

import (
	"log"
	"os"
)

func main() {
	envDir := os.Args[1]
	args := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		log.Printf("failed to read env: %v", err)
		return
	}

	returnCode := RunCmd(args, env)
	os.Exit(returnCode)
}
