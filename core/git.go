package core

import (
	"context"
	"time"
)

// GitSummary describes briefly a repo information.
type GitSummary struct {
	RemoteAddr    string
	BranchName    string
	CommitHash    string
	CommitMessage string
	CommitTime    time.Time
	AuthorName    string
	AuthorEmail   string
}

// GitFetcher fetches the public git repository
type GitFetcher interface {
	// Clone downloads the remote repository,
	// it returns the path to temp directory, the cleanup function
	Clone(ctx context.Context, remoteURL string) (string, func(), error)

	// Summarize extracts the repo summary.
	Summarize(ctx context.Context, dir string) (*GitSummary, error)
}
