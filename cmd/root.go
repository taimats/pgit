package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pgit",
	Short: "psuedo git command",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOut(os.Stderr)
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd, hashObjCmd)
}
