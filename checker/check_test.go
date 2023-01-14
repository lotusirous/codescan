package checker

import (
	"testing"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/testutil"
)

func makeAnalyzers(map[string]*analysis.Analyzer) []*analysis.Analyzer {
	ans := make([]*analysis.Analyzer, 0)
	for k, a := range DefaultRules {
		a.Name = k
		ans = append(ans, a)
	}
	return ans
}

func TestAnalyzeGroup(t *testing.T) {

	checks := []testutil.Test{
		{Dir: "cryptokey", Diags: []*analysis.Diagnostic{
			{
				ByAnalyzer: "G402",
				Path:       "src/cryptokey/main.js",
				Pos:        1,
			},
		}},
	}

	ans := makeAnalyzers(DefaultRules)
	testutil.Run(t, ans, checks)
}

func TestRun(t *testing.T) {
	checks := []testutil.Test{
		{Dir: "cryptokey", Diags: []*analysis.Diagnostic{
			{
				ByAnalyzer: "G402",
				Path:       "src/cryptokey/main.js",
				Pos:        1,
			},
		}},
	}
	Run("testdata")
}
