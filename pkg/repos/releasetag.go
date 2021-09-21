package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pkg/errors"
)

// LatestReleaseTag returns the latest release tag for a given repo.
// The defaultTag is returned if some error occurs or now tags are found.
// Currently only Github releases are supported.
func LatestReleaseTag(repo string) (string, error) {
	repoURL, err := url.Parse(repo)
	if err != nil {
		return "", err
	}

	if repoURL.Hostname() != "github.com" {
		return "", err
	}

	ghAPIPath := "https://" + path.Join("api.github.com/repos", repoURL.Path, "/tags")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ghAPIPath, nil)
	if err != nil {
		return "", err
	}

	client := http.Client{Timeout: time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(ErrStatus(resp.StatusCode), "could not get tags for repo")
	}

	defer resp.Body.Close()

	var releases []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}

	if len(releases) < 1 {
		return "", err
	}

	return releases[0]["name"].(string), nil
}

type ErrStatus int

func (e ErrStatus) Error() string {
	return fmt.Sprintf("unexpected status code: %d", int(e))
}
