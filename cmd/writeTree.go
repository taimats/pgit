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

	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "save directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == cwd {
				return nil
			}
			if d.Name() == PgitDir {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			oid, err := SaveHashObj(cwd, b)
			if err != nil {
				return err
			}
			if d.IsDir() {
				fmt.Fprintf(&buf, "%s %s %s\n", "tree", oid, d.Name())
			} else {
				fmt.Fprintf(&buf, "%s %s %s\n", ObjTypeBlob, oid, d.Name())
			}
			return nil
		})
		if err != nil {
			return err
		}
		oid, err := SaveHashObj(cwd, buf.Bytes())
		if err != nil {
			return err
		}
		fmt.Printf("saved a tree!!\noid: %s\n", oid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
