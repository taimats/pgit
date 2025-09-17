package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var hashObjCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "save a hashed-object",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pgitDir := PgitDirPath()
		if pgitDir == "" {
			return errors.New("need initializing pgit beforehand")
		}
		rootDir := filepath.Dir(pgitDir)
		filename := args[0]
		path := filepath.Join(rootDir, filename)
		fmt.Println("path:", path)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		oid, err := hashObject(content)
		if err != nil {
			return err
		}
		log.Printf("saved a hashed-object!!\noid: %s\n", oid)
		return nil
	},
}

// oid is object ID.
func hashObject(data []byte) (oid string, err error) {
	s := sha1.Sum(data)
	oid = hex.EncodeToString(s[:])
	f, err := os.OpenFile(filepath.Join(RootDir, ObjDir, oid), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.WriteString(oid)
	return oid, nil
}

func init() {
	rootCmd.AddCommand(hashObjCmd)
}
