package run

import (
	"log"
	"bufio"
	"io"
	"os/exec"
	"sync"
)

// Run runs the command and returns a channel of output lines, errors and result of cmd.Wait
func Run(cmd *exec.Cmd) (
	lines chan string, errors chan error, resultCh chan error) {
	lines = make(chan string)
	errors = make(chan error, 1)
	resultCh = make(chan error)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errors <- err
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errors <- err
		return
	}

	err = cmd.Start()
	if err != nil {
		errors <- err
		return
	}

	var wg sync.WaitGroup
	go tailReader(bufio.NewReader(stdout), lines, errors, wg)
	go tailReader(bufio.NewReader(stderr), lines, errors, wg)

	go func() {
		wg.Wait()
		resultCh <- cmd.Wait()
	}()

	return
}

func tailReader(
	r *bufio.Reader, ch chan string, errCh chan error,
	wg sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Printf("** %T | %#v | %v\n", err, err, err)

			if err != io.EOF {
				errCh <- err
			}
			break
		}
		ch <- string(line)
	}
}
