package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var RooDir string = "."

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "starting a pgit project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			RooDir = args[0]
		}
		if _, err := os.Stat(filepath.Join(RooDir, PgitDir)); err == nil {
			return errors.New("already initialized")
		}
		if err := os.MkdirAll(filepath.Join(RooDir, PgitDir, ObjDir), os.ModeDir); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(RooDir, PgitDir, RefDir, TagDir), os.ModeDir); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(RooDir, PgitDir, RefDir, HeadDir), os.ModeDir); err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(RooDir, PgitDir, RefHEAD))
		if err != nil {
			return err
		}
		f.Close()
		log.Println("starting a pgit project!!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
