/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
		err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
			if path == cwd {
				return nil
			}
			if d.Name() == PgitDir {
				return nil
			}
			fmt.Println(d.Name())
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
