package cmd

import (
	"errors"
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
		objdir, err := AbsObjDirPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(objdir); err != nil {
			return errors.New("need initializing first")
		}
		filename := args[0]
		f, err := os.Open(filepath.Join(objdir, filename))
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		fmt.Println(string(content))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
