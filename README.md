# run

Channel-friendly wrapper to run commands in Go.

## example

```Go
cmd := exec.Command("ls", "-l")

// reader
lines := make(chan string)
go func() {
	for line := range lines {
		fmt.Println("LINE: ", line)
	}
}()

// run `cmd`, writing the output to the lines channel
ret, err := Run(cmd, lines)

if err != nil {
	panic(err)
}
fmt.Printf("Command exit value: %s\n", ret)
```
