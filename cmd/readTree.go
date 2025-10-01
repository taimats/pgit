/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

// readTreeCmd represents the readTree command
var readTreeCmd = &cobra.Command{
	Use:   "read-tree",
	Short: "lay out the content of a tree object into the working directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sweepDir("."); err != nil {
			return err
		}
		oid := args[0]
		if err := data.ReadTree(oid, ObjDir, "."); err != nil {
			return err
		}
		fmt.Println("read a tree object!!")
		return nil
	},
}

func sweepDir(rootPath string) error {
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}
		if isIgnored(d.Name()) {
			return filepath.SkipDir
		}
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("sweepDir: %w", err)
	}
	return nil
}

func isIgnored(baseName string) bool {
	return baseName == PgitDir || strings.HasPrefix(baseName, ".")
}

func init() {
	rootCmd.AddCommand(readTreeCmd)
}
