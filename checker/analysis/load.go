package analysis

import (
	"io/fs"
	"path/filepath"
)

// Load loads the project into the pass to work with analyzer.
func Load(base string) (Pass, error) {
	files := []string{}

	err := filepath.Walk(base, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return Pass{}, err
	}

	return Pass{
		Base:  base,
		Files: files,
	}, nil
}
