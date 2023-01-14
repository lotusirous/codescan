package repos

import (
	"github.com/Masterminds/squirrel"
	"github.com/lotusirous/codescan/core"
)

func toParam(repo *core.Repository) squirrel.Eq {
	return squirrel.Eq{
		"user":     repo.User,
		"commit":   repo.Commit,
		"scm":      repo.SCM,
		"http_url": repo.HttpURL,
		"ssh_url":  repo.SSHURL,
		"name":     repo.Name,
		"created":  repo.Created,
		"updated":  repo.Updated,
	}
}
