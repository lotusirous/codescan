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
			Name:    "hello-worId",
			Created: 1673746850,
			Updated: 1673746850,
		}

		err := store.Create(noContext, repo)
		if err != nil {
			t.Error(err)
		}
		if repo.ID == 0 {
			t.Errorf("Want repo ID assigned, got %d", repo.ID)
		}
		t.Run("count", testRepoCount(store))
		t.Run("List", testRepoList(store))
		t.Run("Find", testRepoFind(store, repo))
		t.Run("Delete", testDelete(store, repo))

	}
}

func testRepoList(store *repoStore) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := store.List(noContext)
		if err != nil {
			t.Error(err)
		}
		if len(got) != 1 {
			t.Error("must have 1 records got: 5d", len(got))
		}
	}
}

func testDelete(store *repoStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		err := store.Delete(noContext, repo)
		if err != nil {
			t.Error(err)
		}
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
		if got, want := repo.Name, "hello-worId"; got != want {
			t.Errorf("Want repo name %s got %s", want, got)
		}
		if got, want := repo.HttpURL, "https://github.com/octocat/hello-worId"; got != want {
			t.Errorf("Want repo url %q, got %q", want, got)
		}
		if got, want := repo.Created, int64(1673746850); got != want {
			t.Errorf("Want repo created %q, got %q", want, got)
		}
	}
}
