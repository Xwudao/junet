package cmd

import (
	"github.com/Xwudao/junet/cmd/junet/cmd/goose"
	"github.com/spf13/cobra"

	"github.com/Xwudao/junet/cmd/junet/utils"
)

var GooseCmd = &cobra.Command{
	Use:   "goose",
	Short: "goose is a tool for migration for your database.",
	Run: func(cd *cobra.Command, args []string) {
		utils.Info("goose tools wrapper of up and down migration.")
	},
}

func init() {
	GooseCmd.AddCommand((&goose.Up{}).Cmd())
	GooseCmd.AddCommand((&goose.Down{}).Cmd())

	RootCmd.AddCommand(GooseCmd)
}
