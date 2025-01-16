package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"yi/internal/sdk"
	cjpmPackage "yi/pkg/backend/cjpm/package"
	cjpackage "yi/pkg/package"
	t "yi/pkg/types"
)

var buildOptions *t.BuildOptions
var legacy bool

//var cjpmBuildOptions *cjpmPackage.CJPMBuildOptions

var BuildCommand = &cobra.Command{
	Use:   "build",
	Short: "Compile a local module and all of its dependencies.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err.Error())
		}

		if sdk.GlobalSDKManger.Size() == 0 { // 如果SDK列表为空，则提示首先添加 SDK
			log.Fatal("SDK manager is empty")
		}

		var opt = cjpackage.BuildOptions{IsRelease: buildOptions.IsRelease, RunAfterBuild: false, BuildType: 0, ShowOutput: false}

		if legacy || (!cjpackage.IsFuxoPackage(wd) && cjpackage.IsLegacyPackage(wd)) {
			result, err := cjpmPackage.CJPMProjectBackend{}.Build(&opt)
			if err != nil {
				log.Error(err.Error())
			}

			if !result.Success {
				log.Exit(1)
			}
		} else {
			p, err := cjpackage.LoadPackageFromDir(wd, true)
			if err != nil {
				log.Fatal(err.Error())
			}

			build, err := p.Build(&opt)
			if err != nil {
				return
			}

			if !build.Success {
				log.Exit(1)
			}
		}
	},
}

func init() {
	buildOptions = t.NewBuildOptions()

	cjpmBuildOptions := cjpmPackage.NewCJPMBuildOptions()
	buildOptions.SetBackend(cjpmBuildOptions)

	BuildCommand.Flags().BoolVarP(&buildOptions.IsRelease, "release", "r", false, "Build a release")
	BuildCommand.Flags().BoolVar(&buildOptions.RunAfterBuild, "run", false, "Run after build")
	BuildCommand.Flags().BoolVarP(&legacy, "legacy", "", false, "Use legacy build")
}
