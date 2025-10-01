/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "create a commit object",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckPgitInit(); err != nil {
			return err
		}
		oid, err := NewCommit(message)
		if err != nil {
			return err
		}
		fmt.Println("create a commit!!: ", oid)
		return nil
	},
}

func NewCommit(msg string) (commitOid string, err error) {
	treeOid, err := data.WriteTree(".", ObjDir)
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	ref, err := data.NewRef(data.RefHEADPath)
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	if ref.IsSymbolic {
		ref, err = ref.ResolveSymbolic(ref.Next)
		if err != nil {
			return "", fmt.Errorf("NewCommit: %w", err)
		}
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", ObjTypeTree, treeOid)
	if ref.Oid != "" {
		fmt.Fprintf(&buf, "%s %s\n", "parent", ref.Oid)
	}
	buf.WriteString("\n")
	buf.WriteString(msg)

	commitOid, err = SaveHashObj(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	if err := ref.Update(commitOid); err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	return commitOid, nil
}

var message string

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("message", "m", message, "add a message")
}
