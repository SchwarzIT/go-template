package gotemplate_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/schwarzit/go-template/pkg/repos"
	"github.com/stretchr/testify/assert"
)

func TestGT_PrintVersion(t *testing.T) {
	gt := gotemplate.New()

	buffer := &bytes.Buffer{}
	gt.Out = buffer
	gt.PrintVersion()
	assert.Equal(t, config.Version, strings.Trim(buffer.String(), "\n"))
}

func TestGT_CheckVersion(t *testing.T) {
	tests := []struct {
		name          string
		listerFunc    repos.GithubTagListerFunc
		expectWarning bool
	}{
		{
			name: "warning on error",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return nil, errors.New("some error")
			},
			expectWarning: true,
		},
		{
			name: "warning on version mismatch",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return []string{"someOtherVersion"}, nil
			},
			expectWarning: true,
		},
		{
			name: "success on same version",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return []string{config.Version}, nil
			},
			expectWarning: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			gt := gotemplate.GT{
				Streams: gotemplate.Streams{
					Out: out,
					Err: out,
				},
				GithubTagLister: test.listerFunc,
			}

			gt.CheckVersion()
			if test.expectWarning {
				assert.Contains(t, out.String(), "WARNING")
			} else {
				assert.NotContains(t, out.String(), "WARNING")
			}
		})
	}
}
