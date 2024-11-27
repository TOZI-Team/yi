package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"yi/internal/sdk"
	cjpmPackage "yi/pkg/backend/cjpm/package"
	t "yi/pkg/types"
)

var global bool
var sdkDefault string

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
			pC := t.NewPackageConfigV1()
			pC.SetBackend(cjpmPackage.NewCJPMConfigV1())
			err = pC.LoadFromDir(wd)
			if err != nil {
				log.Fatal(err)
			}

			if s := pC.FindScript(args[0]); s != "" {
				log.Debug("Find the script: ", args[0])
				cmds = strings.Split(s, " ")
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
	RunCommand.Flags().StringVar(&sdkDefault, "compiler-path", "", "Set compiler path")

	//err := viper.BindPFlag("global", RunCommand.Flags().Lookup("global"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = viper.BindPFlag("compiler-path", RunCommand.Flags().Lookup("global"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = viper.BindEnv("compiler-path", "Yi_Compiler_PATH")
	//if err != nil {
	//	log.Fatal(err)
	//}
}
