package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalln("pguard: invalid args <duration> <command> <command args>")
	}

	timeout := os.Args[1]
	cmdName := os.Args[2]
	cmdArgs := os.Args[3:]

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalln("pguard:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("pguard:", err)
	}
}
