package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

type FileCompleter struct{}

func (f *FileCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// Convert to string for easier manipulation
	s := string(line[:pos])

	// Determine the prefix to complete
	segments := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == '|'
	})
	prefix := ""
	if len(segments) > 0 {
		prefix = segments[len(segments)-1]
	}

	// Determine the directory to read from
	dir := "."
	if strings.Contains(prefix, "/") {
		dir, prefix = path.Split(prefix)
	}

	// Resolve ".." and "."
	dir = filepath.Clean(dir)

	// Read the specified directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return [][]rune{}, 0
	}

	// Find files that match the prefix
	var matches [][]rune
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, "./")

			match := path.Join(dir, name)[len(prefix)+len(dir)-1:]
			match = strings.TrimPrefix(match, "./")
			matches = append(matches, []rune(match))
		}
	}

	// Return matches and their common prefix length
	return matches, len(prefix)
}
