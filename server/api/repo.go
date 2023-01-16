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
	Name string `json:"name"`
	Link string `json:"link"`
}

func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if u.Scheme == "" {
		return fmt.Errorf("no url scheme")
	}
	if u.Path == "" {
		return fmt.Errorf("no path to repository")
	}

	paths := strings.Split(u.Path, "/")
	if len(paths) < 3 {
		return fmt.Errorf("wrong format, ex github.com/{user}/{repo}")
	}
	return err
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

		if err := validateURL(in.Link); err != nil {
			render.BadRequestf(w, "bad URL: %s", err.Error())
			return
		}

		repo := &core.Repository{
			HttpURL: in.Link,
			Name:    in.Name,
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
