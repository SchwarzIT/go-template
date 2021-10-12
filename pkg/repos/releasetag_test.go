package repos_test

import (
	"context"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/schwarzit/go-template/pkg/repos"
	"github.com/stretchr/testify/assert"
)

func TestLatestGithubReleaseTag(t *testing.T) {
	tests := []struct {
		name       string
		listerFunc repos.GithubTagListerFunc
		expectErr  bool
		expectTag  *semver.Version
	}{
		{
			name: "error is returned",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return nil, errors.New("some error")
			},
			expectErr: true,
		},
		{
			name: "no tag is returned",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return []string{}, nil
			},
			expectErr: true,
		},
		{
			name: "tags are returned, but not semver",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return []string{"tag0", "tag1"}, nil
			},
			expectErr: true,
		},
		{
			name: "semver tags are returned",
			listerFunc: func(ctx context.Context, owner, repo string) ([]string, error) {
				return []string{"v1.0.0", "v1.0.2"}, nil
			},
			expectErr: false,
			expectTag: semver.MustParse("v1.0.2"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tag, err := repos.LatestGithubReleaseTag(test.listerFunc, "", "")
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectTag, tag)
		})
	}
}
