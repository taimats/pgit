/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
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
		oid, err := saveTree(".")
		if err != nil {
			return fmt.Errorf("failed to write tree: %w", err)
		}
		fmt.Printf("saved a tree!!\noid: %s\n", oid)
		return nil
	},
}

// saveTree is just a high-level layer of function to execute write-tree command.
func saveTree(rootPath string) (oid string, err error) {
	oid, err = data.WriteTree(rootPath, ObjDir)
	if err != nil {
		return "", fmt.Errorf("saveTree: %w", err)
	}
	return oid, nil
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
