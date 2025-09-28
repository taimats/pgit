/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
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
	headOid, err := resolveHEAD()
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", ObjTypeTree, treeOid)
	if headOid != "" {
		fmt.Fprintf(&buf, "%s %s\n", "parent", headOid)
	}
	buf.WriteString("\n")
	buf.WriteString(msg)

	commitOid, err = SaveHashObj(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	ref, err := peekHEADFile()
	if err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	if err := updateHEAD(ref, true, ""); err != nil {
		return "", fmt.Errorf("NewCommit: %w", err)
	}
	return commitOid, nil
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
	if err := WriteFile(filepath.Join(PgitDir, RefDir, HeadDir, ref), []byte(oid)); err != nil {
		return fmt.Errorf("updateRef: %w", err)
	}
	return nil
}

// If isSymbolic is true, updateHEAD saves a symbolic ref (= ref: refs/heads/{ref}) in the HEAD file.
// Otherwise, saves the physical oid directly.
func updateHEAD(ref string, isSymbolic bool, oid string) error {
	if !isSymbolic {
		if err := WriteFile(filepath.Join(PgitDir, RefHEAD), []byte(oid)); err != nil {
			return fmt.Errorf("updateHEAD: %w", err)
		}
		return nil
	}
	symbolic := filepath.Join(RefDir, HeadDir, ref)
	msg := fmt.Sprintf("ref: %s <- HEAD\n", symbolic)
	if err := WriteFile(filepath.Join(PgitDir, RefHEAD), []byte(msg)); err != nil {
		return fmt.Errorf("updateHEAD: %w", err)
	}
	return nil
}

func resolveRef(ref string) (oid string, err error) {
	if ref == "" || ref == "@" || ref == RefHEAD {
		c, err := resolveHEAD()
		if err != nil {
			return "", fmt.Errorf("resolveRef: %w", err)
		}
		return string(c), nil
	}
	c, err := ReadAllFileContent(filepath.Join(PgitDir, RefDir, HeadDir, ref))
	if err != nil {
		return "", fmt.Errorf("resolveRef: %w", err)
	}
	return string(c), nil
}

func resolveHEAD() (oid string, err error) {
	ref, err := peekHEADFile()
	if err != nil {
		return "", fmt.Errorf("resolveHEAD: %w", err)
	}
	if ref == "" {
		return "", nil
	}
	c, err := ReadAllFileContent(filepath.Join(PgitDir, RefDir, HeadDir, string(ref)))
	if err != nil {
		return "", fmt.Errorf("resolveHEAD: %w", err)
	}
	return string(c), nil
}

// extracting a symbolic ref path from the following content in the file:
// ref: refs/heads/{ref} <- HEAD
// If the returned ref is empty, peekHEADFile takes an alternative approach.
// That is, it just reads the whole content and returns it, which is probably an oid.
func peekHEADFile() (ref string, err error) {
	c, err := ReadValueFromFile(filepath.Join(PgitDir, RefHEAD), []byte("ref:"))
	if err != nil {
		return "", fmt.Errorf("peekHEADFile: %w", err)
	}
	if c == nil {
		c, err := ReadAllFileContent(filepath.Join(PgitDir, RefHEAD))
		if err != nil {
			return "", fmt.Errorf("peekHEADFile: %w", err)
		}
		return string(c), nil
	}
	return ref, nil
}

var message string

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("message", "m", message, "add a message")
}
