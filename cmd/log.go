/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var ErrNotFound = errors.New("not found")

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "print commit log list",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ref := RefHEAD
		if len(args) == 1 {
			ref = args[0]
		}
		err := CommitList(ref)
		if err != nil {
			return err
		}
		return nil
	},
}

// priting commit history, namely the list of commit oids
// starting from the provided ref to the initial commit
func CommitList(ref string) error {
	startOid, err := getOidFromRef(ref)
	if err != nil {
		return fmt.Errorf("CommitList func error: %w", err)
	}
	parent, err := commitParent(startOid)
	if err != nil {
		return fmt.Errorf("CommitList func error: %w", err)
	}
	if parent == "" {
		fmt.Println(startOid)
		return nil
	}
	var buf strings.Builder
	fmt.Fprintf(&buf, "%s\n", startOid)
	fmt.Fprintf(&buf, "%s\n", parent)
	for {
		parent, err = commitParent(parent)
		if err != nil {
			return err
		}
		if parent == "" {
			break
		}
		fmt.Fprintf(&buf, "%s\n", parent)
	}
	fmt.Println(buf.String())
	return nil
}

func commitParent(oid string) (parentOid string, err error) {
	if oid == "" {
		return "", fmt.Errorf("commitParent error: oid is empty")
	}
	c, err := ReadAllFileContent(filepath.Join(PgitDir, ObjDir, oid))
	if err != nil {
		return "", fmt.Errorf("commitParent func error: %w", err)
	}
	sc := bufio.NewScanner(bytes.NewReader(c))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		sep := strings.Split(sc.Text(), " ")
		if sep[0] == "parent" {
			return sep[1], nil
		}
	}
	return "", nil
}

func init() {
	rootCmd.AddCommand(logCmd)
}
