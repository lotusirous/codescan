package multirunner

import (
	"sync"

	"github.com/lotusirous/codescan/checker/analysis"
)

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
		actions = append(actions, &action{analyzer: a, pass: pass})
	}
	// Execute in parallel.
	execAll(actions)
	return actions
}

func collectDiagnostics(actions []*action) ([]*analysis.Diagnostic, error) {
	var out []*analysis.Diagnostic
	for _, act := range actions {
		if act.err != nil {
			return nil, act.err
		}

		for _, diag := range act.diagnostics {
			diag.ByAnalyzer = act.analyzer
			out = append(out, diag)
		}

	}
	return out, nil
}

// An action represents one unit of analysis work.
type action struct {
	once        sync.Once
	pass        analysis.Pass
	analyzer    *analysis.Analyzer
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
		act.diagnostics, act.err = act.analyzer.Run(act.pass)
	})
}

func (act *action) String() string { return act.analyzer.Name }
