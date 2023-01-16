package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/server/api/render"
)

type scanRequest struct {
	RepoID int64 `json:"repo_id"`
}

// HandleScanRepo returns an http.HandlerFunc that processes an http.Request
// to ScanRepo the from the system.
func HandleScanRepo(manager core.ScanScheduler, repos core.RepositoryStore, scans core.ScanStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := new(scanRequest)
		ctx := r.Context()
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			render.BadRequest(w, err)
			return
		}

		repo, err := repos.Find(ctx, in.RepoID)
		if errors.Is(err, os.ErrNotExist) {
			render.NotFoundf(w, "not found repo %d", in.RepoID)
			return
		}
		if err != nil {
			render.InternalError(w, err)
			return
		}

		scan, err := manager.ScanRepo(ctx, repo)
		if err != nil {
			render.InternalError(w, err)
			return
		}

		render.JSON(w, scan, http.StatusOK)
	}
}

// HandleFindScan returns an http.HandlerFunc that processes an http.Request
// to FindScan the  from the system.
func HandleFindScan(scans core.ScanStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.BadRequestf(w, "invalid id: %s", id)
			return
		}

		scan, err := scans.Find(r.Context(), int64(id))
		if err != nil {
			render.NotFound(w, err)
			return
		}
		render.JSON(w, scan, http.StatusOK)
	}
}

// HandleListScan returns an http.HandlerFunc that processes an http.Request
// to list all scans from the system.
func HandleListScan(scans core.ScanStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := scans.List(r.Context())
		if err != nil {
			render.InternalError(w, err)
			return
		}
		render.JSON(w, out, http.StatusOK)
	}
}
