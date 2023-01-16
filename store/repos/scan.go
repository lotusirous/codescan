package repos

import (
	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db"
)

func toParam(repo *core.Repository) squirrel.Eq {
	return squirrel.Eq{
		"http_url": repo.HttpURL,
		"created":  repo.Created,
		"updated":  repo.Updated,
	}
}

func scanRow(sc db.Scanner) (*core.Repository, error) {
	repo := new(core.Repository)
	err := sc.Scan(
		&repo.ID,
		&repo.Name,
		&repo.HttpURL,
		&repo.Created,
		&repo.Updated,
	)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
