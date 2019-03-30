package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sourcegraph/sourcegraph/pkg/conf"
	"github.com/sourcegraph/sourcegraph/pkg/extsvc/gitolite"
)

func (s *Server) handleListGitolite(w http.ResponseWriter, r *http.Request) {
	listGitolite(r.Context(), r.URL.Query().Get("gitolite"), w)
}

// listGitolite is effectively a wrapper around gitolite.Client.ListRepos.  This must currently be
// invoked from gitserver, because only gitserver has the SSH key needed to authenticate to the
// Gitolite API.
func listGitolite(ctx context.Context, gitoliteHost string, w http.ResponseWriter) {
	repos := make([]*gitolite.Repo, 0)

	config, err := conf.GitoliteConfigs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, gconf := range config {
		if gconf.Host != gitoliteHost {
			continue
		}
		rp, err := gitolite.NewClient(gconf.Host).ListRepos(ctx, gconf.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		repos = append(repos, rp...)
	}

	if err := json.NewEncoder(w).Encode(repos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
