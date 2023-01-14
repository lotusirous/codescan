package checker

import (
	"fmt"
	"sync"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/rules/cryptoleak"
)

// DefaultRules is the simple rules for detect potential vulnerability.
var DefaultRules = map[string]*analysis.Analyzer{
	"G402": cryptoleak.Analyzer,
}

// InitRules maps the rule name to the analyzer name.
// It replaces the analyzer name to the rule name.
func InitRules(as map[string]*analysis.Analyzer) []*analysis.Analyzer {
	var out []*analysis.Analyzer
	for k, v := range as {
		v.Name = k
		out = append(out, v)
	}
	return out
}

// Run loads the packages specified by args using go/packages,
func Run(dir string, analyzers []*analysis.Analyzer) ([]*analysis.Diagnostic, error) {
	if err := analysis.Validate(analyzers); err != nil {
		return nil, err
	}

	pass, err := analysis.Load(dir)
	if err != nil {
		return nil, err
	}

	actions := analyze(pass, analyzers)
	return collectDiagnostics(actions)
}

func analyze(pass analysis.Pass, analyzers []*analysis.Analyzer) []*action {
	var actions []*action
	for _, a := range analyzers {
		actions = append(actions, &action{a: a, pass: pass})
	}
	// Execute in parallel.
	execAll(actions)
	return actions
}

func collectDiagnostics(actions []*action) ([]*analysis.Diagnostic, error) {
	var out []*analysis.Diagnostic
	for _, a := range actions {
		if a.err != nil {
			return nil, a.err
		}

		for _, diag := range a.diagnostics {
			diag.ByAnalyzer = a.a.Name
			out = append(out, diag)
		}

	}
	return out, nil
}

// An action represents one unit of analysis work.
type action struct {
	once        sync.Once
	pass        analysis.Pass
	a           *analysis.Analyzer
	diagnostics []*analysis.Diagnostic
	err         error
}

// TODO(lotusirous): Do we need sequential execution ?
func execAll(actions []*action) {
	var wg sync.WaitGroup
	for _, act := range actions {
		wg.Add(1)
		work := func(act *action) {
			act.exec()
			wg.Done()
		}
		go work(act)
	}
	wg.Wait()
}

func (act *action) exec() {
	act.once.Do(func() {
		act.diagnostics, act.err = act.a.RunSingle(act.pass)
	})
}

func (act *action) String() string { return fmt.Sprintf("%s", act.a) }
