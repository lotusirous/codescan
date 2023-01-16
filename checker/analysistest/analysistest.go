package analysistest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/lotusirous/codescan/checker/analysis"
)

// TestData returns the abs path of the program's "testdata" directory.
var TestData = func() string {
	testdata, err := filepath.Abs("testdata")
	if err != nil {
		log.Fatal(err)
	}
	return testdata
}

// WriteFiles is a helper function that creates a temporary directory
// On success it returns the name of the directory and a cleanup function to delete it.
func WriteFiles(filemap map[string]string) (dir string, cleanup func(), err error) {
	tmp, err := os.MkdirTemp("", "analysistest")
	if err != nil {
		return "", nil, err
	}
	cleanup = func() { os.RemoveAll(tmp) }

	for name, content := range filemap {
		filename := filepath.Join(tmp, name)
		_ = os.MkdirAll(filepath.Dir(filename), 0777) // ignore error
		if err := os.WriteFile(filename, []byte(content), 0666); err != nil {
			cleanup()
			return "", nil, err
		}
	}
	return tmp, cleanup, nil
}

func TestDiag(got, want []*analysis.Diagnostic) error {
	if len(got) != len(want) {
		return fmt.Errorf("want 2 diags have same length")
	}
	for i := 0; i < len(got); i++ {
		if !reflect.DeepEqual(got[i], want[i]) {
			return fmt.Errorf("not equal at %d got %v - want %v", i, got[i], want[i])
		}
	}
	return nil
}
