package gitolite

import (
	"context"
	"testing"

	"github.com/sourcegraph/sourcegraph/pkg/testutil2"
)

func Test_Client_ListRepos(t *testing.T) {
	tests := []struct {
		client   Client
		expRepos []*Repo
		expErr   error
	}{{
		client: Client{
			Host: "gitolite.example.com",
			clientMock: &clientMock{
				mockCommandOutput: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
					// TODO: want nicer way to specify mock
					return nil, nil
				},
			},
		},
		expRepos: []*Repo{
			// TODO
		},
	}}
	for _, test := range tests {
		repos, err := test.client.ListRepos(context.Background())
		testutil2.CheckDeepEqual(t, test.expRepos, repos, "returned repos did not match")
		testutil2.CheckDeepEqual(t, test.expErr, err, "returned error did not match")
	}
}
