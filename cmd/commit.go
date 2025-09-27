/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func NewCommit(msg string) (commitOid string, err error) {
	treeOid, err := writeTree(".")
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	head, err := getHeadOid()
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", ObjTypeTree, treeOid)
	if head != "" {
		fmt.Fprintf(&buf, "%s %s\n", "parent", head)
	}
	buf.WriteString("\n")
	buf.WriteString(msg)

	oid, err := SaveHashObj(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	if err := updateRef(RefHEAD, oid); err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	return oid, nil
}

func ExtractCommitTree(commitOid string) (treeOid string, err error) {
	b, err := FetchFileContent(commitOid)
	if err != nil {
		return "", fmt.Errorf("ExtractCommitTree: %w", err)
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		sep := strings.Split(sc.Text(), " ")
		if sep[0] == "tree" {
			return sep[1], nil
		}
	}
	return "", nil
}

// Ref is a shorthand for reference, and it is responsible for attaching a name to a specific oid (object ID).
// Ref is classified in two ways: tag or branch. Tag represents a commit oid and branch a HEAD alias.
// The real stuff of tag and branch is just a file in each directory, /refs/tags/{file} and /refs/heads/{file}.
func updateRef(ref string, oid string) error {
	if ref == RefHEAD || ref == "" {
		f, err := os.Create(filepath.Join(PgitDir, RefHEAD))
		if err != nil {
			return fmt.Errorf("updateRef: os.Create %w", err)
		}
		defer f.Close()
		f.WriteString(oid)
		return nil
	}
	f, err := os.Create(filepath.Join(PgitDir, RefDir, TagDir, ref))
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(oid)
	return nil
}

func getOidFromRef(ref string) (oid string, err error) {
	if ref == "" || ref == "@" || ref == RefHEAD {
		c, err := getHeadOid()
		if err != nil {
			return "", fmt.Errorf("getOidFromRef: %w", err)
		}
		return string(c), nil
	}
	c, err := ReadAllFileContent(filepath.Join(PgitDir, RefDir, TagDir, ref))
	if err != nil {
		return "", fmt.Errorf("getOidFromRef: %w", err)
	}
	return string(c), nil
}

func getHeadOid() (oid string, err error) {
	c, err := ReadAllFileContent(filepath.Join(PgitDir, RefHEAD))
	if err != nil {
		return "", fmt.Errorf("getHeadOid: %w", err)
	}
	return string(c), nil
}

var message string

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("message", "m", message, "add a message")
}
