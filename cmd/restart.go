package cmd

import (
	"github.com/spf13/cobra"
)

// restartCmd represents the create command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart crud-mysql server",
	Run: func(cmd *cobra.Command, args []string) {
		//todo restart server

	},
}
