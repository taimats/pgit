/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "print commit details",
	RunE: func(cmd *cobra.Command, args []string) error {
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
		c, err := data.GetCommit(ref.Oid)
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		fromTree, err := data.ParseTreeFile(filepath.Join(ObjDir, c.Parent))
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		toTree, err := data.ParseTreeFile(filepath.Join(ObjDir, c.TreeOid))
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		diffs, err := data.DiffTrees(fromTree, toTree, ObjDir)
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "commit %s\n", c.TreeOid)
		fmt.Fprintln(&buf, "")
		fmt.Fprintf(&buf, "%s\n", c.Msg)
		fmt.Fprintln(&buf, "")
		if len(diffs) == 0 {
			fmt.Fprintln(&buf, "No diffs right now!")
			fmt.Println(buf.String())
			return nil
		}
		for _, diff := range diffs {
			fmt.Fprintf(&buf, "%s\n%s\n", diff.Filename, diff.Diff)
		}
		fmt.Println(buf.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
