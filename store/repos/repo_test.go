package repos

import (
	"context"
	"testing"

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
		dbtest.Reset(conn)
		dbtest.Disconnect(conn)
	}()
	store := New(conn).(*repoStore)

	testCreate(t, store)

}

func testCreate(t *testing.T, store *repoStore) {

}

func testInsert(t *testing.T, store *repoStore) {

}
