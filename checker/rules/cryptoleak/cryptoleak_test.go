package cryptoleak

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/testutil"
)

func TestFromFileSystem(t *testing.T) {
	cases := []testutil.Test{
		{
			Dir: "testdata/a",
			Diags: []*analysis.Diagnostic{
				{
					ByAnalyzer: "cryptoleak",
					Pos:        1,
					Path:       "testdata/a/secret_leak.py",
				},
				{
					ByAnalyzer: "cryptoleak",
					Pos:        2,
					Path:       "testdata/a/secret_leak.py",
				},
			},
		},
	}

	testutil.Run(t, []*analysis.Analyzer{Analyzer}, cases)
}
