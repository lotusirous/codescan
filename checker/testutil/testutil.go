package testutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/multirunner"
)

type Test struct {
	Dir   string
	Diags []*analysis.Diagnostic
}

// TestData returns the abs path of the program's "testdata" directory.
var TestData = func() string {
	testdata, err := filepath.Abs("testdata")
	if err != nil {
		log.Fatal(err)
	}
	return testdata
}

func TestDiagnostic(t *testing.T, got, want []*analysis.Diagnostic) {
	if len(got) != len(want) {
		t.Error("")
	}
	for i := 0; i < len(got); i++ {
		if got[i].Pos != want[i].Pos {
			t.Errorf("invalid pos got: %d - want: %d", got[i].Pos, want[i].Pos)
		}
		if got[i].ByAnalyzer != want[i].ByAnalyzer {
			t.Errorf("invalid pos got: %s - want: %s", got[i].ByAnalyzer, want[i].ByAnalyzer)
		}

		if got[i].Path != want[i].Path {
			t.Errorf("invalid path got: %s - want: %s", got[i].Path, want[i].Path)
		}
	}
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

// Run applies an analysis to the packages
func Run(analyzers []*analysis.Analyzer, tests []Test) func(t *testing.T) {
	return func(t *testing.T) {
		for _, tt := range tests {
			diags, err := run(tt.Dir, analyzers)
			if err != nil {
				t.Errorf("enable to run: %s -  err: %v", tt.Dir, err)
			}

			if err := testDiag(diags, tt.Diags); err != nil {
				t.Error(err)
			}
		}
	}
}

func testDiag(got, want []*analysis.Diagnostic) error {
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

func run(dir string, analyzers []*analysis.Analyzer) ([]*analysis.Diagnostic, error) {
	return multirunner.Run(dir, analyzers)
}
