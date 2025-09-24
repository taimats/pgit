/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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

func NewCommit(msg string) (oid string, err error) {
	treeOid, err := writeTree(".")
	if err != nil {
		return "", fmt.Errorf("failed to write tree: (error: %w)", err)
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", ObjTypeTree, treeOid)
	buf.WriteString("\n")
	buf.WriteString(msg)

	oid, err = SaveHashObj(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to save hash object: (error: %w)", err)
	}
	if err := addHeadFile(oid); err != nil {
		return "", err
	}
	return oid, nil
}

func addHeadFile(oid string) error {
	f, err := os.OpenFile(
		filepath.Join(PgitDir, "HEAD"),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0660,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", oid)
	return nil
}

var message string

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("message", "m", message, "add a message")
}
