/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var ErrNotFound = errors.New("not found")

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "print commit log list",
	RunE: func(cmd *cobra.Command, args []string) error {
		var start string
		if len(args) > 0 {
			start = args[0]
		}
		err := CommitList(start)
		if err != nil {
			return err
		}
		return nil
	},
}

// priting commit history, the list of commit oid starting from the startOid to the initial commit
func CommitList(startOid string) error {
	var err error
	if startOid == "" {
		startOid, err = getOidFromRef(RefHEAD)
		if err != nil {
			return fmt.Errorf("failed to get an oid for HEAD: %w", err)
		}
	}
	parent, err := commitParent(startOid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println(startOid)
			return nil
		} else {
			return err
		}
	}
	var buf strings.Builder
	fmt.Fprintf(&buf, "%s\n", startOid)
	fmt.Fprintf(&buf, "%s\n", parent)
	for {
		parent, err = commitParent(parent)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				break
			} else {
				return err
			}
		}
		fmt.Fprintf(&buf, "%s\n", parent)
	}
	fmt.Println(buf.String())
	return nil
}

func commitParent(oid string) (parentOid string, err error) {
	f, err := os.Open(filepath.Join(PgitDir, ObjDir, oid))
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		sep := strings.Split(sc.Text(), " ")
		if sep[0] == "parent" {
			return sep[1], nil
		}
	}
	return "", ErrNotFound
}

func init() {
	rootCmd.AddCommand(logCmd)
}
