package gocli_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/schwarzit/go-template/pkg/gocli"
	"github.com/stretchr/testify/require"
)

func Test_Semver(t *testing.T) {
	version, err := gocli.Semver()
	require.NoError(t, err)
	// check that the version this test was build with matches the go version provided by goexec.Semver.
	runtimeVersion := semver.MustParse(strings.TrimPrefix(runtime.Version(), "go"))
	require.True(t, version.Equal(runtimeVersion))
}
