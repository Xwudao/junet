package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "",
	Short: "junet help cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("hello, use %s -h for more help\n", os.Args[0])
	},
}

func Execute() error {
	return RootCmd.Execute()
}
