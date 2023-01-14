package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/server/api/render"
)

type scanRequest struct {
	RepoID int64 `json:"repo_id"`
}

type scanResponse struct {
	ScanID int `json:"scan_id"`
}

// HandleScanRepo returns an http.HandlerFunc that processes an http.Request
// to ScanRepo the from the system.
func HandleScanRepo(repos core.RepositoryStore, scans core.ScanStore) http.HandlerFunc {
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

		scan := &core.Scan{
			Repository: repo.ID,
			Status:     core.StatusQueued,
			EnqueuedAt: time.Now().Unix(),
		}

		if err := scans.Create(ctx, scan); err != nil {
			render.InternalError(w, err)
			return
		}

		// TODO(khant): Put to the queue.
		// if err := manager.Enqueue(scan); err != nil {
		// 	render.InternalError(w, err)
		// 	return
		// }

		render.JSON(w, scan, http.StatusOK)
	}
}
