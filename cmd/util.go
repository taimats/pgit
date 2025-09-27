package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	PgitDir = ".pgit"
	ObjDir  = "objects"
	RefDir  = "refs"
	TagDir  = "tags"
	HeadDir = "heads"
)

const (
	RefHEAD = "HEAD"
)

var ErrNeedPgitInit = errors.New("need initializing pgit first")

func CheckPgitInit() error {
	if _, err := os.Stat(PgitDir); err != nil {
		return ErrNeedPgitInit
	}
	return nil
}

// Return all the content of a file
func ReadAllFileContent(path string) (content []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadAllFileContent func error: %w", err)
	}
	defer f.Close()
	content, err = io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("ReadAllFileContent func error: %w", err)
	}
	return content, nil
}
