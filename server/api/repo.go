package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/server/api/render"
)

type createRepoRequest struct {
	RepoURL string `json:"repo_url"`
}

// getRepoName validates and return a name of repository.
// I assume the repo format is:
// https://github.com/user/repo
func getRepoName(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if u.Host != "github.com" {
		return "", fmt.Errorf("not support %s", u.Host)
	}
	if u.Scheme == "" {
		return "", fmt.Errorf("require URL scheme")
	}
	if u.Path == "" {
		return "", fmt.Errorf("require a repo")
	}

	paths := strings.Split(u.Path, "/")
	if len(paths) < 3 {
		return "", fmt.Errorf("invalid format github.com/{user}/{repo}")
	}
	return paths[2], err
}

// HandleCreate returns an http.HandlerFunc that processes an http.Request
// to add the repository from the system.
func HandleCreateRepo(repos core.RepositoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		in := new(createRepoRequest)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			render.BadRequest(w, err)
			return
		}

		name, err := getRepoName(in.RepoURL)
		if err != nil {
			render.BadRequestf(w, "invalid URL: %s", err.Error())
			return
		}

		repo := &core.Repository{
			HttpURL: in.RepoURL,
			Name:    name,
			Created: time.Now().Unix(),
			Updated: time.Now().Unix(),
		}

		if err := repos.Create(r.Context(), repo); err != nil {
			render.InternalError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// HandleFind returns an http.HandlerFunc that processes an http.Request
// to Find the from the system.
func HandleListRepo(repos core.RepositoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := repos.List(r.Context())
		if err != nil {
			render.InternalError(w, err)
			return
		}
		render.JSON(w, out, http.StatusOK)
	}
}

// HandleDeleteRepo returns an http.HandlerFunc that processes an http.Request
// to delete the repository from the system.
func HandleDeleteRepo(repos core.RepositoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
		)
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		repo, err := repos.Find(ctx, int64(id))
		if err != nil {
			render.NotFoundf(w, "not found %d", id)
			return
		}
		repos.Delete(ctx, repo)
		w.WriteHeader(http.StatusNoContent)

	}
}
