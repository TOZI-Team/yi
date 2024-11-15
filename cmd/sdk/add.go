package sdkCmd

import (
	"Yi/internal/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a compiler",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := sdk.GlobalSDKManger.AddSDK(args[0])
		if err != nil {
			log.Fatal(err.Error())
		}

		sdk.WriteGlobal()
	},
}

func init() {
	Command.AddCommand(addCmd)
}
