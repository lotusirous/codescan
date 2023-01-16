package scans

import (
	"context"
	"testing"
	"time"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db/dbtest"
)

var noContext = context.Background()

func TestScanStore(t *testing.T) {
	conn, err := dbtest.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = dbtest.Reset(conn)
		_ = dbtest.Disconnect(conn)
	}()
	store := New(conn).(*scanStore)

	t.Run("Create", testScanCreate(store))
}

func testScanCreate(store *scanStore) func(t *testing.T) {
	return func(t *testing.T) {
		sc := &core.Scan{
			RepoID: 1,
			Status: core.StatusQueued,
		}

		err := store.Create(noContext, sc)
		if err != nil {
			t.Error(err)
		}
		if sc.ID == 0 {
			t.Errorf("Want scan id  assigned, got %d", sc.ID)
		}
		t.Run("Count", testScanCount(store))
		t.Run("Find", testScanFind(store, sc))
		t.Run("FindEnqueued", testScanFindEnqueued(store))
		t.Run("Update", testScanUpdate(store, sc))
	}
}

func testScanFindEnqueued(store *scanStore) func(*testing.T) {
	return func(t *testing.T) {
		out, err := store.FindEnqueued(noContext)
		if err != nil {
			t.Error(err)
		}
		if got := len(out); got != 1 {
			t.Errorf("Want 1 enqueued rows got: %d", len(out))
		}
	}
}

func testScanUpdate(store *scanStore, sc *core.Scan) func(*testing.T) {
	return func(t *testing.T) {
		sc.Status = core.StatusInProgress
		sc.EnqueuedAt = time.Now().Unix()
		err := store.Update(noContext, sc)
		if err != nil {
			t.Error(err)
		}

		out, err := store.Find(noContext, sc.ID)
		if err != nil {
			t.Error(err)
		}
		if got := out.Status; got != core.StatusInProgress {
			t.Errorf("Want status updated in progress got: %s", got)
		}
		if out.EnqueuedAt == 0 {
			t.Error("Want scan enqueued time not 0")
		}
	}
}

func testScanCount(store *scanStore) func(t *testing.T) {
	return func(t *testing.T) {
		count, err := store.Count(noContext)
		if err != nil {
			t.Error(err)
		}
		if got, want := count, int64(1); got != want {
			t.Errorf("Want scans table count %d, got %d", want, got)
		}
	}
}

func testScanFind(store *scanStore, sc *core.Scan) func(t *testing.T) {
	return func(t *testing.T) {
		r, err := store.Find(noContext, sc.ID)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Fields", testScanFields(r))
		}
	}
}

func testScanFields(sc *core.Scan) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := sc.Status, core.StatusQueued; got != want {
			t.Errorf("Want scan status %q, got %q", want, got)
		}
		if got, want := sc.RepoID, int64(1); got != want {
			t.Errorf("Want scan refer to repo id  %q, got %q", want, got)
		}
	}
}
