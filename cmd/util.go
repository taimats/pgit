package cmd

import (
	"errors"
	"os"
)

const PgitDir = ".pgit"
const ObjDir = "objects"

var ErrNeedPgitInit = errors.New("need initializing pgit first")

func CheckPgitInit() error {
	if _, err := os.Stat(PgitDir); err != nil {
		return ErrNeedPgitInit
	}
	return nil
}
