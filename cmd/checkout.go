/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "gets back to the specified commit point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refBranch := filepath.Join(data.RefBranchPath, args[0])
		ref, err := data.NewRef(refBranch)
		if err != nil {
			return fmt.Errorf("invalid ref name: %w", err)
		}
		if ref.IsSymbolic {
			ref, err = ref.ResolveSymbolic(ref.Next)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
		}
		treeOid, err := data.ReadValueFromFile(filepath.Join(ObjDir, ref.Oid), []byte("tree"))
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		if err := data.ReadTree(string(treeOid), ObjDir, "."); err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		head, err := data.NewRef(data.RefHEADPath)
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		if err := head.UpdateSymbolic(refBranch); err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
