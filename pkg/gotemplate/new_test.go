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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/schwarzit/go-template/pkg/gotemplate"
)

var (
	errFoundLeftoverTemplateVar = errors.New("Found a leftover template variable in")
)

const (
	targetDirOptionName = "projectSlug"
	optionName          = "someOption"
)

func TestNewRepositoryOptions_Validate(t *testing.T) {
	t.Run("OutputDir does not exist", func(t *testing.T) {
		opts := gotemplate.NewRepositoryOptions{
			OutputDir: "random-dir-that-does-not-exist",
		}

		require.Error(t, opts.Validate())
	})

	t.Run("OutputDir is not set", func(t *testing.T) {
		opts := gotemplate.NewRepositoryOptions{}

		require.NoError(t, opts.Validate())
	})

	t.Run("OutputDir set to valid dir", func(t *testing.T) {
		opts := gotemplate.NewRepositoryOptions{
			OutputDir: t.TempDir(),
		}

		require.NoError(t, opts.Validate())
	})
}

func TestGT_LoadConfigValuesFromFile(t *testing.T) {
	gt := gotemplate.GT{
		Options: &gotemplate.Options{
			Base: []gotemplate.Option{
				gotemplate.NewOption(optionName, "description", gotemplate.StaticValue("theDefault")),
			},
		},
	}

	t.Run("reads values (base and extensions) from file", func(t *testing.T) {
		optionValue := "someValue"
		categoryName := "someCategory"
		categoryOptionName := "someCategoryOptionName"

		gt.Options.Extensions = []gotemplate.Category{
			{
				Name: categoryName,
				Options: []gotemplate.Option{
					gotemplate.NewOption(
						categoryOptionName,
						"description",
						gotemplate.StaticValue(false),
					),
				},
			},
		}

		optionValues, err := loadValueFromTestFile(t, &gt, fmt.Sprintf(`---
base:
    %s: %s
extensions:
    %s:
        %s: true`, optionName, optionValue, categoryName, categoryOptionName))

		require.NoError(t, err)
		require.Equal(
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

	t.Run("validates that base parameters are not empty", func(t *testing.T) {
		_, err := loadValueFromTestFile(t, &gt, fmt.Sprintf(`---
base:
    %s: ""`, optionName))

		require.ErrorIs(t, err, gotemplate.ErrParameterNotSet)
	})

	t.Run("validates validator if set", func(t *testing.T) {
		gt.Options.Base[0] = gotemplate.NewOption(
			optionName,
			"description",
			gotemplate.StaticValue("theDefault"),
			gotemplate.WithValidator(gotemplate.RegexValidator(
				`[a-z1-9]+(-[a-z1-9]+)*$`,
				"only lowercase letters and dashes",
			)),
		)

		_, err := loadValueFromTestFile(t, &gt, fmt.Sprintf(`---
base:
    %s: "NOT_A_VALID_VALUE"`, optionName))

		require.ErrorIs(t, err, gotemplate.ErrMalformedInput)
	})

	t.Run("sets default values for extensions", func(t *testing.T) {
		gt.Options = &gotemplate.Options{
			Extensions: []gotemplate.Category{
				{
					Name: "test",
					Options: []gotemplate.Option{
						gotemplate.NewOption(
							"string",
							"description",
							gotemplate.StaticValue("default"),
						),
					},
				},
			},
		}

		optionValues, err := loadValueFromTestFile(t, &gt, "")
		require.NoError(t, err)
		require.Equal(t, &gotemplate.OptionValues{
			Extensions: map[string]gotemplate.OptionNameToValue{
				"test": {
					"string": "default",
				},
			},
		}, optionValues)
	})

	t.Run("supports int, string, bool", func(t *testing.T) {
		gt.Options = &gotemplate.Options{
			Base: []gotemplate.Option{
				gotemplate.NewOption(
					"int",
					"description",
					gotemplate.StaticValue(2),
				),
				gotemplate.NewOption(
					"string",
					"description",
					gotemplate.StaticValue("string"),
				),
				gotemplate.NewOption(
					"bool",
					"description",
					gotemplate.StaticValue(false),
				),
			},
		}

		optionValues, err := loadValueFromTestFile(t, &gt, `---
base:
    int: 2
    string: "test"
    bool: true
`)

		require.NoError(t, err)
		require.Equal(t, &gotemplate.OptionValues{
			Base: gotemplate.OptionNameToValue{
				"int":    2,
				"string": "test",
				"bool":   true,
			},
		}, optionValues)
	})

	t.Run("error on type mismatch", func(t *testing.T) {
		gt.Options.Base[0] = gotemplate.NewOption(
			optionName,
			"description",
			gotemplate.StaticValue(false),
		)

		_, err := loadValueFromTestFile(t, &gt, fmt.Sprintf(`---
base:
    %s: "not a bool"`, optionName))

		var errTypeMismatch *gotemplate.ErrTypeMismatch
		require.ErrorAs(t, err, &errTypeMismatch)
	})

	t.Run("error if option is set but shouldDisplay returns false", func(t *testing.T) {
		gt.Options = &gotemplate.Options{
			Extensions: []gotemplate.Category{
				{
					Name: "test",
					Options: []gotemplate.Option{
						gotemplate.NewOption(
							"option",
							"description",
							gotemplate.StaticValue(false),
							gotemplate.WithShouldDisplay(gotemplate.BoolValue(false)),
						),
					},
				},
			},
		}

		_, err := loadValueFromTestFile(t, &gt, `---
extensions:
    test:
        option: true`)

		require.ErrorIs(t, err, gotemplate.ErrParameterSet)
	})

	t.Run("no error if option is set to default value shouldDisplay returns false", func(t *testing.T) {
		gt.Options = &gotemplate.Options{
			Extensions: []gotemplate.Category{
				{
					Name: "test",
					Options: []gotemplate.Option{
						gotemplate.NewOption(
							"option",
							"description",
							gotemplate.StaticValue(true),
							gotemplate.WithShouldDisplay(gotemplate.BoolValue(false)),
						),
					},
				},
			},
		}

		_, err := loadValueFromTestFile(t, &gt, `---
extensions:
    test:
        option: true`)
		require.NoError(t, err)
	})
}

func loadValueFromTestFile(t *testing.T, gt *gotemplate.GT, contents string) (*gotemplate.OptionValues, error) {
	dir := t.TempDir()
	testFile := path.Join(dir, "test.yml")
	err := os.WriteFile(testFile, []byte(contents), os.ModePerm)
	require.NoError(t, err)

	return gt.LoadConfigValuesFromFile(testFile)
}

func TestGT_LoadConfigValuesInteractively(t *testing.T) {
	gt := gotemplate.GT{
		Streams: gotemplate.Streams{Out: &bytes.Buffer{}},
		Options: &gotemplate.Options{},
	}

	optionValue := "someValue with spaces"
	t.Run("reads values (base and extensions) from stdin", func(t *testing.T) {
		optionValue := "someOtherValue"
		categoryName := "grpc"
		categoryOptionName := "base"
		out := &bytes.Buffer{}

		// simulate writing the value to stdin
		gt.InScanner = bufio.NewScanner(strings.NewReader(fmt.Sprintf("%s\n true\n", optionValue)))
		gt.Out = out
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				"description",
				gotemplate.StaticValue("theDefault"),
			),
		}
		gt.Options.Extensions = []gotemplate.Category{
			{
				Name: categoryName,
				Options: []gotemplate.Option{
					gotemplate.NewOption(
						categoryOptionName,
						"description",
						gotemplate.StaticValue(false),
					),
				},
			},
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		require.Equal(
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
		require.Contains(t, out.String(), "CATEGORY")
	})

	t.Run("checks regex if it is set and retry if no match", func(t *testing.T) {
		// simulate writing the value to stdin
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader("DOES_NOT_MATCH\n matches-the-regex\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				"description",
				gotemplate.StaticValue("DOES_NOT_MATCH"),
				gotemplate.WithValidator(gotemplate.RegexValidator(
					`[a-z1-9]+(-[a-z1-9]+)*$`,
					"only lowercase letters and dashes",
				)),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		require.Equal(t, gotemplate.OptionNameToValue{optionName: "matches-the-regex"}, optionValues.Base)
		require.Contains(t, out.String(), "WARNING")
		require.Contains(t, out.String(), "invalid pattern", "should include regex description in warning message")
	})

	t.Run("checks regex on defaults as well", func(t *testing.T) {
		// simulate writing the value to stdin
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader("\nmatches-the-regex"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				"description",
				gotemplate.StaticValue("DOES_NOT_MATCH"),
				gotemplate.WithValidator(gotemplate.RegexValidator(
					`[a-z1-9]+(-[a-z1-9]+)*$`,
					"only lowercase letters and dashes",
				)),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		require.Equal(t, gotemplate.OptionNameToValue{optionName: "matches-the-regex"}, optionValues.Base)
		require.Contains(t, out.String(), "WARNING")
	})

	t.Run("retries to get value on error", func(t *testing.T) {
		out := &bytes.Buffer{}
		gt.Err = out
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "not a bool\ntrue\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(optionName, "description", gotemplate.StaticValue(false)),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		require.Equal(t, gotemplate.OptionNameToValue{optionName: true}, optionValues.Base)
		require.Contains(t, out.String(), "WARNING")
	})

	t.Run("renders dynamic values correctly", func(t *testing.T) {
		templateOptionName := "templatedOption"
		// simulate setting a value for first option and use default for next
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n\n"))
		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				"description",
				gotemplate.StaticValue("theDefault"),
			),
			gotemplate.NewOption(
				templateOptionName,
				"description",
				gotemplate.DynamicValue(func(vals *gotemplate.OptionValues) interface{} {
					return vals.Base[optionName].(string) + "-templated"
				}),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		require.Equal(
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
				"description",
				gotemplate.StaticValue(false),
				gotemplate.WithShouldDisplay(gotemplate.BoolValue(false)),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		// the default value should be used in the optionValues in case there are dependent options
		require.Equal(t, len(optionValues.Base), 1)
		require.NotContains(t, out.String(), dependentOptionName)
	})

	t.Run("parses non string values", func(t *testing.T) {
		intOptionName := "intOption"
		gt.InScanner = bufio.NewScanner(strings.NewReader("false\n4\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				"description",
				gotemplate.StaticValue(true),
			),
			gotemplate.NewOption(
				intOptionName,
				"description",
				gotemplate.StaticValue(3),
			),
		}

		optionValues, err := gt.LoadConfigValuesInteractively()
		require.NoError(t, err)
		require.Equal(t, 2, len(optionValues.Base))
		require.Equal(t, false, optionValues.Base[optionName])
		require.Equal(t, 4, optionValues.Base[intOptionName])
	})
	t.Run("panics if default type is not supported", func(t *testing.T) {
		gt.InScanner = bufio.NewScanner(strings.NewReader("3.0\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Options.Base = []gotemplate.Option{
			gotemplate.NewOption(
				optionName,
				"description",
				// currently float is not supported
				gotemplate.StaticValue(2.0),
			),
		}

		require.PanicsWithValue(t, "unsupported type", func() {
			_, err := gt.LoadConfigValuesInteractively()
			if err != nil {
				t.Error("Error while loading config: %w", err)
			}
		})
	})
}

func TestGT_InitNewProject(t *testing.T) {
	// initialize template.FuncMap
	gt := gotemplate.New()
	gt.Streams.Out = &bytes.Buffer{}

	testValuesBytes, err := os.ReadFile("./testdata/values.yml")
	require.NoError(t, err)

	var optionValues gotemplate.OptionValues
	err = yaml.Unmarshal(testValuesBytes, &optionValues)
	require.NoError(t, err)

	opts := &gotemplate.NewRepositoryOptions{OptionValues: &optionValues}
	t.Run("generates folder in target dir and initializes it with go.mod and .git", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.OutputDir = tmpDir

		err = gt.InitNewProject(opts)
		require.NoError(t, err)

		_, err = os.Stat(path.Join(getTargetDir(tmpDir, opts), ".git"))
		require.NoError(t, err)

		_, err = os.Stat(path.Join(getTargetDir(tmpDir, opts), "go.mod"))
		require.NoError(t, err)
	})

	t.Run("copies hidden files (e.g. .gitignore)", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.OutputDir = tmpDir

		err = gt.InitNewProject(opts)
		require.NoError(t, err)

		testItems := []string{".gitignore", "pkg", "internal", ".golangci.yml"}
		for _, item := range testItems {
			_, err = os.Stat(path.Join(getTargetDir(tmpDir, opts), item))
			require.NoError(t, err)
		}
	})

	t.Run("all templates should be resolved (in files and fileNames)", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.OutputDir = tmpDir

		err := gt.InitNewProject(opts)
		require.NoError(t, err)

		err = filepath.WalkDir(getTargetDir(tmpDir, opts), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.Contains(path, "<no value>") {
				return fmt.Errorf("%w: %s", errFoundLeftoverTemplateVar, path)
			}

			if d.IsDir() || strings.Contains(path, ".git") {
				return nil
			}

			fileBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if strings.Contains(string(fileBytes), "<no value>") {
				return fmt.Errorf("%w: %s", errFoundLeftoverTemplateVar, path)
			}

			return nil
		})
		require.NoError(t, err)
	})

	t.Run("error if target dir already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.OutputDir = tmpDir

		err := os.MkdirAll(getTargetDir(tmpDir, opts), os.ModePerm)
		require.NoError(t, err)

		err = gt.InitNewProject(opts)
		require.Error(t, err)
	})

	t.Run("removes all files on error", func(t *testing.T) {
		tmpDir := t.TempDir()
		// force error with empty values
		err = gt.InitNewProject(
			&gotemplate.NewRepositoryOptions{
				OutputDir: tmpDir,
				OptionValues: &gotemplate.OptionValues{
					Base: gotemplate.OptionNameToValue{
						targetDirOptionName: "testingDir",
					},
				}},
		)
		require.Error(t, err)

		_, err := os.Stat(getTargetDir(tmpDir, opts))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("postHook not executed if value not set", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.OutputDir = tmpDir

		postHookTriggered := false
		gt.Options.Base = append(gt.Options.Base, gotemplate.NewOption(
			"testOption",
			"description",
			gotemplate.StaticValue(true),
			gotemplate.WithPosthook(func(value interface{}, optionValues *gotemplate.OptionValues, targetDir string) error {
				postHookTriggered = true
				return nil
			}),
		))

		err := gt.InitNewProject(opts)
		require.NoError(t, err)
		require.False(t, postHookTriggered, "postHook should not be triggered")
	})

	t.Run("postHook is executed if value is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts.OutputDir = tmpDir
		opts.OptionValues.Base["testOption"] = true

		postHookTriggered := false
		gt.Options.Base = append(gt.Options.Base, gotemplate.NewOption(
			"testOption",
			"description",
			gotemplate.StaticValue(false),
			gotemplate.WithPosthook(func(value interface{}, optionValues *gotemplate.OptionValues, targetDir string) error {
				postHookTriggered = true
				return nil
			}),
		))

		err := gt.InitNewProject(opts)
		require.NoError(t, err)
		require.True(t, postHookTriggered, "postHook should be triggered")
	})
}

func getTargetDir(dir string, opts *gotemplate.NewRepositoryOptions) string {
	return path.Join(dir, opts.OptionValues.Base[targetDirOptionName].(string))
}
