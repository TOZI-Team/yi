package project

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
	"yi/internal/sdk"
	"yi/internal/tui/project"
	cjpmPackage "yi/pkg/backend/cjpm/package"
	t "yi/pkg/types"
)

var isOverWrite bool

var NewCommand = &cobra.Command{
	Use:   "new projectPath [-o]",
	Short: "Creat new Cangjie Project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if sdk.GlobalSDKManger.Size() == 0 { // 如果SDK列表为空，则提示首先添加 SDK
			log.Fatal("SDK manager is empty")
		}

		if isOverWrite {
			log.Warningf("将清空文件夹：%s", args[0])
		}

		iC := t.DefaultInitConfig
		//wd, err := os.Getwd()
		//if err != nil {
		//	log.Fatal(err)
		//}
		iC.Path = args[0]
		iC.Name = path.Base(args[0])
		iC = project.InitGuide(iC) // 启用TUI引导
		c := t.NewPackageConfigV1()
		c.GenerateFromInitConfig(&iC)
		c.SetBackend(cjpmPackage.NewCJPMConfigV1())
		//log.Info(c.Base)
		s, err := os.Stat(args[0]) // 判断是否存在同名文件
		if err == nil {
			if !s.IsDir() {
				log.Fatal("已存在对应文件")
			}

			ds, err := os.ReadDir(args[0]) // 判断同名文件夹是否为空
			if err != nil {
				log.Fatal(err.Error())
			}
			if len(ds) != 0 {
				if !isOverWrite {
					log.Fatal("指定文件夹非空")
				}
				for _, d := range ds { // 清空文件夹
					err = os.RemoveAll(path.Join(args[0], d.Name()))
					if err != nil {
						log.Fatal(err.Error())
					}
				}
			}
		}

		err = os.MkdirAll(path.Join(args[0], "./src"), os.ModePerm) //创建相关文件
		if err != nil {
			log.Fatal(err.Error())
		}
		if err := os.WriteFile(path.Join(args[0], "./src/demo.cj"), []byte(fmt.Sprintf("package %s\n\n// You can write Cangjie code here.\n", c.Name)), os.ModePerm); err != nil {
			log.Fatal(err.Error())
		}

		//log.Debug(c.Base)
		if err := c.WriteToDisk(); err != nil {
			log.Fatal(err.Error())
		}
		log.Info("创建项目成功")
	},
}

func init() {
	NewCommand.Flags().BoolVarP(&isOverWrite, "overwrite", "o", false, "Overwrite existing project")
}
