package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("pguard: invalid args")
	}

	_ = os.Args[1]
	cmdName := os.Args[2]
	cmdArgs := os.Args[3:]

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("pguard:", err)
	}
}
