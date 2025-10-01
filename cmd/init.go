package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
)

var RootDir string = "."

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "starting a pgit project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			RootDir = filepath.Join(RootDir, args[0])
		}
		if _, err := os.Stat(PgitDir); err == nil {
			return errors.New("already initialized")
		}
		if err := os.MkdirAll(filepath.Join(RootDir, ObjDir), os.ModeDir); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(RootDir, TagDir), os.ModeDir); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(RootDir, HeadDir), os.ModeDir); err != nil {
			return err
		}
		refPath := filepath.Join(RootDir, HeadDir, "master")
		if err := data.WriteFile(refPath, []byte("")); err != nil {
			return err
		}
		content := fmt.Sprintf("ref: %s <- HEAD\n", refPath)
		if err := data.WriteFile(data.RefHEADPath, []byte(content)); err != nil {
			return err
		}
		log.Println("starting a pgit project!!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
