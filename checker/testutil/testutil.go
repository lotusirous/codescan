package testutil

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
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
	tmp, err := ioutil.TempDir("", "analysistest")
	if err != nil {
		return "", nil, err
	}
	cleanup = func() { os.RemoveAll(tmp) }

	for name, content := range filemap {
		filename := filepath.Join(tmp, name)
		os.MkdirAll(filepath.Dir(filename), 0777) // ignore error
		if err := ioutil.WriteFile(filename, []byte(content), 0666); err != nil {
			cleanup()
			return "", nil, err
		}
	}
	return tmp, cleanup, nil
}

// Run applies an analysis to the packages
func Run(t *testing.T, analyzers []*analysis.Analyzer, tests []Test) {
	for _, test := range tests {
		var got []*analysis.Diagnostic
		for _, a := range analyzers {
			pass, err := analysis.Load(test.Dir)
			if err != nil {
				t.Error(err)
			}

			diag, err := a.RunSingle(pass)
			if err != nil {
				t.Errorf("analyzer %s failed: %v", a.Name, err)
			}

			for _, d := range diag {
				got = append(got, d)
			}

			TestDiagnostic(t, got, test.Diags)
		}

	}

}
