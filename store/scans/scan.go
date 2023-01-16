package scans

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func New(db *sql.DB) core.ScanStore {
	return &scanStore{db: db}
}

type scanStore struct {
	db *sql.DB
}

const queryBase = `SELECT scan_id, repo_id, status, enqueued_at, started_at, finished_at `

func scanRow(sc db.Scanner) (*core.Scan, error) {
	out := new(core.Scan)
	err := sc.Scan(&out.ID, &out.RepoID, &out.Status, &out.EnqueuedAt, &out.StartedAt, &out.FinishedAt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *scanStore) FindEnqueued(ctx context.Context) ([]*core.Scan, error) {
	query := queryBase + ` FROM scans WHERE status = 'Queued' ORDER BY enqueued_at`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	out := []*core.Scan{}
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

// Find returns a scan from datastore.
func (s *scanStore) Find(ctx context.Context, id int64) (*core.Scan, error) {
	out := new(core.Scan)
	query := queryBase + `FROM scans WHERE scan_id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&out.ID,
		&out.RepoID,
		&out.Status,
		&out.EnqueuedAt,
		&out.StartedAt,
		&out.FinishedAt,
	)
	return out, err
}

// Count counts the scan in the datastore.
func (s *scanStore) Count(ctx context.Context) (int64, error) {
	queryCount := `SELECT COUNT(*) FROM scans`
	var count int64
	err := s.db.QueryRowContext(ctx, queryCount).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes the scan from datastore.
func (s *scanStore) Delete(ctx context.Context, scan *core.Scan) error {
	queryDelete := `DELETE FROM scans WHERE scan_id = ?`
	r, err := s.db.ExecContext(ctx, queryDelete, scan.ID)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("no row affected for scan_id: %d", scan.ID)
	}
	return nil
}

// UpdateStatus update scan to datastore.
func (s *scanStore) Update(ctx context.Context, scan *core.Scan) error {
	var errMsg string
	if scan.Error != nil {
		errMsg = scan.Error.Error()
	}
	b := squirrel.Update("scans").
		SetMap(squirrel.Eq{
			"status":      scan.Status,
			"enqueued_at": scan.EnqueuedAt,
			"started_at":  scan.StartedAt,
			"finished_at": scan.FinishedAt,
			"scan_error":  errMsg,
		}).Where(squirrel.Eq{"scan_id": scan.ID})
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
// It requires a db transaction level to create a scan option.
// multiple client can submit the same scan id. However, the execution is isolated
// in the transaction.
func (s *scanStore) Create(ctx context.Context, scan *core.Scan) error {
	query, args, err := squirrel.Insert("scans").SetMap(squirrel.Eq{
		"repo_id":     scan.RepoID,
		"status":      scan.Status,
		"enqueued_at": scan.EnqueuedAt,
		"started_at":  scan.StartedAt,
		"finished_at": scan.FinishedAt,
	}).ToSql()
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	r, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	scan.ID = id
	return tx.Commit()
}

// List returns a all stored scans.
// It returns fs.ErrNotExist if the scan does not exist.
// The caller owns the returned value.
func (s *scanStore) List(ctx context.Context) ([]*core.Scan, error) {
	query := queryBase + `FROM scans`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	out := []*core.Scan{}
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
