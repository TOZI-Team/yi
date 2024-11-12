package project

import (
	"Yi/internal/sdk"
	cjpmPackage "Yi/pkg/backend/cjpm/package"
	t "Yi/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var buildOptions *t.BuildOptions

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

		log.Debug("Write backend config to disk.")
		p := t.NewPackageConfig() // 获取包的设置
		err = p.LoadFromDir(wd)
		if err != nil {
			log.Fatal(err.Error())
		}
		p.SetBackend(cjpmPackage.NewCJPMConfig())
		p.SyncToBackendConfig()
		err = p.WriteBackendConfigToDisk()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Debug("检查缓存")
		err = p.CheckCache()
		if err != nil {
			log.Warn(err.Error())
			sdk, err := sdk.GlobalSDKManger.FindByVersion(p.Base.ComVer)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = p.ResetCache(sdk.Path)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		uSdk := p.GetCacheSDK()
		log.Infof("Use SDK: %s", uSdk.Path)
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err.Error())
		}
		output, err := uSdk.BuildProject(workDir, *buildOptions)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("Successfully build project successfully, output: %s", output)
	},
}

func init() {
	buildOptions = t.NewBuildOptions()

	cjpmBuildOptions := cjpmPackage.NewCJPMBuildOptions()
	buildOptions.SetBackend(cjpmBuildOptions)

	BuildCommand.Flags().BoolVarP(&buildOptions.IsRelease, "release", "r", false, "Build a release")
	BuildCommand.Flags().BoolVar(&buildOptions.RunAfterBuild, "run", false, "Run after build")
}
