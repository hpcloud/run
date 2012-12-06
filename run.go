package run

import (
	"bufio"
	"io"
	"os/exec"
)

// Run runs the command and returns a channel of output lines, errors and result of cmd.Wait
func Run(cmd *exec.Cmd) (chan string, chan error, chan error) {
	lines := make(chan string)
	errors := make(chan error, 1)
	resultCh := make(chan error)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errors <- err
		return lines, errors, resultCh
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errors <- err
		return lines, errors, resultCh
	}

	err = cmd.Start()
	if err != nil {
		errors <- err
		return lines, errors, resultCh
	}

	go tailReader(bufio.NewReader(stdout), lines, errors)
	go tailReader(bufio.NewReader(stderr), lines, errors)

	go func() {
		resultCh <- cmd.Wait()
	}()

	return lines, errors, resultCh
}

func tailReader(r *bufio.Reader, ch chan string, errCh chan error) {
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err != io.EOF {
				errCh <- err
			}
			break
		}
		ch <- string(line)
	}
}
