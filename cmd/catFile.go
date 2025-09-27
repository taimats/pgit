package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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
		content, err := FetchFileContent(oid)
		if err != nil {
			return fmt.Errorf("failed to fetch file content: (error: %w)", err)
		}
		fmt.Println(string(content))
		return nil
	},
}

// search the path /.pgit/objects/{oid} for the content of a file
func FetchFileContent(oid string) ([]byte, error) {
	f, err := os.Open(filepath.Join(PgitDir, ObjDir, oid))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil, fmt.Errorf("FetchFileContent: no such an oid")
		} else {
			return nil, fmt.Errorf("FetchFileContent: %w", err)
		}
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("FetchFileContent: %w", err)
	}
	return b, nil
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
