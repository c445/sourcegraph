package gitolite

import (
	context "context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sourcegraph/sourcegraph/pkg/extsvc/gitolite/mock_gitolite"
	"github.com/sourcegraph/sourcegraph/pkg/testutil2"
)

func Test_Client_ListRepos(t *testing.T) {
	tests := []struct {
		client   func(*gomock.Controller) *Client
		expRepos []*Repo
		expErr   error
	}{
		{
			client: func(ctrl *gomock.Controller) *Client {
				return &Client{
					Host: "git@gitolite.example.com",
					ClientDeps: func() ClientDeps {
						m := mock_gitolite.NewMockClientDeps(ctrl)
						m.
							EXPECT().
							CommandOutput(gomock.Any(), gomock.Eq("ssh"), gomock.Eq("git@gitolite.example.com"), gomock.Eq("info")).
							Return(
								[]byte(`hello admin, this is git@gitolite-799486b5db-ghrxg running gitolite3 v3.6.6-0-g908f8c6 on git 2.7.4

 R W    gitolite-admin
 R W    repowith@sign
 R W    testing
`),
								nil,
							)
						return m
					}(),
				}
			},
			expRepos: []*Repo{
				{Name: "gitolite-admin", URL: "git@gitolite.example.com:gitolite-admin"},
				{Name: "repowith@sign", URL: "git@gitolite.example.com:repowith@sign"},
				{Name: "testing", URL: "git@gitolite.example.com:testing"},
			},
		},
		{
			client: func(ctrl *gomock.Controller) *Client {
				return &Client{
					Host: "git@gitolite.example.com",
					ClientDeps: func() ClientDeps {
						m := mock_gitolite.NewMockClientDeps(ctrl)
						m.
							EXPECT().
							CommandOutput(gomock.Any(), gomock.Eq("ssh"), gomock.Eq("git@gitolite.example.com"), gomock.Eq("info")).
							Return(
								[]byte(``),
								nil,
							)
						return m
					}(),
				}
			},
			expRepos: nil,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := test.client(ctrl)

			repos, err := client.ListRepos(context.Background())
			testutil2.CheckDeepEqual(t, test.expRepos, repos, "returned repos did not match")
			testutil2.CheckDeepEqual(t, test.expErr, err, "returned error did not match")
		})
	}
}
