package cmd

import (
	"errors"
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
