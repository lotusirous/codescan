package github

import (
	"context"
	"testing"
)

var noContext = context.TODO()

func TestGitFetcher(t *testing.T) {

	cases := []struct {
		URL string
	}{
		{
			URL: "https://github.com/octocat/Hello-World",
		},
	}

	gh := githubFetcher{"testdata"}

	for _, tt := range cases {
		t.Run(tt.URL, func(t *testing.T) {
			dir, cleanup, err := gh.Clone(noContext, tt.URL)
			if err != nil {
				t.Error(err)
			}

			if dir == "" {
				t.Error("Empty tempdir after clone")
			}
			t.Logf("clone to dir: %s", dir)
			cleanup()
		})
	}
}
