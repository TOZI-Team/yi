package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"yi/internal/sdk"
	"yi/internal/tui/project"
	cjpmPackage "yi/pkg/backend/cjpm/package"
	cjpackage "yi/pkg/package"
	fuxo "yi/pkg/repo"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update project deps",
	Run: func(cmd *cobra.Command, args []string) {

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if legacy || (!cjpackage.IsFuxoPackage(wd) && cjpackage.IsLegacyPackage(wd)) {
			p := cjpmPackage.CJPMConfigV1{}
			err := p.LoadFromDir(wd)
			if err != nil {
				log.Fatal(err)
			}

			version, err := sdk.GlobalSDKManger.FindByVersion(p.Package.ComVer)
			if err != nil {
				return
			}

			err = version.RunCommand([]string{"cjpm", "update"}, wd)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = project.NewRepoUpdateJobFromRepo(fuxo.GlobalConfig().GetDefault()).Run()
		if err != nil {
			log.Fatal(err)
		}

		p, err := cjpackage.LoadPackageFromDir(wd, false)
		if err != nil {
			log.Fatal(err)
		}

		err = p.MakeCache(true)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	UpdateCmd.Flags().BoolVar(&legacy, "legacy", false, "Use legacy dependencies")
}
