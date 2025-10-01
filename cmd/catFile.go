package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "print the content of a named file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := CheckPgitInit()
		if err != nil {
			return err
		}
		oid := args[0]
		content, err := data.ReadAllFileContent(filepath.Join(ObjDir, oid))
		if err != nil {
			return fmt.Errorf("failed to fetch file content: (error: %w)", err)
		}
		fmt.Println(string(content))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
