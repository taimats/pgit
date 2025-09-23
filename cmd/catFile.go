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
		objdir, err := AbsObjDirPath()
		if err != nil {
			return err
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
		obj, err := Decode(IdentBlob, content)
		if err != nil {
			return err
		}
		fmt.Println(string(obj.data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
