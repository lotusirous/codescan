package github

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/lotusirous/codescan/core"
)

// New inits the github fetcher
func New(dir, pattern string) core.GitFetcher {
	return &github{dir, pattern}
}

type github struct {
	dir, pattern string
}

// Summarize extracts meta data from a git repo
func (gh *github) Summarize(dir string) (*core.GitSummary, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	out := new(core.GitSummary)
	out.CommitHash = ref.Hash().String()
	out.Branch = ref.Name().String()
	out.CommitTime = commit.Committer.When

	return out, nil
}

// Clone fetches the remote repo to local dir.
// The caller will clean up the temp dir.
func (gh *github) Clone(remoteURL string) (string, func() error, error) {
	dir, err := os.MkdirTemp(gh.dir, gh.pattern)
	if err != nil {
		return "", func() error { return nil }, err
	}
	cleanup := func() error {
		return os.RemoveAll(dir)
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:               remoteURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return dir, cleanup, err
	}

	return dir, cleanup, nil
}
