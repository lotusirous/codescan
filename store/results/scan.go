package results

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/fs"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func scanRow(sc db.Scanner) (*core.ScanResult, error) {

	var id, repoID, scanID, created, updated int64
	var commit string
	var findingByte []byte

	err := sc.Scan(
		&id,
		&repoID,
		&scanID,
		&commit,
		&findingByte,
		&created,
		&updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		err = fs.ErrNotExist
	}
	if err != nil {
		return nil, err
	}

	var findings []core.Finding
	if err := json.Unmarshal(findingByte, &findings); err != nil {
		return nil, err
	}
	return &core.ScanResult{
		ID:       id,
		RepoID:   repoID,
		ScanID:   scanID,
		Commit:   commit,
		Findings: findings,
		Created:  created,
		Updated:  updated,
	}, nil
}

func scanRows(rows *sql.Rows) ([]*core.ScanResult, error) {
	out := []*core.ScanResult{}
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
