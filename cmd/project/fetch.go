package project

import "github.com/spf13/cobra"

var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch dependencies",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
