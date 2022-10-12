package exec_test

import (
	"errors"
	"os/exec"
	"testing"

	ownexec "github.com/schwarzit/go-template/pkg/exec"
	"github.com/stretchr/testify/require"
)

var (
	errDummy = errors.New("dummy")
)

func Test_CommandGroup_RunWith(t *testing.T) {
	t.Run("no command is executed if prerun fails", func(t *testing.T) {
		cg := ownexec.CommandGroup{
			PreRun: func() error {
				return errDummy
			},
			Commands: []*exec.Cmd{
				exec.Command("anything"),
			},
		}

		anyCommandExecuted := false
		err := cg.RunWith(ownexec.CmdRunnerFunc(func(cmd *exec.Cmd) (string, error) {
			anyCommandExecuted = true
			return "", nil
		}))

		require.Error(t, err)
		require.False(t, anyCommandExecuted)
	})
	t.Run("all commands are executed if there's no error", func(t *testing.T) {
		cg := ownexec.CommandGroup{
			Commands: []*exec.Cmd{
				exec.Command("anything"),
				exec.Command("sthelse"),
			},
		}

		anythingExecuted, sthelseExecuted := false, false
		err := cg.RunWith(ownexec.CmdRunnerFunc(func(cmd *exec.Cmd) (string, error) {
			switch cmd.Path {
			case "anything":
				anythingExecuted = true
			case "sthelse":
				sthelseExecuted = true
			}
			return "", nil
		}))

		require.NoError(t, err)
		require.True(t, anythingExecuted)
		require.True(t, sthelseExecuted)
	})
	t.Run("no commands are executed after the first one fails", func(t *testing.T) {
		cg := ownexec.CommandGroup{
			Commands: []*exec.Cmd{
				exec.Command("anything"),
				exec.Command("sthelse"),
			},
		}

		sthelseExecuted := false
		err := cg.RunWith(ownexec.CmdRunnerFunc(func(cmd *exec.Cmd) (string, error) {
			switch cmd.Path {
			case "anything":
				return "", errDummy
			case "sthelse":
				sthelseExecuted = true
			}
			return "", nil
		}))

		require.Error(t, err)
		require.False(t, sthelseExecuted)
	})
}
