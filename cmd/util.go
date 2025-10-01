package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/taimats/pgit/data"
)

const (
	PgitDir = ".pgit"
	ObjDir  = ".pgit/objects"
	RefDir  = ".pgit/refs"
	TagDir  = ".pgit/refs/tags"
	HeadDir = ".pgit/refs/heads"

	HEAD = "HEAD"
)

var ErrNeedPgitInit = errors.New("need initializing pgit first")

func CheckPgitInit() error {
	if _, err := os.Stat(PgitDir); err != nil {
		return ErrNeedPgitInit
	}
	return nil
}

// converts content to a hashed-object under the hood, and
// save it as a file in the object storage (= .pgit/objects)
func SaveHashObj(content []byte) (oid string, err error) {
	oid, err = data.SaveBlobObj(ObjDir, content)
	if err != nil {
		return "", fmt.Errorf("SaveHashObj: %w", err)
	}
	return oid, nil
}
