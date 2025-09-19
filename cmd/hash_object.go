package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
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
		objdir, err := AbsObjDirPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(objdir); err != nil {
			return errors.New("need initializing first")
		}
		filename := args[0]
		path, err := filepath.Abs(filename)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		oid := hashObject(content)
		if err != nil {
			return err
		}
		newfile, err := os.Create(filepath.Join(objdir, oid))
		if err != nil {
			return err
		}
		defer newfile.Close()
		newfile.Write(content)

		log.Printf("saved a hashed-object!!\noid: %s\n", oid)
		return nil
	},
}

// oid is an so-called object ID.
func hashObject(data []byte) (oid string) {
	s := sha1.Sum(data)
	oid = hex.EncodeToString(s[:])
	return oid
}

func init() {
	rootCmd.AddCommand(hashObjCmd)
}
