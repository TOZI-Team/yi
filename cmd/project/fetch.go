package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"yi/internal/tui/project"
	cjpackage "yi/pkg/package"
)

var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		dir, err := cjpackage.LoadPackageFromDir(wd, true)
		if err != nil {
			log.Fatal(err)
		}

		dependencies, err := dir.AllOnlineDependencies()
		if err != nil {
			log.Fatal(err)
		}

		err = project.NewPackageDownloadTUI(dependencies).Run()
		if err != nil {
			log.Fatal(err)
		}

		err = project.NewDepMakeTUI(dependencies).Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}
