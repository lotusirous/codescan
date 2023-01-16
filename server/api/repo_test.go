package api

import "testing"

func TestValidateGithubURL(t *testing.T) {

	cases := []struct {
		Name     string
		URL      string
		RepoName string
		WantErr  bool
	}{
		{
			Name:     "valid-url",
			URL:      "https://github.com/octocat/Hello-World",
			RepoName: "Hello-World",
			WantErr:  false,
		},
		{
			Name:     "no-proto",
			URL:      "github.com/octocat/Hello-World",
			RepoName: "Hello-World",
			WantErr:  true,
		},
		{
			Name:    "malformed-url",
			URL:     "http//github.com/octocat/Hello-World",
			WantErr: true,
		},
		{
			Name:    "missing-repo-name",
			URL:     "http//github.com/octocat",
			WantErr: true,
		},
		{
			Name:     "long-url",
			URL:      "https://github.com/octocat/Hello-World/very/long/path/in/url",
			RepoName: "Hello-World",
			WantErr:  false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			repoName, err := getRepoName(test.URL)
			if err != nil {
				if !test.WantErr {
					t.Error(err)
				}
			} else {
				if test.RepoName != repoName {
					t.Errorf("got name %s - want %s", repoName, test.RepoName)
				}
			}
		})
	}

}

func TestHandleListRepo(t *testing.T) {

}

func TestHandleCreateRepo(t *testing.T) {

}
