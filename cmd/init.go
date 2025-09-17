package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const RootDir = ".pgit"
const ObjDir = "objects"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "starting a pgit project",
	RunE: func(cmd *cobra.Command, args []string) error {
		pgitDir := PgitDirPath()
		if pgitDir != "" {
			return errors.New("pgit is already initialiezed")
		}
		if err := os.MkdirAll(filepath.Join(pgitDir, ObjDir), 0400); err != nil {
			return err
		}
		log.Println("starting a pgit project!!")
		return nil
	},
}

// The returned path is an absolute one including ".pgit" directory.
func PgitDirPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	path := filepath.Join(cwd, RootDir)
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func init() {
	rootCmd.AddCommand(initCmd)
}
