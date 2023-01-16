package cryptoleak

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/testutil"
)

func TestFromFileSystem(t *testing.T) {
	checks := []testutil.Test{
		{
			Dir: "testdata/a",
			Diags: []*analysis.Diagnostic{
				{
					ByAnalyzer: Analyzer,
					Pos:        1,
					Path:       "testdata/a/secret_leak.py",
				},
				{
					ByAnalyzer: Analyzer,
					Pos:        2,
					Path:       "testdata/a/secret_leak.py",
				},
			},
		},
	}

	t.Run("test-run", testutil.Run([]*analysis.Analyzer{Analyzer}, checks))

}
