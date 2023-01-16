package api

import "testing"

func TestValidateGithubURL(t *testing.T) {

	cases := []struct {
		Name    string
		URL     string
		WantErr bool
	}{
		{
			Name:    "valid-url",
			URL:     "https://github.com/octocat/Hello-World",
			WantErr: false,
		},
		{
			Name:    "no-proto",
			URL:     "github.com/octocat/Hello-World",
			WantErr: true,
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
			Name:    "long-url",
			URL:     "https://github.com/octocat/Hello-World/very/long/path/in/url",
			WantErr: false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			err := validateURL(test.URL)
			if err != nil {
				if !test.WantErr {
					t.Error(err)
				}
			}
		})
	}

}

func TestHandleListRepo(t *testing.T) {

}

func TestHandleCreateRepo(t *testing.T) {

}
