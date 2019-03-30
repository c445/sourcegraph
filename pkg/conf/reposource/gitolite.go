package reposource

import (
	"strings"

	"github.com/sourcegraph/sourcegraph/pkg/api"
	"github.com/sourcegraph/sourcegraph/schema"
)

type Gitolite struct {
	*schema.GitoliteConnection
}

var _ RepoSource = Gitolite{}

func (c Gitolite) CloneURLToRepoName(cloneURL string) (repoName api.RepoName, err error) {
	parsedCloneURL, err := parseCloneURL(cloneURL)
	if err != nil {
		return "", err
	}
	parsedHostURL, err := parseCloneURL(c.Host)
	if err != nil {
		return "", err
	}
	if parsedHostURL.Hostname() != parsedCloneURL.Hostname() {
		return "", nil
	}
	return GitoliteRepoName(c.Prefix, strings.TrimPrefix(strings.TrimSuffix(parsedCloneURL.Path, ".git"), "/")), nil
}

func GitoliteRepoName(prefix, gitoliteName string) api.RepoName {
	gitoliteNameWithNoIllegalChars := strings.Replace(gitoliteName, "@", "-", -1)
	return api.RepoName(strings.NewReplacer(
		"{prefix}", prefix,
		"{gitoliteName}", gitoliteNameWithNoIllegalChars,
	).Replace("{prefix}{gitoliteName}"))
}
