/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
		return "", fmt.Errorf("NewCommit failed to write tree: %w", err)
	}
	head, err := getHeadOid()
	if err != nil {
		return "", fmt.Errorf("NewCommit failed to getHeadOid: %w", err)
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
		return "", fmt.Errorf("failed to save hash object: (error: %w)", err)
	}
	if err := updateRef(RefHEAD, oid); err != nil {
		return "", err
	}
	return oid, nil
}

func ExtractCommitTree(commitOid string) (treeOid string, err error) {
	b, err := FetchFileContent(commitOid)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file content: %w", err)
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		sep := strings.Split(sc.Text(), " ")
		if sep[0] == "tree" {
			return sep[1], nil
		}
	}
	return "", ErrNotFound
}

// Ref is a shorthand for reference, and it is responsible for attaching a name to a specific oid (object ID).
// Ref is classified in two ways: tag or branch. Tag represents a commit oid and branch a HEAD alias.
// The real stuff of tag and branch is just a file in each directory, /refs/tags/{file} and /refs/heads/{file}.
func updateRef(ref string, oid string) error {
	if oid == "" {
		headOid, err := getOidFromRef(RefHEAD)
		if err != nil {
			return err
		}
		oid = headOid
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
	if ref == "" || ref == "@" {
		ref = RefHEAD
	}
	f, err := os.Open(filepath.Join(PgitDir, RefDir, TagDir, ref))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return "", ErrNotFound
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

func getHeadOid() (oid string, err error) {
	c, err := ReadAllFileContent(filepath.Join(PgitDir, RefHEAD))
	if err != nil {
		return "", fmt.Errorf("getHeadOid func error: %w", err)
	}
	return string(c), nil
}

var message string

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("message", "m", message, "add a message")
}
