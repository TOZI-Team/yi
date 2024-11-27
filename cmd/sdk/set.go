package sdkCmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"syscall"
	"yi/internal/sdk"
	"yi/internal/tui/box/compiler"
)

var tmp bool

var setCmd = &cobra.Command{
	Use:  "set",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		if args[0] == "default" {
			if len(args) == 2 {
				sdk.GlobalSDKManger.Default = args[1]
			} else {
				s := compiler.ChooseCompiler(sdk.GlobalSDKManger.GetSDKs())
				sdk.GlobalSDKManger.Default = s.Path
			}
			if tmp {
				err := syscall.Setenv("Yi_Compiler_PATH", sdk.GlobalSDKManger.Default)
				if err != nil {
					log.Fatal(err)
				}

				log.Debug("Set env Yi_Compiler_PATH=", sdk.GlobalSDKManger.Default)
				return
			}
		} else {
			log.Fatal("Invalid argument")
		}

		sdk.WriteGlobal()
	},
}

func init() {
	Command.AddCommand(setCmd)
	setCmd.Flags().BoolVarP(&tmp, "temp", "t", false, "Temp")
}
