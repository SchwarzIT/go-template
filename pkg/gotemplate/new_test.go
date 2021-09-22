package gotemplate_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/schwarzit/go-template/pkg/option"
	"github.com/stretchr/testify/assert"
)

func TestGT_LoadOptionToValueFromFile(t *testing.T) {
	optionName := "someOption"

	gt := gotemplate.GT{
		Options: []option.Option{
			{
				Name:    optionName,
				Default: "theDefault",
			},
		},
	}

	dir := t.TempDir()
	testFile := path.Join(dir, "test.yml")
	t.Run("reads values from file", func(t *testing.T) {
		optionValue := "someOtherValue"
		err := os.WriteFile(testFile, []byte(fmt.Sprintf(`%s: %s`, optionName, optionValue)), os.ModePerm)
		assert.NoError(t, err)

		values, err := gt.LoadOptionToValueFromFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: optionValue}, values)
	})
}

func TestGT_GetOptionToValueInteractively(t *testing.T) {
	optionName := "someOption"

	gt := gotemplate.GT{
		Streams: gotemplate.Streams{Out: &bytes.Buffer{}},
	}

	optionValue := "someValue with spaces"
	t.Run("reads values from file", func(t *testing.T) {
		// simulate writing the value to stdin
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n"))
		gt.Options = []option.Option{
			{
				Name:    optionName,
				Default: "theDefault",
			},
		}

		values, err := gt.GetOptionToValueInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: optionValue}, values)
	})

	t.Run("applies templates from earlier options and uses default if not set", func(t *testing.T) {
		templateOptionName := "templatedOption"
		templatedOptionDefault := fmt.Sprintf(`{{.%s}}-templated`, optionName)
		// simulate setting a value for first option and use default for next
		gt.InScanner = bufio.NewScanner(strings.NewReader(optionValue + "\n\n"))
		gt.Options = []option.Option{
			{
				Name:    optionName,
				Default: "theDefault",
			},
			{
				Name:    templateOptionName,
				Default: templatedOptionDefault,
			},
		}

		values, err := gt.GetOptionToValueInteractively()
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{optionName: optionValue, templateOptionName: fmt.Sprintf("%s-templated", optionValue)}, values)
	})

	t.Run("does not display options that have dependencies that are not met", func(t *testing.T) {
		dependentOptionName := "dependentOption"
		// simulate accepting the defaults
		gt.InScanner = bufio.NewScanner(strings.NewReader("\n\n"))

		out := &bytes.Buffer{}
		gt.Out = out

		gt.Options = []option.Option{
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

		values, err := gt.GetOptionToValueInteractively()
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

		gt.Options = []option.Option{
			{
				Name:    optionName,
				Default: true,
			},
			{
				Name:    intOptionName,
				Default: 3,
			},
		}

		values, err := gt.GetOptionToValueInteractively()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(values))
		assert.Equal(t, false, values[optionName])
		assert.Equal(t, 4, values[intOptionName])
	})
}
