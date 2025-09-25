/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "attach a name to an oid",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var ref, oid string
		if len(args) == 2 {
			ref, oid = args[0], args[1]
		} else {
			ref = args[0]
		}
		if err := updateRef(ref, oid); err != nil {
			return err
		}
		fmt.Println("created a tag!!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
