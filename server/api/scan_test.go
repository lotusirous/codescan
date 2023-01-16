package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/mock"
)

var mockScan = &core.Scan{
	ID:     1,
	RepoID: 1,
	Status: core.StatusQueued,
}

func TestHandleScanRepo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var (
		repos   = mock.NewMockRepositoryStore(controller)
		manager = mock.NewMockScanScheduler(controller)
	)

	repos.EXPECT().Find(gomock.Any(), gomock.Any()).Do(func(_ context.Context, id int64) (*core.Repository, error) {
		return mockRepo, nil
	})

	manager.EXPECT().ScanRepo(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, repo *core.Repository) error {
		return nil
	})

	// prepare data
	in := new(bytes.Buffer)
	if err := json.NewEncoder(in).Encode(mockScan); err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	HandleScanRepo(manager, repos)(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

}
