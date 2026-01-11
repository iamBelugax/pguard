package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	timeout  time.Duration
	graceful time.Duration
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	flag.DurationVar(&timeout, "timeout", -1, "Maximum runtime for the command (e.g. 10s, 1m). -1 means no timeout")
	flag.DurationVar(&graceful, "graceful", -1, "Grace period before force kill after timeout")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Println("pguard: no command specified")
		flag.Usage()
		os.Exit(1)
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	ctx, cancel := makeContext(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("pguard:", err)
	}
}

func makeContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithCancel(ctx)
}

func usage() {
	fmt.Fprintf(os.Stderr, `pguard - process guard with timeout

Usage:
  pguard [flags] <command> [command args...]

Flags:
`)
	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, `
Examples:
  pguard --timeout=10s --graceful=5s -- sleep 30
  pguard --timeout=1m command --foo bar
`)
}
