package sdkCmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"yi/internal/sdk"
	"yi/internal/tui/box/compiler"
)

var setCmd = &cobra.Command{
	Use:  "set",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		defer sdk.WriteGlobal()

		if args[0] == "default" {
			if len(args) == 2 {
				sdk.GlobalSDKManger.Default = args[1]
			} else {
				s := compiler.ChooseCompiler(sdk.GlobalSDKManger.GetSDKs())
				sdk.GlobalSDKManger.Default = s.Path
			}
		} else {
			log.Fatal("Invalid argument")
		}
	},
}

func init() {
	Command.AddCommand(setCmd)
}
