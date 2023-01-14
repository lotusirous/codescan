package analysis

import (
	"fmt"
	"unicode"
)

// Validate reports an error if any of the analyzers are misconfigured.
func Validate(analyzers []*Analyzer) error {
	visit := func(a *Analyzer) error {
		if a == nil {
			return fmt.Errorf("nil Analyzer")
		}

		if !validName(a.Name) {
			return fmt.Errorf("invalid analyzer name %v", a)
		}

		if a.Meta.Description == "" {
			return fmt.Errorf("analyzer %v is undocumented", a)
		}

		if a.Meta.Severity == "" {
			return fmt.Errorf("analyzer %v has no level: HIGH, DANGER", a)
		}

		if a.Run == nil && a.RunSingle == nil {
			return fmt.Errorf("analyzer %v wants a Run or RunSingle", a)
		}

		return nil
	}

	for _, a := range analyzers {
		if err := visit(a); err != nil {
			return err
		}
	}

	return nil
}

func validName(name string) bool {
	for i, r := range name {
		if !(r == '_' || unicode.IsLetter(r) || i > 0 && unicode.IsDigit(r)) {
			return false
		}
	}
	return name != ""
}
