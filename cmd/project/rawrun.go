package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"yi/internal/sdk"
	t "yi/pkg/types"
)

var RawRunCmd = &cobra.Command{
	Use:                "rawrun",
	Short:              "Run a command",
	Aliases:            []string{"rr"},
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.PanicLevel)

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if t.IsProjectDir(wd) {
			pC := t.NewPackageConfig()
			err = pC.LoadFromDir(wd)
			if err != nil {
				log.Fatal(err)
			}

			err = pC.GetCacheSDK().RunCommand(args, wd)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if sdk.GlobalSDKManger.Size() == 0 {
				log.Fatal("未找到默认编译器")
			}

			err := (*sdk.GlobalSDKManger.GetSDKs())[0].RunCommand(args, wd)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
