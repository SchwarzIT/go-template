package gotemplate_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/schwarzit/go-template/pkg/option"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

const (
	targetDirOptionName = "projectSlug"
	optionName          = "someOption"
)

func TestNewRepositoryOptions_Validate(t *testing.T) {
	t.Run("CWD does not exist", func(t *testing.T) {
		opts := gotemplate.NewRepositoryOptions{
			CWD: "random-dir-that-does-not-exist",
		}

		assert.Error(t, opts.Validate())
	})

	t.Run("CWD is not set", func(t *testing.T) {
		opts := gotemplate.NewRepositoryOptions{}

		assert.NoError(t, opts.Validate())
	})

	t.Run("CWD set to valid dir", func(t *testing.T) {
		opts := gotemplate.NewRepositoryOptions{
			CWD: t.TempDir(),
		}

		assert.NoError(t, opts.Validate())
	})
}

func TestGT_LoadConfigValuesFromFile(t *testing.T) {
	gt := gotemplate.GT{
		Configs: option.Configuration{
			Parameters: []option.Option{
				{
					Name:    optionName,
					Default: "theDefault",
				},
			},
		},
	}

	t.Run("reads values from file", func(t *testing.T) {
		dir := t.TempDir()
		testFile := path.Join(dir, "test.yml")
		optionValue := "someOtherValue"
		testFileContent := fmt.Sprintf(`---
parameters:
    %s: %s
`, optionName, optionValue)
		err := os.WriteFile(testFile, []byte(testFileContent), os.ModePerm)
		assert.NoError(t, err)

		values, err := gt.LoadConfigValuesFromFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: optionValue}, values)
	})

	t.Run("validates that parameters are not empty", func(t *testing.T) {
		dir := t.TempDir()
		testFile := path.Join(dir, "test.yml")
		testFileContent := fmt.Sprintf(`---
parameters:
    %s: ""
`, optionName)
		err := os.WriteFile(testFile, []byte(testFileContent), os.ModePerm)
		assert.NoError(t, err)

		_, err = gt.LoadConfigValuesFromFile(testFile)
		assert.ErrorIs(t, err, gotemplate.ErrParameterNotSet)
	})

	t.Run("validates regex if set", func(t *testing.T) {
		gt.Configs.Parameters[0].Regex = option.Regex{
			Pattern:     `[a-z1-9]+(-[a-z1-9]+)*$`,
			Description: "only lowercase letters and dashes",
		}

		dir := t.TempDir()
		testFile := path.Join(dir, "test.yml")
		testFileContent := fmt.Sprintf(`---
parameters:
    %s: "NOT_A_VALID_VALUE"
`, optionName)
		err := os.WriteFile(testFile, []byte(testFileContent), os.ModePerm)
		assert.NoError(t, err)

		_, err = gt.LoadConfigValuesFromFile(testFile)
		assert.ErrorIs(t, err, gotemplate.ErrMalformedInput)
	})
}

func TestGT_LoadConfigValuesInteractively(t *testing.T) {
	gt := gotemplate.GT{
		Streams: gotemplate.Streams{Out: &bytes.Buffer{}},
	}

	optionValue := "someValue with spaces"
	t.Run("reads values from stdin", func(t *testing.T) {
		// simulate writing the value to stdin
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n"))
		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: "theDefault",
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: optionValue}, values)
	})

	t.Run("checks regex if it is set and retry if no match", func(t *testing.T) {
		// simulate writing the value to stdin
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader("DOES_NOT_MATCH\n matches-the-regex\n"))
		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: "theDefault",
				Regex: option.Regex{
					Pattern:     `[a-z1-9]+(-[a-z1-9]+)*$`,
					Description: "only lowercase letters and dashes",
				},
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: "matches-the-regex"}, values)
		assert.Contains(t, out.String(), "WARNING")
		assert.Contains(t, out.String(), "only lowercase letters and dashes", "should include regex description in warning message")
	})

	t.Run("checks regex on defaults as well", func(t *testing.T) {
		// simulate writing the value to stdin
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader("\nmatches-the-regex"))
		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: "DOES_NOT_MATCH",
				Regex: option.Regex{
					Pattern:     `[a-z1-9]+(-[a-z1-9]+)*$`,
					Description: "only lowercase letters and dashes",
				},
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: "matches-the-regex"}, values)
		assert.Contains(t, out.String(), "WARNING")
	})

	t.Run("retries to get value on error", func(t *testing.T) {
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "not a bool\ntrue\n"))
		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: false,
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: true}, values)
		assert.Contains(t, out.String(), "WARNING")
	})

	t.Run("applies templates from earlier options and uses default if not set", func(t *testing.T) {
		templateOptionName := "templatedOption"
		templatedOptionDefault := fmt.Sprintf(`{{.%s}}-templated`, optionName)
		// simulate setting a value for first option and use default for next
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n\n"))
		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: "theDefault",
			},
			{
				Name:    templateOptionName,
				Default: templatedOptionDefault,
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: optionValue, templateOptionName: fmt.Sprintf("%s-templated", optionValue)}, values)
	})

	t.Run("does not display options that have dependencies that are not met", func(t *testing.T) {
		dependentOptionName := "dependentOption"
		// simulate accepting the defaults
		gt.InScanner = bufio.NewScanner(strings.NewReader("\n\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: false,
			},
			{
				Name:      dependentOptionName,
				Default:   false,
				DependsOn: []string{optionName},
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, len(values), 1)
		assert.Contains(t, out.String(), optionName)
		assert.NotContains(t, out.String(), dependentOptionName)
	})

	t.Run("parses non string values", func(t *testing.T) {
		intOptionName := "intOption"
		// simulate accepting the defaults
		gt.InScanner = bufio.NewScanner(strings.NewReader("false\n4\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Configs.Parameters = []option.Option{
			{
				Name:    optionName,
				Default: true,
			},
			{
				Name:    intOptionName,
				Default: 3,
			},
		}

		values, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(values))
		assert.Equal(t, false, values[optionName])
		assert.Equal(t, 4, values[intOptionName])
	})
}

func TestGT_InitNewProject(t *testing.T) {
	// initialize template.FuncMap
	gt := gotemplate.New()
	gt.Streams.Out = &bytes.Buffer{}

	testValuesBytes, err := os.ReadFile("./testdata/values.yml")
	assert.NoError(t, err)

	fileStruct := gotemplate.FileStruct{}
	err = yaml.Unmarshal(testValuesBytes, &fileStruct)
	assert.NoError(t, err)

	opts := &gotemplate.NewRepositoryOptions{ConfigValues: fileStruct.Parameters}
	t.Run("generates folder in target dir and initializes it with go.mod and .git", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.CWD = tmpDir

		err = gt.InitNewProject(opts)
		assert.NoError(t, err)

		_, err = os.Stat(path.Join(getTargetDir(tmpDir, opts), ".git"))
		assert.NoError(t, err)

		_, err = os.Stat(path.Join(getTargetDir(tmpDir, opts), "go.mod"))
		assert.NoError(t, err)
	})

	t.Run("all templates should be resolved (in files and fileNames)", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.CWD = tmpDir

		err := gt.InitNewProject(opts)
		assert.NoError(t, err)

		err = filepath.WalkDir(getTargetDir(tmpDir, opts), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.Contains(path, "<no value>") {
				return fmt.Errorf("found a leftover template variable in %s", path)
			}

			if d.IsDir() || strings.Contains(path, ".git") {
				return nil
			}

			fileBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if strings.Contains(string(fileBytes), "<no value>") {
				return fmt.Errorf("found a leftover template variable in %s", path)
			}

			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("error if target dir already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.CWD = tmpDir

		err := os.MkdirAll(getTargetDir(tmpDir, opts), os.ModePerm)
		assert.NoError(t, err)

		err = gt.InitNewProject(opts)
		assert.Error(t, err)
	})

	t.Run("removes all files on error", func(t *testing.T) {
		tmpDir := t.TempDir()
		// force error with empty values
		err = gt.InitNewProject(
			&gotemplate.NewRepositoryOptions{CWD: tmpDir, ConfigValues: map[string]interface{}{
				targetDirOptionName: "testingDir",
			}},
		)
		assert.Error(t, err)

		_, err := os.Stat(getTargetDir(tmpDir, opts))
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("files for integrations are properly deleted or added", func(t *testing.T) {
		for _, opt := range gt.Configs.Parameters {
			if _, ok := opt.Default.(bool); !ok {
				continue
			}

			for _, enabled := range []bool{true, false} {
				t.Run(fmt.Sprintf("%s: %t", opt.Name, enabled), func(t *testing.T) {
					tmpDir := t.TempDir()
					fileStruct.Parameters[opt.Name] = enabled

					opts := &gotemplate.NewRepositoryOptions{CWD: tmpDir, ConfigValues: fileStruct.Parameters}
					err := gt.InitNewProject(opts)
					assert.NoError(t, err)

					for _, file := range opt.Files.Add {
						_, err := os.Stat(path.Join(getTargetDir(tmpDir, opts), file))
						if enabled {
							assert.NoErrorf(t, err, "file %s should exist", file)
						} else {
							assert.ErrorIsf(t, err, os.ErrNotExist, "file %s should be gone", file)
						}
					}

					for _, file := range opt.Files.Remove {
						_, err := os.Stat(path.Join(getTargetDir(tmpDir, opts), file))
						if enabled {
							assert.ErrorIsf(t, err, os.ErrNotExist, "file %s should be gone", file)
						} else {
							assert.NoErrorf(t, err, "file %s should exist", file)
						}
					}
				})
			}
		}
	})
}

func getTargetDir(dir string, opts *gotemplate.NewRepositoryOptions) string {
	return path.Join(dir, opts.ConfigValues[targetDirOptionName].(string))
}
