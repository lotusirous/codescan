package core

import "context"

// Repository represents a source code repository.
type Repository struct {
	ID      int64  `json:"id"`
	User    string `json:"user"`   // submit by user
	Commit  string `json:"commit"` // the latest commit
	HttpURL string `json:"git_http_url"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

// RepoParam defines repo query parameters.
type RepoParam struct {
	Sort bool
	Page int64
	Size int64
}

// RepositoryStore defines operations for working with repositories.
type RepositoryStore interface {

	// Find a repository by a id.
	Find(ctx context.Context, id int64) (*Repository, error)

	// Create persists a new repository to the datastore.
	Create(ctx context.Context, repo *Repository) error

	// List returns a list of repositories from the datastore.
	List(context.Context) ([]*Repository, error)

	// ListRange returns a range of repo from the datastore.
	// ListRange(context.Context, RepoParam) ([]*Repository, error)

	// Delete deletes a repository from the datastore.
	Delete(context.Context, *Repository) error

	// Update persists repository changes to the datastore.
	Update(context.Context, *Repository) error

	// Count returns a count of activated repositories.
	Count(context.Context) (int64, error)
}
