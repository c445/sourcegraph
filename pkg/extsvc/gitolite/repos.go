package gitolite

import (
	"context"
	"os/exec"
	"strings"

	log15 "gopkg.in/inconshreveable/log15.v2"
)

// Client is a client for the Gitolite API.
//
// IMPORTANT: in order to authenticate to the Gitolite API, the client must be invoked from a
// service in an environment that contains a Gitolite-authorized SSH key. As of writing, only
// gitserver meets this criterion. (I.e., only invoke this from gitserver.)
type Client struct {
	Host string

	*clientMock
}

type Repo struct {
	Name string // the name of the repository as it is returned by `ssh git@GITOLITE_HOST info`
	URL  string // the clone URL of the repository
}

func (c *Client) ListRepos(ctx context.Context) ([]*Repo, error) {
	out, err := c.commandOutput(ctx, "ssh", c.Host, "info")
	if err != nil {
		log15.Error("listing gitolite failed", "error", err, "out", string(out))
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var repos []*Repo
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "R" {
			continue
		}
		name := fields[len(fields)-1]
		if len(fields) >= 2 && fields[0] == "R" {
			repos = append(repos, &Repo{
				Name: name,
				URL:  c.Host + ":" + name,
			})
		}
	}

	return repos, nil
}

type clientMock struct {
	mockCommandOutput func(ctx context.Context, name string, arg ...string) ([]byte, error)
}

func (c *clientMock) commandOutput(ctx context.Context, name string, arg ...string) ([]byte, error) {
	if c != nil {
		return c.mockCommandOutput(ctx, name, arg...)
	}
	return exec.CommandContext(ctx, name, arg...).Output()
}
