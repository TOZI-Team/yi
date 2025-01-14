package sdkCmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path/filepath"
	"yi/internal/sdk"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a compiler",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatal(err)
		}

		err = sdk.GlobalSDKManger.AddSDK(p)
		if err != nil {
			log.Fatal(err.Error())
		}

		sdk.WriteGlobal()
	},
}

func init() {
	Command.AddCommand(addCmd)
}
