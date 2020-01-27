package cmd

import (
	"github.com/spf13/cobra"
)

// stopCmd represents the create command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop crud-mysql server",
	Run: func(cmd *cobra.Command, args []string) {
		//todo stop server

	},
}
