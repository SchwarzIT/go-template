package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	targetDir := t.TempDir()
	outputDir := t.TempDir()
	outputFile := path.Join(outputDir, "test.go")

	dotfiles := []string{
		path.Join(targetDir, ".gitignore"),
		path.Join(targetDir, ".dockerignore"),
		path.Join(targetDir, ".github"),
		path.Join(targetDir, ".azure-pipelines.yml"),
		path.Join(targetDir, ".gitlab-ci.yml"),
	}
	for _, file := range dotfiles {
		_, err := os.Create(path.Join(file))
		assert.NoError(t, err)
	}

	err := run([]string{"-target", targetDir, "-o", outputFile})
	assert.NoError(t, err)

	fileContents, err := os.ReadFile(outputFile)
	assert.NoError(t, err)

	assert.Contains(t, string(fileContents), fmt.Sprintf("//go:embed %s\n", strings.Join(sortStrings(append(dotfiles, targetDir)), " ")))
}
