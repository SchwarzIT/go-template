package exec

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// CommandGroup contains commands that are run one after another.
// As soon as one command fails the whole CommandGroup will be stopped and all other
// not yet executed commands are skipped.
type CommandGroup struct {
	// func run before any of the commands is executed
	PreRun func() error
	// Commands to run
	Commands []*exec.Cmd
	// TargetDir to run in
	TargetDir string
}

func (cg *CommandGroup) Run() error {
	return cg.RunWith(NewExecCmdRunner())
}

func (cg *CommandGroup) RunWith(runner CmdRunner) error {
	if len(cg.Commands) == 0 {
		return nil
	}

	if cg.PreRun != nil {
		if err := cg.PreRun(); err != nil {
			var skipsCmds []string
			for _, cmd := range cg.Commands {
				skipsCmds = append(skipsCmds, fmt.Sprintf("`%s`", strings.Join(cmd.Args, " ")))
			}

			return errors.Wrapf(err, "skipping %s", strings.Join(skipsCmds, ", "))
		}
	}

	for _, cmd := range cg.Commands {
		if cg.TargetDir != "" {
			cmd.Dir = cg.TargetDir
		}

		if _, err := runner.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}
