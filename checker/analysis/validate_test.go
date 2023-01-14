package analysis

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	var (
		run = func(p Pass) (interface{}, error) {
			return nil, nil
		}

		simple = &Analyzer{
			Name: "secretcheck",
			Meta: Meta{
				Severity:    "HIGH",
				Description: "foobar",
			},
			Run: run,
		}
	)

	cases := []struct {
		analyzers []*Analyzer
		wantErr   bool
	}{
		{
			[]*Analyzer{simple},
			false,
		},
	}

	for _, c := range cases {
		got := Validate(c.analyzers)
		if !c.wantErr {
			if got == nil {
				continue
			}
			t.Errorf("got unexpected error while validating analyzers %v: %v", c.analyzers, got)
		}

		if got == nil {
			t.Errorf("expected error while validating analyzers %v, but got nil", c.analyzers)
		}
	}

}

func TestValidateEmptyMeta(t *testing.T) {
	withoutDoc := &Analyzer{
		Name: "noMeta",
		Run: func(p Pass) (interface{}, error) {
			return nil, nil
		},
	}
	err := Validate([]*Analyzer{withoutDoc})
	if err == nil || !strings.Contains(err.Error(), "is undocumented") {
		t.Errorf("got unexpected error while validating analyzers withoutDoc: %v", err)
	}
}

func TestValidateNoRun(t *testing.T) {
	withoutRun := &Analyzer{
		Name: "withoutRun",
		Meta: Meta{
			Description: "No run",
			Severity:    "NO",
		},
	}
	err := Validate([]*Analyzer{withoutRun})
	if err == nil || !strings.Contains(err.Error(), "has nil Run") {
		t.Errorf("got unexpected error while validating analyzers withoutRun: %v", err)
	}
}
