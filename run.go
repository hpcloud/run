package run

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"syscall"
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

	go tailReader(bufio.NewReader(stdout), lines, errors)
	go tailReader(bufio.NewReader(stderr), lines, errors)

	go func() {
		resultCh <- cmd.Wait()
	}()

	return
}

func tailReader(r *bufio.Reader, ch chan string, errCh chan error) {
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			// log.Printf("** %T | %#v | %v\n", err, err, err)
			
			// Apparently EOF and EBADF are the two errors we must
			// safely ignore when reading from a Stdout or Stderr pipe
			// of a exec.Command - because they seem to be returned by
			// io.Reader when the command exits. Other errors are
			// considered abnormal, so we return them to the caller
			// via errCh.
			if err == io.EOF {
				break
			}
			if e, ok := err.(*os.PathError); ok && e.Err == syscall.EBADF {
				break
			}
			// unknown error
			errCh <- err
			break
		}
		ch <- string(line)
	}
}
