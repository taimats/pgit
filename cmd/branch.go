/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "attach a name to a commit point that HEAD always refers to",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		_, err := NewBranch(name)
		if err != nil {
			return err
		}
		return nil
	},
}

func NewBranch(name string) (path string, err error) {
	headOid, err := resolveHEAD()
	if err != nil {
		return "", fmt.Errorf("NewBranch: %w", err)
	}
	path = filepath.Join(PgitDir, RefDir, HeadDir, name)
	if err := WriteFile(path, []byte(headOid)); err != nil {
		return "", fmt.Errorf("NewBranch: %w", err)
	}
	return path, nil
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
