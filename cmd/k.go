/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// kCmd represents the k command
var kCmd = &cobra.Command{
	Use:   "k",
	Short: "visualizing refs and relavant oids",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootPath := filepath.Join(PgitDir, RefDir, TagDir)
		var buf strings.Builder
		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == rootPath {
				return nil
			}
			fmt.Fprintf(&buf, "%s\n", d.Name())
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Println(buf.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(kCmd)
}
