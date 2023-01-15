package checker

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/testutil"
)

var checks = []testutil.Test{
	{Dir: "cryptokey", Diags: []*analysis.Diagnostic{
		{
			ByAnalyzer: "G402",
			Path:       "src/cryptokey/main.js",
			Pos:        1,
		},
		{
			ByAnalyzer: "G402",
			Path:       "src/cryptokey/main.py",
			Pos:        1,
		},
		{
			ByAnalyzer: "G402",
			Path:       "src/cryptokey/main.py",
			Pos:        2,
		},
	}},
}

func makeAnalyzers(map[string]*analysis.Analyzer) []*analysis.Analyzer {
	ans := make([]*analysis.Analyzer, 0)
	for k, a := range DefaultRules {
		a.Name = k
		ans = append(ans, a)
	}
	return ans
}

func TestAnalyzeGroup(t *testing.T) {
	ans := makeAnalyzers(DefaultRules)
	testutil.Run(t, ans, checks)
}

func TestRun(t *testing.T) {
	for _, tt := range checks {
		var testdata = testutil.TestData
		got, err := Run(testdata(), makeAnalyzers(DefaultRules))
		if err != nil {
			t.Error(err)
			return
		}
		testutil.TestDiagnostic(t, got, tt.Diags)
	}

}
