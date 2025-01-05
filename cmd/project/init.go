package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
	"yi/internal/sdk"
	cjpackage "yi/pkg/package"
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

		_, err = os.Stat(path.Join(wd, "./fuxo.toml"))
		if os.IsNotExist(err) {

		} else if err != nil {
			log.Fatal(err)
		} else {
			if !isOverWrite {
				log.Fatal("配置文件已存在")
			}
		}

		c := cjpackage.InitConfigToProjectConfig(path.Base(wd), "1.0.0", sdk.GlobalSDKManger.Sdks[0].Ver)
		//c := t.NewPackageConfigV1()
		//c.SetBackend(cjpmPackage.NewCJPMConfigV1())
		//c.GenerateFromInitConfig(&iC)

		if err := c.WriteToDisk(wd, false); err != nil {
			log.Fatal(err.Error())
		}
		log.Info("项目初始化成功")
	},
}

func init() {
	InitCmd.Flags().BoolVarP(&isOverWrite, "overwrite", "o", false, "Overwrite existing project")
}
