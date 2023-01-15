package checker

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/testutil"
)

var checks = []testutil.Test{
	{Dir: "cryptokey", Diags: []*analysis.Diagnostic{
		{
			Path: "src/cryptokey/main.js",
			Pos:  1,
		},
		{
			Path: "src/cryptokey/main.py",
			Pos:  1,
		},
		{
			Path: "src/cryptokey/main.py",
			Pos:  2,
		},
	}},
}

func TestAnalyzeGroup(t *testing.T) {
	testutil.Run(t, DefaultRules(), checks)
}

func TestRun(t *testing.T) {
	for _, tt := range checks {
		var testdata = testutil.TestData
		got, err := Run(testdata(), DefaultRules())
		if err != nil {
			t.Error(err)
			return
		}
		testutil.TestDiagnostic(t, got, tt.Diags)
	}

}
