package cmd

import (
	"errors"
	"os"
	"path/filepath"
)

const PgitDir = ".pgit"
const ObjDir = "objects"

var ErrNeedPgitInit = errors.New("need initializing pgit first")

func AbsObjDirPath() (string, error) {
	name := filepath.Join(PgitDir, ObjDir)
	path, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	return path, nil
}

func CheckPgitInit() error {
	objdir, err := AbsObjDirPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(objdir); err != nil {
		return ErrNeedPgitInit
	}
	return nil
}
