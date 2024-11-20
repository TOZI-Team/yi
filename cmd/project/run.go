package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"yi/internal/sdk"
	t "yi/pkg/types"
)

var global bool

var RunCommand = &cobra.Command{
	Use:     "run",
	Short:   "Run a command",
	Aliases: []string{"r"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmds := strings.Split(args[0], " ")

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if (!global) && t.IsProjectDir(wd) {
			pC := t.NewPackageConfig()
			err = pC.LoadFromDir(wd)
			if err != nil {
				log.Fatal(err)
			}

			err = pC.GetCacheSDK().RunCommand(cmds, wd)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if sdk.GlobalSDKManger.Size() == 0 {
				log.Fatal("未找到默认编译器")
			}

			err := sdk.GlobalSDKManger.GetDefault().RunCommand(cmds, wd)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	RunCommand.Flags().BoolVarP(&global, "global", "g", false, "Use default SDK")
}
