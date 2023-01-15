package results

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func scanRow(sc db.Scanner) (*core.ScanResult, error) {
	var b []byte
	out := new(core.ScanResult)

	err := sc.Scan(
		&out.ID,
		&out.ScanID,
		&out.RepoID,
		&out.Created,
		&out.Updated,
		&b,
	)
	if errors.Is(err, sql.ErrNoRows) {
		err = os.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	var findings []core.Finding
	if err := json.Unmarshal(b, findings); err != nil {
		return nil, err
	}
	out.Findings = findings
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
