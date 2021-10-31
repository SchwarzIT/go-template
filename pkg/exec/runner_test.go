package exec_test

import (
	"os/exec"
	"runtime"
	"testing"

	ownexec "github.com/schwarzit/go-template/pkg/exec"
	"github.com/stretchr/testify/require"
)

func Test_execCmdRunner_Run(t *testing.T) {
	t.Run("error if command not found", func(t *testing.T) {
		_, err := ownexec.NewExecCmdRunner().Run(exec.Command("does-not-exist"))
		var errWithStderr *ownexec.ErrWithStderr
		require.ErrorAs(t, err, &errWithStderr)
		require.ErrorIs(t, err, exec.ErrNotFound)
	})
	t.Run("returns command's stdout", func(t *testing.T) {
		output, err := ownexec.NewExecCmdRunner().Run(exec.Command("go", "version"))
		require.NoError(t, err)
		require.Contains(t, output, runtime.Version())
	})
}
