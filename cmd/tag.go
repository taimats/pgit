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

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "attach a name to an oid",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, oid string
		if len(args) == 2 {
			name, oid = args[0], args[1]
		} else {
			name = args[0]
		}
		if oid == "" {
			head, err := data.NewRef(data.RefHEADPath)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			if head.IsSymbolic {
				head, err = head.ResolveSymbolic(head.Next)
				if err != nil {
					return fmt.Errorf("internal error: %w", err)
				}
			}
			oid = head.Oid
		}
		if err := data.WriteFile(filepath.Join(TagDir, name), []byte(oid)); err != nil {
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
		}
		fmt.Println("created a tag!!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
