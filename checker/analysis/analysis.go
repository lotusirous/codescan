package analysis

// Meta represents a description of meta data for the rules.
type Meta struct {
	Description string
	Severity    string // HIGH, MEDIUM, LOW
}

// A Pass provides information to the Run function
// that applies a specific analyzer.
type Pass struct {
	Base  string
	Files []string // file path from base, ex: src/a/main.go
	// you might store the AST here, ex: ast.File
}

// An Analyzer describes an analysis function and its options.
type Analyzer struct {
	Name string
	Meta Meta

	// If we want to resolve the dependency between analyzer the returns value should be
	Run func(Pass) (any, error)

	// RunSingle allows a single analyzer to report after execution.
	// it is different from the Run that supports a chain execution.
	RunSingle func(Pass) ([]*Diagnostic, error)
}

func (a *Analyzer) String() string { return a.Name }

// A Diagnostic is a message associated with a source location or range.
//
// An Analyzer may return a variety of diagnostics;
// It is primarily intended to make it easy to look up documentation.
//The Report function emits a diagnostic, a message associated with a
// source position. For most analyses, diagnostics are their primary
// result.

type Diagnostic struct {
	ByAnalyzer *Analyzer // detected by analyzer.
	Path       string    // location to the file from base.
	Pos        int
	Message    string
	// End      int    // optional
	// Category should be a constant, may be used to classify them.
	// Category string // optional
}
