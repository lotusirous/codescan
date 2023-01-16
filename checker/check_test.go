package checker

// import (
// 	"testing"

// 	"github.com/lotusirous/codescan/checker/analysis"
// 	"github.com/lotusirous/codescan/checker/testutil"
// )

// var checks = []testutil.Test{
// 	{
// 		Dir: "cryptokey",
// 		Diags: []*analysis.Diagnostic{
// 			{
// 				Path: "/cryptokey/main.js",
// 				Pos:  1,
// 			},
// 			{
// 				Path: "/cryptokey/main.py",
// 				Pos:  1,
// 			},
// 			{
// 				Path: "/cryptokey/main.py",
// 				Pos:  2,
// 			},
// 		},
// 	},
// 	{
// 		Dir: "crypto-with-comment",
// 		Diags: []*analysis.Diagnostic{
// 			{
// 				Path: "/crypto-with-comment/main.js",
// 				Pos:  2,
// 			},
// 		},
// 	},
// }

// func TestRun(t *testing.T) {
// 	for _, test := range checks {
// 		t.Run(test.Dir, func(t *testing.T) {
// 			got, err := Run(testutil.TestData(), DefaultRules())
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			testDiagnostic(t, got, test.Diags)
// 		})
// 	}
// }

// func testDiagnostic(t *testing.T, got, want []*analysis.Diagnostic) {
// 	if len(got) != len(want) {
// 		t.Errorf("diag len not matched got %d - want %d", len(got), len(want))
// 		return
// 	}
// 	for i := 0; i < len(got); i++ {
// 		if got[i].Pos != want[i].Pos {
// 			t.Errorf("invalid pos got: %d - want: %d", got[i].Pos, want[i].Pos)
// 		}
// 		if got[i].Path != want[i].Path {
// 			t.Errorf("invalid path got: %s - want: %s", got[i].Path, want[i].Path)
// 		}
// 	}
// }
