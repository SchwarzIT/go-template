package gocli

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

var ErrMalformedGoVersionOutput = errors.New("malformed go version output")

func Semver() (*semver.Version, error) {
	stdout := &bytes.Buffer{}

	goVersion := exec.Command("go", "version")
	goVersion.Stdout = stdout

	err := goVersion.Run()
	if err != nil {
		return nil, errors.Wrap(err, "failed checking go version")
	}

	versionParts := strings.Split(stdout.String(), " ")
	if len(versionParts) != 4 { //nolint:gomnd // go version output has exactly 4 parts (e.g. "go version go1.17.2 darwin/amd64")
		return nil, errors.Wrap(ErrMalformedGoVersionOutput, stdout.String())
	}

	goSemverString := strings.TrimPrefix(versionParts[2], "go")
	goSemver, err := semver.NewVersion(goSemverString)
	if err != nil {
		return nil, errors.Wrap(ErrMalformedGoVersionOutput, stdout.String())
	}

	return goSemver, nil
}
