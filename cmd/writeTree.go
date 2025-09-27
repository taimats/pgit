/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "turn the current directory into a tree object and save it",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := CheckPgitInit()
		if err != nil {
			return err
		}
		rootPath := "."
		oid, err := saveTree(filepath.Clean(rootPath))
		if err != nil {
			return fmt.Errorf("failed to write tree: %w", err)
		}
		fmt.Printf("saved a tree!!\noid: %s\n", oid)
		return nil
	},
}

// saveTree is just a high-level layer of function to execute write-tree command.
func saveTree(rootPath string) (oid string, err error) {
	oid, err = writeTree(rootPath)
	if err != nil {
		return "", fmt.Errorf("saveTree: %w", err)
	}
	return oid, nil
}

// Tree represents a directory in this package's context.
// writeTree walks through the roothPath and do the following things for each file (or directory):
// ・save each file as a hashed object in storage
// ・if the given file is a directory, then recursively do the same
// ・at the end, save the whole directory (i.e. rootPath) as a hashed object in storage
func writeTree(rootPath string) (treeOid string, err error) {
	var buf bytes.Buffer
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}
		if isExcluded(d.Name()) {
			return filepath.SkipDir
		}
		if d.IsDir() {
			oid, err := writeTree(path)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(&buf, "%s %s %s\n", ObjTypeTree, oid, d.Name()); err != nil {
				return err
			}
			return nil
		}
		b, err := ReadAllFileContent(path)
		if err != nil {
			return fmt.Errorf("writeTree: %w", err)
		}
		oid, err := SaveHashObj(b)
		if err != nil {
			return fmt.Errorf("writeTree: %w", err)
		}
		fmt.Fprintf(&buf, "%s %s %s\n", ObjTypeBlob, oid, d.Name())
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("writeTree: %w", err)
	}
	treeOid, err = SaveHashObj(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("writeTree: %w", err)
	}
	return treeOid, nil
}

func isExcluded(baseName string) bool {
	return baseName == PgitDir || strings.HasPrefix(baseName, ".")
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
