package github

import (
	"testing"

	"github.com/lotusirous/codescan/core"
)

func TestGitFetcher(t *testing.T) {

	cases := []struct {
		URL  string
		Want *core.GitSummary
	}{
		{
			URL: "https://github.com/octocat/Hello-World",
			Want: &core.GitSummary{
				Branch:     "refs/heads/master",
				CommitHash: "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			},
		},
		{
			URL: "https://github.com/octocat/hello-worId",
			Want: &core.GitSummary{
				Branch:     "refs/heads/master",
				CommitHash: "7e068727fdb347b685b658d2981f8c85f7bf0585",
			},
		},
	}

	gh := &github{dir: "testdata", pattern: "gh"}

	t.Parallel()
	for _, tt := range cases {
		t.Run(tt.URL, func(t *testing.T) {
			dir, cleanup, err := gh.Clone(tt.URL)
			if err != nil {
				t.Error(err)
			}

			if dir == "" {
				t.Error("Empty tempdir after clone")
			}

			got, err := gh.Summarize(dir)
			if err != nil {
				t.Error(err)
				return
			}

			if got.CommitHash != tt.Want.CommitHash {
				t.Errorf("commit hash got: %s - want: %s", got.CommitHash, tt.Want.CommitHash)
			}

			if got.Branch != tt.Want.Branch {
				t.Errorf("branch got: %s - want: %s", got.Branch, tt.Want.Branch)
			}

			if err := cleanup(); err != nil {
				t.Error(err)
			}

		})
	}
}
