package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pgit",
	Short: "psuedo git command",
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd, hashObjCmd)
}
