package core

import "context"

// GitFetcher fetches the public git repository
type GitFetcher interface {
	// Clone downloads the remote repository,
	// it returns the path to temp directory, the cleanup function
	Clone(ctx context.Context, remoteURL string) (string, func(), error)
}
