/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
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
	head, err := getHead()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", ObjTypeTree, treeOid)
	if head != "" {
		fmt.Fprintf(&buf, "%s %s\n", "parent", head)
	}
	buf.WriteString("\n")
	buf.WriteString(msg)

	oid, err = SaveHashObj(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to save hash object: (error: %w)", err)
	}
	if err := updateHead(oid); err != nil {
		return "", err
	}
	return oid, nil
}

// HEAD represetns the latest commit, so the HEAD file always has one oid.
func updateHead(oid string) error {
	f, err := os.Create(filepath.Join(PgitDir, "HEAD"))
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(oid)
	return nil
}

func getHead() (string, error) {
	f, err := os.Open(filepath.Join(PgitDir, "HEAD"))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return "", nil
		} else {
			return "", err
		}
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var message string

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("message", "m", message, "add a message")
}
