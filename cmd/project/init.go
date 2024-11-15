package project

import (
	"Yi/internal/sdk"
	cjpmPackage "Yi/pkg/backend/cjpm/package"
	t "Yi/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if isOverWrite {
			log.Warn("将覆盖当前配置")
		}

		_, err = os.Stat(path.Join(wd, "./project.yml"))
		if os.IsNotExist(err) {

		} else if err != nil {
			log.Fatal(err)
		} else {
			if !isOverWrite {
				log.Fatal("配置文件已存在")
			}
		}
		iC := t.DefaultInitConfig
		iC.Path = wd
		iC.Name = path.Base(wd)
		iC.ComVer = sdk.GlobalSDKManger.Sdks[0].Ver
		iC.SDK = &sdk.GlobalSDKManger.Sdks[0]

		c := t.NewPackageConfig()
		c.GenerateFromInitConfig(&iC, cjpmPackage.NewCJPMConfig())

		if err := c.WriteToDisk(); err != nil {
			log.Fatal(err.Error())
		}
		log.Info("项目初始化成功")
	},
}

func init() {
	InitCmd.Flags().BoolVarP(&isOverWrite, "overwrite", "o", false, "Overwrite existing project")
}
