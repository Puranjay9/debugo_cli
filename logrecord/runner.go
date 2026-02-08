package logrecord

import (
	"context"
	"fmt"
	"os/exec"
)

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	traceChan := make(chan StackTrace, 100)
	processDone := make(chan struct{})

	go ScanStream(stdout, false, traceChan)
	go ScanStream(stderr, true, traceChan)

	go func() {
		cmd.Wait()
		close(processDone)
	}()

	PromptAndCapture(traceChan, processDone, cancel)

	return nil
}
