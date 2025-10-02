/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "undoes a commit, making the current branch point to a specified oid",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		oid := args[0]
		ref, err := data.NewRef(data.RefHEADPath)
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		if ref.IsSymbolic {
			ref, err = ref.ResolveSymbolic(ref.Next)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
		}
		if err := ref.Update(oid); err != nil {
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
