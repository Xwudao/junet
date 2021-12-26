package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Xwudao/junet/cmd/junet/cmd/gen"
	"github.com/Xwudao/junet/cmd/junet/utils"
)

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen something",
	Run: func(cd *cobra.Command, args []string) {
		utils.Info("gen something")
	},
}

func init() {
	GenCmd.AddCommand((&gen.Route{}).Cmd())
	GenCmd.AddCommand((&gen.DB{}).Cmd())

	RootCmd.AddCommand(GenCmd)
}
