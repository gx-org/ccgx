// Package exec encapsulates go stdlib exec package with debug information.
package exec

import (
	"ccgx/internal/cmd/debug"
	"log"
	"os/exec"
	"strings"
)

// Command creates a new command.
func Command(name string, args ...string) *exec.Cmd {
	if debug.Debug {
		cmdS := append([]string{"DEBUG", "go"}, args...)
		log.Println(strings.Join(cmdS, " "))
	}
	return exec.Command(name, args...)
}
