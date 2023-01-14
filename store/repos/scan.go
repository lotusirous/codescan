package repos

import (
	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
)

func toParam(repo *core.Repository) squirrel.Eq {
	return squirrel.Eq{
		"commit":   repo.Commit,
		"http_url": repo.HttpURL,
		"created":  repo.Created,
		"updated":  repo.Updated,
	}
}
