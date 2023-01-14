package repos

import (
	"context"
	"testing"

	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/store/db/dbtest"
)

var noContext = context.TODO()

func TestScanStore(t *testing.T) {
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
	store := New(conn).(*repoStore)

	t.Run("Create", testRepoCreate(store))

}

func testRepoCreate(store *repoStore) func(t *testing.T) {
	return func(t *testing.T) {
		repo := &core.Repository{
			HttpURL: "https://github.com/octocat/hello-worId",
			Commit:  "7e068727fdb347b685b658d2981f8c85f7bf0585",
		}

		err := store.Create(noContext, repo)
		if err != nil {
			t.Error(err)
		}
		if repo.ID == 0 {
			t.Errorf("Want repo ID assigned, got %d", repo.ID)
		}
		t.Run("count", testRepoCount(store))
		t.Run("Find", testRepoFind(store, repo))

	}
}

func testRepoCount(store *repoStore) func(t *testing.T) {
	return func(t *testing.T) {
		count, err := store.Count(noContext)
		if err != nil {
			t.Error(err)
		}
		if got, want := count, int64(1); got != want {
			t.Errorf("Want repo table count %d, got %d", want, got)
		}
	}
}

func testRepoFind(store *repoStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		r, err := store.Find(noContext, repo.ID)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Fields", testRepo(r))
		}
	}
}

func testRepo(repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := repo.HttpURL, "https://github.com/octocat/hello-worId"; got != want {
			t.Errorf("Want repo url %q, got %q", want, got)
		}
		if got, want := repo.Commit, "7e068727fdb347b685b658d2981f8c85f7bf0585"; got != want {
			t.Errorf("Want repo commit %q, got %q", want, got)
		}
	}
}
