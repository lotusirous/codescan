package repos

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func New(db *db.DB) core.RepositoryStore {
	return &repoStore{db: db}
}

type repoStore struct {
	db *db.DB
}

// Count implements core.RepositoryStore
func (s *repoStore) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM repos`
	var count int64
	err := s.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CREATE TABLE repositories (
// 	repo_id INTEGER,
// 	user VARCHAR(255),
// 	commit VARCHAR(255),
// 	scm VARCHAR(255),
// 	http_url VARCHAR(255),
// 	ssh_url VARCHAR(255),
// 	name VARCHAR(255),
// 	created INTEGER,
// 	updated INTEGER
//   );

// Create implements core.RepositoryStore
func (s *repoStore) Create(ctx context.Context, repo *core.Repository) error {
	query, args, err := squirrel.Insert("repos").SetMap(squirrel.Eq{
		"user":     repo.User,
		"commit":   repo.Commit,
		"scm":      repo.SCM,
		"http_url": repo.HttpURL,
		"ssh_url":  repo.SSHURL,
		"name":     repo.Name,
		"created":  repo.Created,
		"updated":  repo.Updated,
	}).PlaceholderFormat(squirrel.Question).ToSql()
	if err != nil {
		return err
	}
	return s.db.Exec(ctx, query, args...)
}

// Delete implements core.RepositoryStore
func (s *repoStore) Delete(ctx context.Context, repo *core.Repository) error {
	query := `DELETE FROM repos WHERE repo_id = ?`
	return s.db.Exec(ctx, query, repo.ID)
}

const stmtFind = `
SELECT repo_id
	,user
	,commit
	,scm
	,http_url
	,ssh_url
	,name
	,created
	,updated
FROM repos
`

// List implements core.RepositoryStore
func (s *repoStore) List(ctx context.Context) ([]*core.Repository, error) {
	rows, err := s.db.Query(ctx, stmtFind)
	if err != nil {
		return nil, err
	}
	out := make([]*core.Repository, 0)
	for rows.Next() {
		var repo *core.Repository
		err := rows.Scan(&repo.ID, &repo.User, &repo.Commit, &repo.SCM,
			&repo.HttpURL, &repo.SSHURL, &repo.Name,
			&repo.Created, &repo.Updated)
		if err != nil {
			return nil, err
		}
		out = append(out)

	}
	return out, nil
}

// ListRange implements core.RepositoryStore
func (s *repoStore) ListRange(ctx context.Context, param core.RepoParam) ([]*core.Repository, error) {
	panic("unimplemented")
}

// Update implements core.RepositoryStore
func (s *repoStore) Update(ctx context.Context, repo *core.Repository) error {
	query, args, err := squirrel.Update("repos").SetMap(toParam(repo)).ToSql()
	if err != nil {
		return err
	}
	return s.db.ExecTX(ctx, query, args)
}
