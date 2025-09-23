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

	"github.com/spf13/cobra"
)

// readTreeCmd represents the readTree command
var readTreeCmd = &cobra.Command{
	Use:   "read-tree",
	Short: "lay out the content of a tree object into the working directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		oid := args[0]
		treeContent, err := FetchFileContent(oid)
		if err != nil {
			return fmt.Errorf("failed to find a tree content: (error: %s)", err)
		}
		if err := os.Mkdir(oid, os.ModeDir); err != nil {
			return err
		}
		sc := bufio.NewScanner(bytes.NewReader(treeContent))
		sc.Split(bufio.ScanLines)
		for sc.Scan() {
			line := sc.Bytes()
			sep := bytes.Split(line, []byte{' '}) //separated bytes hold "objType, oid, filename"
			if len(sep) != 3 {
				return fmt.Errorf("invalid data: { object: %s}", sep)
			}
			fc, err := FetchFileContent(string(sep[1]))
			if err != nil {
				return err
			}
			f, err := os.Create(filepath.Join(oid, string(sep[1])))
			if err != nil {
				return err
			}
			f.Write(fc)
			f.Close()
		}
		fmt.Println("read a tree object!!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(readTreeCmd)
}
