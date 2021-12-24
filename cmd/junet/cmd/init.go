package cmd

import (
	"github.com/spf13/cobra"
)

var initProject = &cobra.Command{
	Use:   "init",
	Short: "init junet project",
}

func init() {
	rootCmd.AddCommand(initProject)
}
