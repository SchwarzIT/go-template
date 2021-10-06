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
		Options: &gotemplate.Options{
			Base: []gotemplate.Option{
				gotemplate.NewOption(optionName, gotemplate.StringValue("description"), gotemplate.StaticValue("theDefault")),
			},
		},
	}

	t.Run("reads values (base and extensions) from file", func(t *testing.T) {
		optionValue := "someOtherValue"
		categoryName := "grpc"
		categoryOptionName := "base"

		gt.Options.Extensions = []gotemplate.Category{
			{
				Name: categoryName,
				Options: []gotemplate.Option{
					gotemplate.NewOption(
						categoryOptionName,
						gotemplate.StringValue("desc"),
						gotemplate.StaticValue(false),
					),
				},
			},
		}

		optionValues, err := loadValueFromTestFile(t, gt, fmt.Sprintf(`---
base:
    %s: %s
extensions:
    %s:
        %s: true`, optionName, optionValue, categoryName, categoryOptionName))

		assert.NoError(t, err)
		assert.Equal(
			t,
			&gotemplate.OptionValues{
				Base: gotemplate.OptionNameToValue{optionName: optionValue},
				Extensions: map[string]gotemplate.OptionNameToValue{
					categoryName: {
						categoryOptionName: true,
					},
				},
			},
			optionValues,
		)
	})

	t.Run("validates that parameters are not empty", func(t *testing.T) {
		_, err := loadValueFromTestFile(t, gt, fmt.Sprintf(`---
base:
    %s: ""`, optionName))

		assert.ErrorIs(t, err, gotemplate.ErrParameterNotSet)
	})

	t.Run("validates validator if set", func(t *testing.T) {
		gt.Options.Base[0] = gotemplate.NewOption(
			optionName,
			gotemplate.StringValue("description"),
			gotemplate.StaticValue("theDefault"),
			gotemplate.WithValidator(gotemplate.RegexValidator(
				`[a-z1-9]+(-[a-z1-9]+)*$`,
				"only lowercase letters and dashes",
			)),
		)

		_, err := loadValueFromTestFile(t, gt, fmt.Sprintf(`---
base:
    %s: "NOT_A_VALID_VALUE"`, optionName))

		assert.ErrorIs(t, err, gotemplate.ErrMalformedInput)
	})

	t.Run("error on type mismatch", func(t *testing.T) {
		gt.Options.Base[0] = gotemplate.NewOption(
			optionName,
			gotemplate.StringValue("description"),
			gotemplate.StaticValue(false),
		)

		_, err := loadValueFromTestFile(t, gt, fmt.Sprintf(`---
base:
    %s: "not a bool"`, optionName))

		var errTypeMismatch *gotemplate.ErrTypeMismatch
		assert.ErrorAs(t, err, &errTypeMismatch)
	})

	t.Run("error if option is set but shouldDisplay returns false", func(t *testing.T) {
		gt.Options.Base[0] = gotemplate.NewOption(
			optionName,
			gotemplate.StringValue("description"),
			gotemplate.StaticValue(true),
			gotemplate.WithShouldDisplay(gotemplate.BoolValue(false)),
		)

		_, err := loadValueFromTestFile(t, gt, fmt.Sprintf(`---
base:
    %s: true`, optionName))

		assert.ErrorIs(t, err, gotemplate.ErrParameterSet)
	})
}

func loadValueFromTestFile(t *testing.T, gt gotemplate.GT, contents string) (*gotemplate.OptionValues, error) {
	dir := t.TempDir()
	testFile := path.Join(dir, "test.yml")
	err := os.WriteFile(testFile, []byte(contents), os.ModePerm)
	assert.NoError(t, err)

	return gt.LoadConfigValuesFromFile(testFile)
}

func TestGT_LoadConfigValuesInteractively(t *testing.T) {
	gt := gotemplate.GT{
		Streams: gotemplate.Streams{Out: &bytes.Buffer{}},
		Options: &gotemplate.Options{},
	}

	optionValue := "someValue with spaces"
	t.Run("reads values from stdin", func(t *testing.T) {
		// simulate writing the value to stdin
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue("theDefault"),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, gotemplate.OptionNameToValue{optionName: optionValue}, optionValues.Base)
	})

	t.Run("checks regex if it is set and retry if no match", func(t *testing.T) {
		// simulate writing the value to stdin
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader("DOES_NOT_MATCH\n matches-the-regex\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue("DOES_NOT_MATCH"),
				gotemplate.WithValidator(gotemplate.RegexValidator(
					`[a-z1-9]+(-[a-z1-9]+)*$`,
					"only lowercase letters and dashes",
				)),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, gotemplate.OptionNameToValue{optionName: "matches-the-regex"}, optionValues.Base)
		assert.Contains(t, out.String(), "WARNING")
		assert.Contains(t, out.String(), "invalid pattern", "should include regex description in warning message")
	})

	t.Run("checks regex on defaults as well", func(t *testing.T) {
		// simulate writing the value to stdin
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader("\nmatches-the-regex"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue("DOES_NOT_MATCH"),
				gotemplate.WithValidator(gotemplate.RegexValidator(
					`[a-z1-9]+(-[a-z1-9]+)*$`,
					"only lowercase letters and dashes",
				)),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, gotemplate.OptionNameToValue{optionName: "matches-the-regex"}, optionValues.Base)
		assert.Contains(t, out.String(), "WARNING")
	})

	t.Run("retries to get value on error", func(t *testing.T) {
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "not a bool\ntrue\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(optionName, gotemplate.StringValue("desc"), gotemplate.StaticValue(false)),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, gotemplate.OptionNameToValue{optionName: true}, optionValues.Base)
		assert.Contains(t, out.String(), "WARNING")
	})

	t.Run("renders dynamic values correctly", func(t *testing.T) {
		templateOptionName := "templatedOption"
		// simulate setting a value for first option and use default for next
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue("theDefault"),
			),
			gotemplate.NewOption(
				templateOptionName,
				gotemplate.StringValue("desc"),
				gotemplate.DynamicValue(func(vals *gotemplate.OptionValues) interface{} {
					return vals.Base[optionName].(string) + "-templated"
				}),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(
			t,
			gotemplate.OptionNameToValue{
				optionName:         optionValue,
				templateOptionName: fmt.Sprintf("%s-templated", optionValue),
			},
			optionValues.Base,
		)
	})

	t.Run("does not display options that have shouldDisplay returning false", func(t *testing.T) {
		dependentOptionName := "dependentOption"
		// simulate accepting the defaults
		gt.InScanner = bufio.NewScanner(strings.NewReader("\n\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				dependentOptionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue(false),
				gotemplate.WithShouldDisplay(gotemplate.BoolValue(false)),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, len(optionValues.Base), 0)
		assert.NotContains(t, out.String(), dependentOptionName)
	})

	t.Run("parses non string values", func(t *testing.T) {
		intOptionName := "intOption"
		// simulate accepting the defaults
		gt.InScanner = bufio.NewScanner(strings.NewReader("false\n4\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue(true),
			),
			gotemplate.NewOption(
				intOptionName,
				gotemplate.StringValue("desc"),
				gotemplate.StaticValue(3),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(optionValues.Base))
		assert.Equal(t, false, optionValues.Base[optionName])
		assert.Equal(t, 4, optionValues.Base[intOptionName])
	})
}

func TestGT_InitNewProject(t *testing.T) {
	// initialize template.FuncMap
	gt := gotemplate.New()
	gt.Streams.Out = &bytes.Buffer{}

	testValuesBytes, err := os.ReadFile("./testdata/values.yml")
	assert.NoError(t, err)

	var optionValues gotemplate.OptionValues
	err = yaml.Unmarshal(testValuesBytes, &optionValues)
	assert.NoError(t, err)

	opts := &gotemplate.NewRepositoryOptions{OptionValues: &optionValues}
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
			&gotemplate.NewRepositoryOptions{
				CWD: tmpDir,
				OptionValues: &gotemplate.OptionValues{
					Base: gotemplate.OptionNameToValue{
						targetDirOptionName: "testingDir",
					},
				}},
		)
		assert.Error(t, err)

		_, err := os.Stat(getTargetDir(tmpDir, opts))
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("postHook not executed if value not set", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.CWD = tmpDir

		postHookTriggered := false
		gt.Options.Base = append(gt.Options.Base, gotemplate.NewOption(
			"testOption",
			gotemplate.StringValue("desc"),
			gotemplate.StaticValue(true),
			gotemplate.WithPosthook(func(value interface{}, optionValues *gotemplate.OptionValues, targetDir string) error {
				postHookTriggered = true
				return nil
			}),
		))

		err := gt.InitNewProject(opts)
		assert.NoError(t, err)
		assert.False(t, postHookTriggered, "postHook should not be triggered")
	})

	t.Run("postHook is executed if value is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.CWD = tmpDir
		opts.OptionValues.Base["testOption"] = true

		postHookTriggered := false
		gt.Options.Base = append(gt.Options.Base, gotemplate.NewOption(
			"testOption",
			gotemplate.StringValue("desc"),
			gotemplate.StaticValue(false),
			gotemplate.WithPosthook(func(value interface{}, optionValues *gotemplate.OptionValues, targetDir string) error {
				postHookTriggered = true
				return nil
			}),
		))

		err := gt.InitNewProject(opts)
		assert.NoError(t, err)
		assert.True(t, postHookTriggered, "postHook should be triggered")
	})
}

func getTargetDir(dir string, opts *gotemplate.NewRepositoryOptions) string {
	return path.Join(dir, opts.OptionValues.Base[targetDirOptionName].(string))
}
