package multirunner

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/analysistest"
)

func TestError(t *testing.T) {
	fail := &analysis.Analyzer{
		Name: "fail",
		Meta: analysis.Meta{
			Description: "fail analysis",
			Severity:    "LOW",
		},
		Run: func(pass analysis.Pass) ([]*analysis.Diagnostic, error) {
			return nil, nil
		},
	}

	files := map[string]string{
		"src/test.go": `public_key := "foobar"`,
	}

	testdata, cleanup, err := analysistest.WriteFiles(files)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := Run(testdata, []*analysis.Analyzer{fail}); err != nil {
		t.Error(err)
	}

	cleanup()
}
