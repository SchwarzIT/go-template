package repos

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// LatestReleaseTag returns the latest release tag for a given repo.
// The defaultTag is returned if some error occurs or now tags are found.
// Currently only Github releases are supported.
func LatestReleaseTag(repo string, defaultTag string) string {
	repoURL, err := url.Parse(repo)
	if err != nil {
		return defaultTag
	}

	if repoURL.Hostname() != "github.com" {
		return defaultTag
	}

	ghAPIPath := "https://" + path.Join("api.github.com/repos", repoURL.Path, "/tags")
	resp, err := http.Get(ghAPIPath)
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
