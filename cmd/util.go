package cmd

import "path/filepath"

const PgitDir = ".pgit"
const ObjDir = "objects"

func AbsObjDirPath() (string, error) {
	name := filepath.Join(PgitDir, ObjDir)
	path, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	return path, nil
}
