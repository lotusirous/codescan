package core

import (
	"time"
)

// GitSummary describes briefly a repo information.
type GitSummary struct {
	Branch     string
	CommitHash string
	CommitTime time.Time
}

// GitFetcher fetches the public git repository
type GitFetcher interface {
	// Clone downloads the remote repository,
	// it returns the path to temp directory, the cleanup function
	Clone(remoteURL string) (string, func() error, error)

	// Summarize extracts the repo summary.
	Summarize(dir string) (*GitSummary, error)
}
