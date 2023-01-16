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
	"github.com/rs/zerolog/log"
)

type scanRequest struct {
	RepoID int64 `json:"repoID"`
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

type findScanResponse struct {
	ID         int64          `json:"id"`
	RepoName   string         `json:"repoName"`
	RepoURL    string         `json:"repoURL"`
	Status     string         `json:"status"` // refer to status scanning job
	EnqueuedAt int64          `json:"enqueuedAt"`
	StartedAt  int64          `json:"startedAt"`
	FinishedAt int64          `json:"finishedAt"`
	Findings   []core.Finding `json:"findings,omitempty"`
	Commit     string         `json:"commit,omitempty"`
}

// HandleFindScan returns an http.HandlerFunc that processes an http.Request
// to FindScan the  from the system.
func HandleFindScan(scans core.ScanStore, repos core.RepositoryStore, results core.ScanResultStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		scan, err := scans.Find(ctx, int64(id))
		if err != nil {
			render.NotFoundf(w, "not found %d", id)
			return
		}

		repo, err := repos.Find(ctx, scan.RepoID)
		if err != nil {
			log.Error().Err(err).Msgf("unable to find repo with scan %d", scan.RepoID)
			render.InternalError(w, err)
			return
		}

		findings := []core.Finding{}
		var commit string
		if scan.Status == core.StatusSuccess {
			res, err := results.Find(ctx, scan.ID)
			if err != nil {
				render.InternalError(w, err)
				return
			}
			findings = res.Findings
			commit = res.Commit
		}

		render.JSON(w, findScanResponse{
			ID:         scan.ID,
			RepoName:   repo.Name,
			RepoURL:    repo.HttpURL,
			Status:     scan.Status,
			EnqueuedAt: scan.EnqueuedAt,
			StartedAt:  scan.StartedAt,
			FinishedAt: scan.FinishedAt,
			Findings:   findings,
			Commit:     commit,
		}, http.StatusOK)
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

// HandleDeleteScan returns an http.HandlerFunc that processes an http.Request
// to delete the scan from the system.
func HandleDeleteScan(scans core.ScanStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.BadRequest(w, err)
			return
		}
		sc, err := scans.Find(ctx, int64(id))
		if err != nil {
			render.NotFoundf(w, "not found scan id %d", id)
			return
		}

		if err := scans.Delete(ctx, sc); err != nil {
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
