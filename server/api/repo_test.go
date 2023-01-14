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
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			err := validateGithubURL(test.URL)
			got := err != nil
			if test.WantErr != got {
				t.Errorf("url %s is required to %v", test.URL, test.WantErr)
			}
		})
	}

}
