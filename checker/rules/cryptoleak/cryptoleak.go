package cryptoleak

import (
	"bufio"
	"os"
	"strings"

	"github.com/lotusirous/codescan/checker/analysis"
)

const Doc = `find private_key or private info`

var secretPatterns = []string{"private_key", "public_key"}

var Analyzer = &analysis.Analyzer{
	Name: "cryptoleak",
	Meta: analysis.Meta{
		Description: "Leak the cryptography keys",
		Severity:    "HIGH",
	},
	Run: run,
}

func run(pass analysis.Pass) ([]*analysis.Diagnostic, error) {
	var out []*analysis.Diagnostic

	for _, file := range pass.Files {
		diag, err := scanFile(file)
		if err != nil {
			return nil, err
		}
		out = append(out, diag...)
	}
	return out, nil
}

func scanFile(path string) ([]*analysis.Diagnostic, error) {
	var out []*analysis.Diagnostic
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	lineNum := 1

	for sc.Scan() {
		line := sc.Text()
		words := strings.Split(line, " ")
		for _, w := range words {
			for _, pattern := range secretPatterns {
				if strings.HasPrefix(w, pattern) {
					out = append(out, &analysis.Diagnostic{
						Pos:     lineNum,
						Path:    path,
						Message: "Secret key might be leaked",
					})
				}
			}
		}
		lineNum++
	}
	return out, nil
}
