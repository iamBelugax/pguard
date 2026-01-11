package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

var (
	timeout  time.Duration
	graceful time.Duration
)

func main() {
	log.SetFlags(0)

	flag.DurationVar(&timeout, "timeout", -1, "Maximum runtime (e.g. 10s, 1m). -1 means no timeout.")
	flag.DurationVar(&graceful, "graceful", 5*time.Second, "Grace period before force kill.")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		log.Println("pguard: no command specified")
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := makeContext(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, flag.Arg(0), flag.Args()[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalln("pguard:", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- cmd.Wait()
	}()

	select {
	case err := <-doneCh:
		exit(err)

	case <-sigCh:
		log.Println("pguard: interrupt received")

	case <-ctx.Done():
		log.Println("pguard: timeout reached")
	}

	_ = cmd.Process.Signal(os.Interrupt)

	select {
	case err := <-doneCh:
		exit(err)

	case <-time.After(graceful):
		log.Println("pguard: force killing process")
		doneCh <- cmd.Process.Kill()
		exit(<-doneCh)
	}
}

func makeContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithCancel(ctx)
}

func exit(err error) {
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func usage() {
	fmt.Fprintf(os.Stderr, `pguard - process guard with timeout

Usage:
  pguard [flags] <command> [args...]

Flags:
`)
	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, `
Examples:
  pguard --timeout=10s --graceful=5s -- sleep 30
  pguard --timeout=1m command --foo bar
`)
}
