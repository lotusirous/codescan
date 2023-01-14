package scans

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func New(db *db.DB) core.ScanStore {
	return &repoStore{db: db}
}

type repoStore struct {
	db *db.DB
}

// Count counts the scan in the datastore.
func (s *repoStore) Count(ctx context.Context) (int64, error) {
	queryCount := `SELECT COUNT(*) FROM scans`
	var count int64
	err := s.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes the scan from datastore.
func (s *repoStore) Delete(ctx context.Context, scan *core.Scan) error {
	queryDelete := `DELETE FROM repos WHERE repo_id = ?`
	err := s.db.Exec(ctx, queryDelete, scan.ID)
	return err
}

// UpdateStatus implements core.ScanStore
func (s *repoStore) Update(ctx context.Context, scan *core.Scan) error {
	b := squirrel.Update("scans").SetMap(squirrel.Eq{"scan_id": 1})
	query, args, err := b.PlaceholderFormat(squirrel.Question).ToSql()
	if err != nil {
		return err
	}
	return s.db.ExecTX(ctx, query, args...)
}

// Create persists a scan to datastore.
func (s *repoStore) Create(ctx context.Context, scan *core.Scan) error {
	query, args, err := squirrel.Insert("scan").SetMap(squirrel.Eq{
		"repo_id":     scan.Repository.ID,
		"status":      scan.Status,
		"enqueued_at": 0,
		"started_at":  0,
		"finished_at": 0,
	}).ToSql()
	if err != nil {
		return err
	}
	return s.db.Exec(ctx, query, args...)

}

// List returns a all stored scans.
// It returns fs.ErrNotExist if the scan does not exist.
// The caller owns the returned value.
func (s *repoStore) List(ctx context.Context) ([]*core.Scan, error) {
	query := `SELECT scan_id, status, enqueued_at, started_at, finished_at FROM scans`
	rows, err := s.db.Query(ctx, query)
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
