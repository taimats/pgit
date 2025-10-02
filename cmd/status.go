/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "print a status of the current branch",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := currentBranchName()
		if err != nil {
			return fmt.Errorf("no such a branch: %w", err)
		}
		fmt.Printf("on branch %s\n", current)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
