package core

import "context"

// Fetchers fetches the public git repository
type Fetcher interface {
	// Fetch downloads the public repository,
	// The returns a downloaded directory, cleanup function and the error.
	Fetch(ctx context.Context, target string) (string, func(), error)
}
