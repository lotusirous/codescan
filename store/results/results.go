package results

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func New(conn *sql.DB) core.ScanResultStore {
	return &resultStore{db: conn}
}

type resultStore struct {
	db *sql.DB
}

func (s *resultStore) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM scan_results`
	var cnt int64
	err := s.db.QueryRowContext(ctx, query).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// Create return the status of adding a scan result to datastore.
func (s *resultStore) Create(ctx context.Context, result *core.ScanResult) error {

	findingByte, err := json.Marshal(result.Findings)
	if err != nil {
		return err
	}

	query, args, err := squirrel.Insert("scan_results").SetMap(squirrel.Eq{
		"scan_id":  result.ScanID,
		"repo_id":  result.RepoID,
		"commit":   result.Commit,
		"created":  result.Created,
		"updated":  result.Updated,
		"findings": findingByte,
	}).PlaceholderFormat(squirrel.Question).ToSql()
	if err != nil {
		return fmt.Errorf("%w - sql: %s", err, query)
	}

	r, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	result.ID = id
	return err

}

const queryBase = `SELECT scan_result_id, scan_id, repo_id, commit, findings, created, updated `

// Find returns a scan result from datastore.
func (s *resultStore) Find(ctx context.Context, id int64) (*core.ScanResult, error) {
	query := queryBase + `
	FROM scan_results
	WHERE scan_result_id = ?
	`
	r := s.db.QueryRowContext(ctx, query, id)
	return scanRow(r)
}

// Find returns a scan result from datastore.
func (s *resultStore) FindScan(ctx context.Context, scanID int64) (*core.ScanResult, error) {
	query := queryBase + `
	FROM scan_results
	WHERE scan_id = ?
	`
	r := s.db.QueryRowContext(ctx, query, scanID)
	return scanRow(r)
}

// List returns a list of scan result from datastore.
// func (s *resultStore) List(ctx context.Context) ([]*core.ScanResult, error) {
// 	query, args, err := squirrel.
// 		Select("scan_result_id, scan_id, repo_id, commit, findings, created, updated").
// 		From("scan_results").
// 		PlaceholderFormat(squirrel.Question).
// 		ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := s.db.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return scanRows(rows)
// }

// DeleteByScan removes the scan by given id.
func (s *resultStore) DeleteByScan(ctx context.Context, scanID int64) error {
	query := `DELETE FROM scan_results WHERE scan_id = ?`

	r, err := s.db.ExecContext(ctx, query, scanID)
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
	return nil
}
