package results

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
)

func New(conn *sql.DB) core.ScanResultStore {
	return &resultStore{db: conn}
}

type resultStore struct {
	db *sql.DB
}

// Create return the status of adding a scan result to datastore.
func (s *resultStore) Create(ctx context.Context, result *core.ScanResult) error {
	query, args, err := squirrel.Insert("scan_results").SetMap(squirrel.Eq{
		"scan_id":  result.ScanID,
		"repo_id":  result.RepoID,
		"created":  result.Created,
		"updated":  result.Updated,
		"findings": result.Findings,
	}).ToSql()
	if err != nil {
		return err
	}
	r, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("insert scan result not affected for %d", result.RepoID)
	}
	return err

}

// Find returns a scan result from datastore.
func (s *resultStore) Find(ctx context.Context, id int64) (*core.ScanResult, error) {
	query, args, err := squirrel.
		Select("scan_result_id,scan_id,repo_id,created,updated").
		From("scan_results").
		Where(squirrel.Eq{"scan_result_id": id}).
		PlaceholderFormat(squirrel.Colon).
		ToSql()

	if err != nil {
		return nil, err
	}

	r := s.db.QueryRowContext(ctx, query, args...)
	return scanRow(r)
}

// List returns a list of scan result from datastore.
func (s *resultStore) List(ctx context.Context) ([]*core.ScanResult, error) {
	query, args, err := squirrel.
		Select("scan_result_id,scan_id,repo_id,created,updated").
		From("scan_results").
		PlaceholderFormat(squirrel.Colon).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scanRows(rows)
}
