package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "starting a pgit project",
	RunE: func(cmd *cobra.Command, args []string) error {
		objdir, err := AbsObjDirPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(objdir); err == nil {
			return errors.New("already initialized")
		}
		if err := os.MkdirAll(objdir, os.ModeDir); err != nil {
			return err
		}
		log.Println("starting a pgit project!!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
