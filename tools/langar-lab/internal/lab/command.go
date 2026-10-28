package lab

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"time"
)

type Executor interface {
	Run(context.Context, string, []string, io.Reader, io.Writer, io.Writer) (int, error)
}

type OSExecutor struct{}

func (OSExecutor) Run(ctx context.Context, name string, args []string, in io.Reader, out, errOut io.Writer) (int, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = in, out, errOut
	// The Linux worker handles interruption and kills its own build/test groups.
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return interruptProcess(cmd.Process)
	}
	cmd.WaitDelay = 3 * time.Second
	err := cmd.Run()
	if err == nil {
		return 0, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode(), nil
	}
	return -1, err
}
