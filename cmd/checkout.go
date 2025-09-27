/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "gets back to the specified commit point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commitOid := args[0]
		treeOid, err := ExtractCommitTree(commitOid)
		if err != nil {
			return err
		}
		if treeOid == "" {
			return fmt.Errorf("no tree oid")
		}
		if err := readTreeCmd.RunE(readTreeCmd, []string{treeOid}); err != nil {
			return err
		}
		if err := updateRef(RefHEAD, commitOid); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
