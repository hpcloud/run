package run

import (
	"bufio"
	"io"
	"os/exec"
	"sync"
)

// Run runs the command and returns a channel of output lines, errors and result of cmd.Wait
func Run(cmd *exec.Cmd, lines chan string) (error, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)
	go tailReader(bufio.NewReader(stdout), lines, errCh, &wg)
	go tailReader(bufio.NewReader(stderr), lines, errCh, &wg)
	wg.Wait()
	select {
		case err := <-errCh:
		return nil, err
		default:
	}
	return cmd.Wait(), nil
}

func tailReader(
	r *bufio.Reader, ch chan string, errCh chan error,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			// log.Printf("** %T | %#v | %v\n", err, err, err)

			if err != io.EOF {
				errCh <- err
			}
			break
		}
		ch <- string(line)
	}
}
