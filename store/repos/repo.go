package repos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
)

func New(db *sql.DB) core.RepositoryStore {
	return &repoStore{db: db}
}

type repoStore struct {
	db *sql.DB
}

// Find returns a repository from the datastore.
func (s *repoStore) Find(ctx context.Context, id int64) (*core.Repository, error) {
	query, args, err := squirrel.Select("repo_id, commit, http_url, created, updated").From("repos").ToSql()
	if err != nil {
		return nil, err
	}
	out := new(core.Repository)
	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&out.ID,
		&out.HttpURL,
		&out.Created,
		&out.Updated,
	)
	if err != nil {
		return nil, err
	}
	return out, nil

}

// Count returns the number of repository in the datastore.
func (s *repoStore) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM repos`
	var count int64
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create adds a repository to datastore.
func (s *repoStore) Create(ctx context.Context, repo *core.Repository) error {
	query, args, err := squirrel.Insert("repos").SetMap(squirrel.Eq{
		"http_url": repo.HttpURL,
		"created":  repo.Created,
		"updated":  repo.Updated,
	}).PlaceholderFormat(squirrel.Question).ToSql()
	if err != nil {
		return err
	}
	r, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}
	repo.ID = lastID
	return nil

}

// Delete removes a repository from datastore.
func (s *repoStore) Delete(ctx context.Context, repo *core.Repository) error {
	query := `DELETE FROM repos WHERE repo_id = ?`

	r, err := s.db.ExecContext(ctx, query, repo.ID)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("no rows affected for repo id: %d", repo.ID)
	}
	return nil
}

const stmtFind = `
SELECT repo_id
	,http_url
	,created
	,updated
FROM repos
`

// List implements core.RepositoryStore
func (s *repoStore) List(ctx context.Context) ([]*core.Repository, error) {
	rows, err := s.db.QueryContext(ctx, stmtFind)
	if err != nil {
		return nil, err
	}
	out := make([]*core.Repository, 0)
	for rows.Next() {
		var repo *core.Repository
		err := rows.Scan(&repo.ID,
			&repo.HttpURL,
			&repo.Created, &repo.Updated)
		if err != nil {
			return nil, err
		}
		out = append(out)

	}
	return out, nil
}

// Update updates the repo to datastore.
func (s *repoStore) Update(ctx context.Context, repo *core.Repository) error {
	query, args, err := squirrel.Update("repos").SetMap(toParam(repo)).ToSql()
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	tx.ExecContext(ctx, query, args...)
	return tx.Commit()
}

// ListRange implements core.RepositoryStore
// func (s *repoStore) ListRange(ctx context.Context, param core.RepoParam) ([]*core.Repository, error) {
// 	panic("unimplemented")
// }
