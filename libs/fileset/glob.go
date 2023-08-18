package fileset

import (
	"io/fs"
	"os"
	"path/filepath"
)

type GlobSet struct {
	root     string
	patterns []string
}

func NewGlobSet(root string, includes []string) (*GlobSet, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	for k := range includes {
		includes[k] = filepath.Join(absRoot, filepath.FromSlash(includes[k]))
	}
	return &GlobSet{absRoot, includes}, nil
}

// Return all files which matches defined glob patterns
func (s *GlobSet) All() ([]File, error) {
	files := make([]File, 0)
	for _, pattern := range s.patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return files, err
		}

		for _, match := range matches {
			matchRel, err := filepath.Rel(s.root, match)
			if err != nil {
				return files, err
			}

			stat, err := os.Stat(match)
			if err != nil {
				return files, err
			}
			files = append(files, File{fs.FileInfoToDirEntry(stat), match, matchRel})
		}
	}

	return files, nil
}
