package run

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestSimple(t *testing.T) {
	cmd := exec.Command("ls", "-l")

	lines := make(chan string)
	go func() {
		for line := range lines {
			fmt.Println("LINE: ", line)
		}
	}()
	ret, err := Run(cmd, lines)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Command exit value: %s\n", ret)
}
