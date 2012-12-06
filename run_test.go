package run

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestSimple(t *testing.T) {
	cmd := exec.Command("ls", "-l")

	lines, errors, resultCh := Run(cmd)

	done := make(chan bool)
	go func() {
		for {
			select {
			case line := <-lines:
				fmt.Println("LINE: ", line)
			case err := <-errors:
				t.Fatal(err)
				break
			case <-done:
				break
			}
		}
	}()

	err := <-resultCh
	done <- true
	fmt.Printf("Command exit value: %s\n", err)
}
