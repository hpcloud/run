package run

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestSimple(t *testing.T) {
	cmd := exec.Command("ls", "-l")

	lines, errors, resultCh := Run(cmd)

	for line := range lines {
		fmt.Println("LINE: ", line)
	}
	select {
	case err := <-errors:
		t.Fatal(err)
	default:
	}

	err := <-resultCh
	fmt.Printf("Command exit value: %s\n", err)
}
