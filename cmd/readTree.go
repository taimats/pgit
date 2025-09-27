/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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
		treeContent, err := FetchFileContent(oid)
		if err != nil {
			return fmt.Errorf("no tree content: %w", err)
		}
		sc := bufio.NewScanner(bytes.NewReader(treeContent))
		sc.Split(bufio.ScanLines)
		for sc.Scan() {
			line := sc.Bytes()
			sep := bytes.Split(line, []byte{' '}) //separated bytes hold "objType, oid, filename"
			if len(sep) != 3 {
				return fmt.Errorf("invalid data: { object: %s}", sep)
			}
			_, oid, filename := sep[0], sep[1], sep[2]
			fc, err := FetchFileContent(string(oid))
			if err != nil {
				return err
			}
			f, err := os.Create(filepath.Clean(string(filename)))
			if err != nil {
				return err
			}
			f.Write(fc)
			f.Close()
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
