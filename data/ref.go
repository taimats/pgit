package data

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const (
	PgitDirBase = ".pgit"
	ObjDirBase  = "objects"
	RefDirBase  = "refs"
	HeadDirBase = "heads"
	TagDirBase  = "tags"

	HEAD      = "HEAD"
	HEADAlias = "@"
)

var (
	RefBranchPath = filepath.Join(PgitDirBase, RefDirBase, HeadDirBase)
	RefTagPath    = filepath.Join(PgitDirBase, RefDirBase, TagDirBase)
	RefHEADPath   = filepath.Join(PgitDirBase, HEAD)
)

// Ref is a shorhand for reference, and its main feature is to
// generalize a file reference. It reads and writes an oid or a symbolic ref in a referenced file.
type Ref struct {
	Path       string //referenced file path
	Oid        string //object id written in the referenced file, if any
	IsSymbolic bool   //reports whether this ref returns an oid or another ref
	Next       string //pointing to a symbolic ref, if any
}

// Note: If a returned Ref is nil, that means there is no such a file reffered by the path.
// When a ref with a specified path is needed, the file should be generated beforehand in the path.
func NewRef(path string) (*Ref, error) {
	ref := &Ref{Path: path, IsSymbolic: false}
	c, err := ReadAllFileContent(path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil, nil
		} else {
			return nil, fmt.Errorf("NewRef: %w", err)
		}
	}
	if isSymbolic(c) {
		symbolic := getSymbolicRefPath(c)
		ref.IsSymbolic = true
		ref.Next = string(symbolic)
		return ref, nil
	}
	ref.Oid = string(c)
	return ref, nil
}

func (r *Ref) Update(oid string) error {
	if !r.IsSymbolic {
		if err := WriteFile(r.Path, []byte(oid)); err != nil {
			return fmt.Errorf("Ref Update: %w", err)
		}
		return nil
	}
	ref, err := r.ResolveSymbolic(r.Next)
	if err != nil {
		return fmt.Errorf("Ref Update: %w", err)
	}
	if err := WriteFile(ref.Path, []byte(oid)); err != nil {
		return fmt.Errorf("Ref Update: %w", err)
	}
	return nil
}

func (r *Ref) ResolveSymbolic(next string) (*Ref, error) {
	if !r.IsSymbolic {
		return nil, nil
	}
	var resolved *Ref
	current := r
	for {
		ref, err := current.ResolveSymbolic(current.Path)
		if err != nil {
			return nil, fmt.Errorf("ResolveSymbolic: %w", err)
		}
		if !ref.IsSymbolic {
			resolved = ref
			break
		}
		current = ref
	}
	return resolved, nil
}

func isSymbolic(b []byte) bool {
	return bytes.HasPrefix(b, []byte("ref:"))
}

func getSymbolicRefPath(b []byte) []byte {
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := sc.Bytes()
		sep := bytes.Split(line, []byte(" "))
		if len(sep) > 2 && bytes.Equal(sep[0], []byte("ref:")) {
			return sep[1]
		}
	}
	return nil
}
