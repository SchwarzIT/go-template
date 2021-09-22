package gotemplate_test

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"strings"
	"testing"
)

func TestGT_PrintVersion(t *testing.T) {
	gt := gotemplate.New()

	buffer := &bytes.Buffer{}
	gt.Out = buffer
	gt.PrintVersion()
	assert.Equal(t, config.Version, strings.Trim(buffer.String(), "\n"))
}
