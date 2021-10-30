package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

var (
	_ CmdRunner = (*execCmdRunner)(nil)
	_ CmdRunner = CmdRunnerFunc(nil)
)

// CmdRunner is an interface to safely abstract exec.Cmd calls.
type CmdRunner interface {
	Run(cmd *exec.Cmd) (string, error)
}

type CmdRunnerFunc func(cmd *exec.Cmd) (string, error)

func (f CmdRunnerFunc) Run(cmd *exec.Cmd) (string, error) {
	return f(cmd)
}

type execCmdRunner struct{}

func (r *execCmdRunner) Run(cmd *exec.Cmd) (string, error) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr

	if err := cmd.Run(); err != nil {
		return "", &ErrWithStderr{
			Wrapped: err,
			Args:    cmd.Args,
			StdErr:  stderr.Bytes(),
		}
	}

	return stdout.String(), nil
}

func NewExecCmdRunner() CmdRunner {
	return &execCmdRunner{}
}

type ErrWithStderr struct {
	Wrapped error
	StdErr  []byte
	Args    []string
}

func (e *ErrWithStderr) Error() string {
	if len(e.StdErr) > 0 {
		return fmt.Sprintf("failed running `%s`, %q: %s", strings.Join(e.Args, " "), e.StdErr, e.Wrapped.Error())
	}

	return fmt.Sprintf("failed running `%s`, make sure %s is available: %s", strings.Join(e.Args, " "), e.Args[0], e.Wrapped.Error())
}

func (e *ErrWithStderr) Unwrap() error {
	return e.Wrapped
}
