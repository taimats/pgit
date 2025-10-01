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
	"github.com/taimats/pgit/data"
)

var ErrNotFound = errors.New("not found")

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "print commit log list",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) == 1 {
			name = args[0]
		}
		var refPath string
		if name == data.HEADAlias || name == HEAD || name == "" {
			refPath = filepath.Join(data.RefHEADPath)
		} else {
			refPath = filepath.Join(data.RefBranchPath, name)
		}
		ref, err := data.NewRef(refPath)
		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}
		if ref.IsSymbolic {
			ref, err = ref.ResolveSymbolic(ref.Next)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
		}
		var buf strings.Builder
		current := ref.Oid
		for {
			fmt.Fprintf(&buf, "%s", current)
			oid, err := commitParent(current)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			if oid == "" {
				break
			}
			current = oid
		}
		fmt.Println(buf.String())
		return nil
	},
}

func commitParent(oid string) (parentOid string, err error) {
	c, err := data.ReadAllFileContent(filepath.Join(ObjDir, oid))
	if err != nil {
		return "", fmt.Errorf("commitParent: %w", err)
	}
	sc := bufio.NewScanner(bytes.NewReader(c))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		sep := strings.Split(sc.Text(), " ")
		if len(sep) < 2 {
			continue
		}
		if sep[0] == "parent" {
			return sep[1], nil
		}
	}
	return "", nil
}

func init() {
	rootCmd.AddCommand(logCmd)
}
