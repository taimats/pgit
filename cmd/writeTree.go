/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "turn a directory into a tree object and save it",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := CheckPgitInit()
		if err != nil {
			return err
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		oid, err := saveTree(cwd)
		if err != nil {
			return fmt.Errorf("failed to write tree: %w", err)
		}
		fmt.Printf("saved a tree!!\noid: %s\n", oid)
		return nil
	},
}

func saveTree(rootPath string) (oid string, err error) {
	objdir, err := AbsObjDirPath()
	if err != nil {
		return "", err
	}
	oid, err = writeTree(objdir, rootPath)
	if err != nil {
		return "", err
	}
	return oid, nil
}

func writeTree(storePath, rootPath string) (oid string, err error) {
	var buf bytes.Buffer
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}
		if isExcluded(d.Name()) {
			return nil
		}
		if d.IsDir() && !isExcluded(d.Name()) {
			oid, err := writeTree(storePath, path)
			if err != nil {
				return fmt.Errorf("failed to handle dir tree: %w\n{ path: %s }", err, path)
			}
			if _, err := fmt.Fprintf(&buf, "%s %s %s\n", "tree", oid, d.Name()); err != nil {
				return err
			}
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open a file: %w", err)
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read a file: %w", err)
		}
		oid, err := SaveHashObj(storePath, b)
		if err != nil {
			return err
		}
		fmt.Fprintf(&buf, "%s %s %s\n", ObjTypeBlob, oid, d.Name())
		return nil
	})
	if err != nil {
		return "", err
	}
	oid, err = SaveHashObj(storePath, buf.Bytes())
	if err != nil {
		return "", err
	}
	return oid, nil
}

func isExcluded(baseName string) bool {
	return baseName == PgitDir || baseName == ObjDir || strings.HasPrefix(baseName, ".")
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
