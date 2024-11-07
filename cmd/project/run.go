package project

import "github.com/spf13/cobra"

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Compile and run an executable product",
	Args:  cobra.ExactArgs(1),
}
