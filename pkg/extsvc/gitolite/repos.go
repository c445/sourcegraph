package gitolite

import (
	"context"
	"os/exec"
	"strings"

	log15 "gopkg.in/inconshreveable/log15.v2"
)

type Client struct {
	Host string
}

type Repo struct {
	Name string // the name of the repository as it is returned by `ssh git@GITOLITE_HOST info`
	URL  string // the clone URL of the repository
}

func (c *Client) ListRepos(ctx context.Context, blacklistStr string) ([]*Repo, error) {
	out, err := exec.CommandContext(ctx, "ssh", c.Host, "info").Output()
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
