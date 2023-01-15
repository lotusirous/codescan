package results

import (
	"database/sql"
	"errors"
	"io/fs"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func scanRow(sc db.Scanner) (*core.ScanResult, error) {
	out := new(core.ScanResult)
	err := sc.Scan(
		&out.ID,
		&out.RepoID,
		&out.ScanID,
		&out.Commit,
		&out.Created,
		&out.Updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		err = fs.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}

func scanRows(rows *sql.Rows) ([]*core.ScanResult, error) {
	out := make([]*core.ScanResult, 0)
	for rows.Next() {
		r, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
