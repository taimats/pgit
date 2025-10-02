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
	"github.com/taimats/pgit/data"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "attach a name to a commit point that HEAD always refers to",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) == 1 {
			name = args[0]
		}
		if name == "" {
			current, err := currentBranchName()
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			list, err := ListBranches(string(current))
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			list[0] = fmt.Sprintf("*%s", list[0])
			str := strings.Join(list, "\n")
			fmt.Println(str)
			return nil
		}
		_, err := NewBranch(name)
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		return nil
	},
}

func NewBranch(name string) (path string, err error) {
	ref, err := data.NewRef(data.RefHEADPath)
	if err != nil {
		return "", fmt.Errorf("NewBranch: %w", err)
	}
	if ref.IsSymbolic {
		ref, err = ref.ResolveSymbolic(ref.Next)
		if err != nil {
			return "", fmt.Errorf("NewBranch: %w", err)
		}
	}
	path = filepath.Join(HeadDir, name)
	if err := data.WriteFile(path, []byte(ref.Oid)); err != nil {
		return "", fmt.Errorf("NewBranch: %w", err)
	}
	return path, nil
}

// a list of all the branches with the current one at the top
func ListBranches(currentBranch string) ([]string, error) {
	var fns []string
	fns = append(fns, currentBranch)
	err := filepath.WalkDir(data.RefBranchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == data.RefBranchPath {
			return nil
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		if d.Name() == currentBranch {
			return nil
		}
		fns = append(fns, d.Name())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fns, nil
}

// reads the HEAD file and returns a current branch name like this:
// [ ref: .pgit/ref/heads/{name} ] ===> name
func currentBranchName() (string, error) {
	current, err := data.ReadValueFromFile(data.RefHEADPath, []byte("ref:"))
	if err != nil {
		return "", fmt.Errorf("internal error: %w", err)
	}
	return filepath.Base(string(current)), nil
}

// returning all file names with any dir names excluded
func allFileNames(dirPath string) ([]string, error) {
	var fns []string
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == dirPath {
			return nil
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		fns = append(fns, d.Name())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fns, nil
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
