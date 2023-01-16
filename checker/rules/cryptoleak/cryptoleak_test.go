package cryptoleak

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/analysistest"
	"github.com/lotusirous/codescan/checker/multirunner"
)

type Test struct {
	Dir   string
	Diags []*analysis.Diagnostic
}

func TestFromFileSystem(t *testing.T) {
	checks := []Test{
		{
			Dir: "testdata/a",
			Diags: []*analysis.Diagnostic{
				{
					ByAnalyzer: Analyzer,
					Pos:        1,
					Path:       "testdata/a/README.md",
					Message:    "Secret key might be leaked",
				},
				{
					ByAnalyzer: Analyzer,
					Pos:        1,
					Path:       "testdata/a/secret.js",
					Message:    "Secret key might be leaked",
				},
				{
					ByAnalyzer: Analyzer,
					Pos:        1,
					Path:       "testdata/a/secret_leak.py",
					Message:    "Secret key might be leaked",
				},
			},
		},
	}

	for _, tt := range checks {
		t.Run(tt.Dir, func(t *testing.T) {
			got, err := multirunner.Run(tt.Dir, []*analysis.Analyzer{Analyzer})
			if err != nil {
				t.Error(err)
			}

			if err := analysistest.TestDiag(got, tt.Diags); err != nil {
				t.Error(err)
			}
		})
	}

}
