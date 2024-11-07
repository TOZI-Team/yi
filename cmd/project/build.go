package project

import (
	"Yi/internal/sdk"
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
		uSdk := (*sdk.GlobalSDKManger.GetSDKs())[0]
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
	//cjpmBuildOptions = cjpmPackage.NewCJPMBuildOptions()

	BuildCommand.Flags().BoolVarP(&buildOptions.IsRelease, "release", "r", false, "Build a release")
}
