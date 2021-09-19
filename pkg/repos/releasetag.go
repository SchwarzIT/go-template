package repos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"
)

// LatestReleaseTag returns the latest release tag for a given repo.
// The defaultTag is returned if some error occurs or now tags are found.
// Currently only Github releases are supported.
func LatestReleaseTag(repo, defaultTag string) string {
	repoURL, err := url.Parse(repo)
	if err != nil {
		return defaultTag
	}

	if repoURL.Hostname() != "github.com" {
		return defaultTag
	}

	ghAPIPath := "https://" + path.Join("api.github.com/repos", repoURL.Path, "/tags")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ghAPIPath, nil)
	if err != nil {
		return defaultTag
	}

	client := http.Client{Timeout: time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return defaultTag
	}

	defer resp.Body.Close()

	var releases []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return defaultTag
	}

	if len(releases) < 1 {
		return defaultTag
	}

	return releases[0]["name"].(string)
}
