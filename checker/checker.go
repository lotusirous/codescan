package checker

import (
	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/checker/multirunner"
	"github.com/lotusirous/codescan/checker/rules/cryptoleak"
)

// DefaultRules is the simple rules for detect potential vulnerability.
var allRules = map[string]*analysis.Analyzer{
	"G402": cryptoleak.Analyzer,
}

// DefaultRules defines the defaults rule of system.
func DefaultRules() []*analysis.Analyzer {
	var out []*analysis.Analyzer
	for k, v := range allRules {
		v.Name = k
		out = append(out, v)
	}
	return out
}

// Run starts the group of analyzers to analyze a directory.
func Run(dir string, analyzers []*analysis.Analyzer) ([]*analysis.Diagnostic, error) {
	return multirunner.Run(dir, analyzers)
}
