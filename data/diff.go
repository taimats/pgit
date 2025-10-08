package data

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Diff struct {
	Filename string
	Diff     string
}

// comparing two trees and generating the differences in a clear way.
func DiffTrees(from Tree, to Tree, srcDir string) ([]*Diff, error) {
	difs := make([]*Diff, 0, len(from)+len(to))
	for name, fromElem := range from {
		toElem, ok := to[name]
		if !ok {
			continue
		}
		diff, err := DiffFiles(filepath.Join(srcDir, fromElem.Oid), filepath.Join(srcDir, toElem.Oid))
		if err != nil {
			return nil, fmt.Errorf("DiffTrees: %w", err)
		}
		if diff == "" {
			continue
		}
		difs = append(difs, &Diff{Filename: name, Diff: diff})
	}
	return difs, nil
}

// comparing the content of files between fromPath and toPath, and generating an output of differences
func DiffFiles(fromPath string, toPath string) (diff string, err error) {
	from, err := ReadAllFileContent(fromPath)
	if err != nil {
		return "", fmt.Errorf("DiffFiles: %w", err)
	}
	to, err := ReadAllFileContent(toPath)
	if err != nil {
		return "", fmt.Errorf("DiffFiles: %w", err)
	}
	dmp := diffmatchpatch.New()
	fromChars, toChars, list := dmp.DiffLinesToChars(string(from), string(to))
	diffs := dmp.DiffMain(fromChars, toChars, false)
	diffs = dmp.DiffCharsToLines(diffs, list)

	return diffReport(diffs), nil
}

// converting multiple diffs into a human-readable line-by-line report in the following way:
// +this is a example text.
// -this is a example test.
// Note:
// +(plus) represents inserted strings while -(minus) indicates deleted strings.
func diffReport(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			lines := strings.Split(text, "\n")
			for _, line := range lines {
				fmt.Fprintf(&buff, "+%s\n", line)
			}
		case diffmatchpatch.DiffDelete:
			lines := strings.Split(text, "\n")
			for _, line := range lines {
				fmt.Fprintf(&buff, "-%s\n", line)
			}
		}
	}
	return buff.String()
}
