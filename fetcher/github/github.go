package github

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/lotusirous/codescan/core"
)

// New inits the github fetcher
func New(pattern string) core.GitFetcher {
	return &githubFetcher{pattern}
}

type githubFetcher struct {
	pattern string
}

func (gh *githubFetcher) Clone(ctx context.Context, remoteURL string) (string, func(), error) {
	dir, err := os.MkdirTemp("", gh.pattern)
	if err != nil {
		return "", func() {}, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:               remoteURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return "", cleanup, err
	}

	return dir, cleanup, nil
}
