package results

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db/dbtest"
)

var sample = &core.ScanResult{
	ScanID:  1,
	RepoID:  1,
	Commit:  "02aeacffe4dfae05956c28421c949c38c69d354c",
	Created: 1673746850,
	Updated: 1673746850,
	Findings: []core.Finding{
		{
			Type:   "static",
			RuleID: "G402",
			Location: core.Location{
				Path: "src/main.js",
				Positions: core.Positions{
					Begin: core.Begin{Line: 1},
				},
			},
			Metadata: core.Metadata{
				Description: "foobar",
				Severity:    "DANGER",
			},
		},
	},
}

var noContext = context.TODO()

func TestResultStore(t *testing.T) {
	conn, err := dbtest.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := dbtest.Reset(conn); err != nil {
			t.Error(err)
		}
		dbtest.Disconnect(conn)
	}()
	store := New(conn).(*resultStore)

	t.Run("Create", testStoreCreate(store))
}

func testStoreCreate(store *resultStore) func(t *testing.T) {
	return func(t *testing.T) {
		r := sample
		err := store.Create(noContext, r)
		if err != nil {
			t.Error(err)
		}
		if r.ID == 0 {
			t.Errorf("Want repo ID assigned, got %d", r.ID)
		}
		t.Run("Find", testStoreFind(store, r))

	}
}

func testStoreFind(store *resultStore, r *core.ScanResult) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := store.Find(noContext, r.ID)
		if err != nil {
			t.Error(err)
		}
		testResult(t, got, sample)
	}
}

func testResult(t *testing.T, got, want *core.ScanResult) {
	if got.Commit != want.Commit {
		t.Errorf("commit not match got: %s - want: %s", got.Commit, want.Commit)
	}
	if got.Created != want.Created {
		t.Errorf("created not match got: %d - want: %d", got.Created, want.Created)
	}

	if got.RepoID != want.RepoID {
		t.Errorf("repo id not match got: %d - want: %d", got.RepoID, want.RepoID)
	}

	if got.ScanID != want.ScanID {
		t.Errorf("scan id not match got: %d - want: %d", got.ScanID, want.ScanID)
	}
	if got.ID != want.ID {
		t.Errorf("id not match got: %d - want: %d", got.ScanID, want.ScanID)
	}

	if cmp.Diff(got.Findings, want.Findings) != "" {
		t.Errorf("findings not match got: %v - want: %v", got.Findings, want.Findings)
	}

}
