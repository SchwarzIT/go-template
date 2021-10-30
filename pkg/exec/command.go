package exec

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrInvalidArgs = errors.New("invalid command args")

var _ Command = (*execCmd)(nil)

type Command interface {
	// Run executes the command and returns stdout as string as well as an error if the command failed.
	// The error is of type ErrWithStderr and can be type asserted to that to get the exact stderr output.
	Run() (string, error)
	// Args returns all the command's args including the executable name as the first item.
	Args() []string
}

type execCmd struct {
	*exec.Cmd
}

func (ec *execCmd) Run() (string, error) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	ec.Stdout, ec.Stderr = stdout, stderr

	if err := ec.Cmd.Run(); err != nil {
		return "", &ErrWithStderr{
			Wrapped: err,
			Args:    ec.Cmd.Args,
			StdErr:  stderr.Bytes(),
		}
	}

	return stdout.String(), nil
}

func (ec *execCmd) Args() []string {
	return ec.Cmd.Args
}

func NewExecCmd(args []string, opts ...NewExecCmdOption) Command {
	cmd := exec.Command(args[0], args[1:]...)
	for _, opt := range opts {
		opt(cmd)
	}

	return &execCmd{Cmd: cmd}
}

type NewExecCmdOption func(*exec.Cmd)

func WithTargetDir(targetDir string) NewExecCmdOption {
	return func(c *exec.Cmd) {
		c.Dir = targetDir
	}
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
