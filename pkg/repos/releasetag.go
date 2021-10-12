package repos

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

var ErrNoTagsAvailable = errors.New("no tags available")

type GithubTagLister interface {
	ListTags(ctx context.Context, owner, repo string) ([]string, error)
}

type GithubTagListerFunc func(ctx context.Context, owner, repo string) ([]string, error)

func (f GithubTagListerFunc) ListTags(ctx context.Context, owner, repo string) ([]string, error) {
	return f(ctx, owner, repo)
}

// LatestGithubReleaseTag returns the latest release tag for a given repo.
func LatestGithubReleaseTag(lister GithubTagLister, owner, repo string) (*semver.Version, error) {
	tags, err := lister.ListTags(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}

	if len(tags) < 1 {
		return nil, errors.Wrap(ErrNoTagsAvailable, repo)
	}

	latest := &semver.Version{}
	for _, tag := range tags {
		currentVersion, err := semver.NewVersion(tag)
		if err != nil {
			return nil, err
		}

		if currentVersion.GreaterThan(latest) {
			latest = currentVersion
		}
	}

	return latest, nil
}
