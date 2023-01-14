package scans

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
)

func New(db *sql.DB) core.ScanStore {
	return &repoStore{db: db}
}

type repoStore struct {
	db *sql.DB
}

// Count counts the scan in the datastore.
func (s *repoStore) Count(ctx context.Context) (int64, error) {
	queryCount := `SELECT COUNT(*) FROM scans`
	var count int64
	err := s.db.QueryRowContext(ctx, queryCount).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes the scan from datastore.
func (s *repoStore) Delete(ctx context.Context, scan *core.Scan) error {
	queryDelete := `DELETE FROM repos WHERE repo_id = ?`
	r, err := s.db.ExecContext(ctx, queryDelete, scan.ID)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("no row affected for repo_id: %d", scan.ID)
	}
	return nil
}

// UpdateStatus implements core.ScanStore
func (s *repoStore) Update(ctx context.Context, scan *core.Scan) error {
	b := squirrel.Update("scans").SetMap(squirrel.Eq{"scan_id": 1})
	query, args, err := b.PlaceholderFormat(squirrel.Question).ToSql()
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
	if rows != 1 {
		return err
	}
	return tx.Commit()

}

// Create persists a scan to datastore.
func (s *repoStore) Create(ctx context.Context, scan *core.Scan) error {
	query, args, err := squirrel.Insert("scan").SetMap(squirrel.Eq{
		"repo_id":     scan.Repository,
		"status":      scan.Status,
		"enqueued_at": scan.EnqueuedAt,
		"started_at":  scan.StartedAt,
		"finished_at": scan.FinishedAt,
	}).ToSql()
	if err != nil {
		return err
	}
	r, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	scan.ID = id
	return nil
}

// List returns a all stored scans.
// It returns fs.ErrNotExist if the scan does not exist.
// The caller owns the returned value.
func (s *repoStore) List(ctx context.Context) ([]*core.Scan, error) {
	query := `SELECT scan_id, status, enqueued_at, started_at, finished_at FROM scans`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	out := make([]*core.Scan, 0)
	for rows.Next() {
		var scan *core.Scan
		rows.Scan(&scan.ID, &scan.Status, &scan.EnqueuedAt, &scan.StartedAt, &scan.FinishedAt)
		out = append(out, scan)
	}

	return out, nil
}
