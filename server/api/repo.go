package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/server/api/render"
)

type createRepoRequest struct {
	RepoURL string `json:"repo_url"`
}

func validateGithubURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if u.Host != "github.com" {
		return fmt.Errorf("not support %s", u.Host)
	}
	if u.Scheme == "" {
		return fmt.Errorf("require URL scheme")
	}
	if u.Path == "" {
		return fmt.Errorf("require a repo")
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

		if err := validateGithubURL(in.RepoURL); err != nil {
			render.BadRequestf(w, "invalid URL: %s", err.Error())
			return
		}

		repo := &core.Repository{
			HttpURL: in.RepoURL,
			Created: time.Now().Unix(),
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
		}
		render.JSON(w, out, http.StatusOK)

	}
}
