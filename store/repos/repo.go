package repos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func New(db *sql.DB) core.RepositoryStore {
	return &repoStore{db: db}
}

type repoStore struct {
	db *sql.DB
}

const baseColumns = `repo_id, name, http_url, created, updated`

// Find returns a repository from the datastore.
// The caller should handle fs.ErrNotExist when there is no rows.
func (s *repoStore) Find(ctx context.Context, id int64) (*core.Repository, error) {
	query, args, err := squirrel.Select(baseColumns).
		From("repos").
		Where(squirrel.Eq{"repo_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanRow(row)
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
		"name":     repo.Name,
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

// List implements core.RepositoryStore
func (s *repoStore) List(ctx context.Context) ([]*core.Repository, error) {
	query, args, err := squirrel.Select(baseColumns).
		From("repos").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	out := []*core.Repository{}
	for rows.Next() {
		r, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Close(); err != nil {
		return nil, err
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

	r, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return db.ErrOptimisticLock
	}
	return tx.Commit()
}

// ListRange implements core.RepositoryStore
// func (s *repoStore) ListRange(ctx context.Context, param core.RepoParam) ([]*core.Repository, error) {
// 	panic("unimplemented")
// }
