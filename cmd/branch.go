/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
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
	headOid, err := getOidFromRef(RefHEAD)
	if err != nil {
		return "", fmt.Errorf("NewBranch: %w", err)
	}
	path = filepath.Join(PgitDir, RefDir, HeadDir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("NewBranch: %w", err)
	}
	defer f.Close()

	f.WriteString(headOid)
	return path, nil
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
