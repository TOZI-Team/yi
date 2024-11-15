package sdkCmd

import (
	"fmt"
	"github.com/kirsle/configdir"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/walle/targz"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"yi/internal/sdk"
	t "yi/pkg/types"
)

func copyFolder(source string, destination string) error {
	// 遍历源文件夹
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 计算目标路径
		relativePath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destination, relativePath)
		if info.IsDir() {
			// 如果是文件夹，则创建对应的目标文件夹
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// 如果是文件，则复制文件内容
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(srcFile *os.File) {
				err := srcFile.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(srcFile)
			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer func(destFile *os.File) {
				err := destFile.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(destFile)
			_, err = io.Copy(destFile, srcFile)
			return err
		}
	})
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sdk from local or Internet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer sdk.WriteGlobal()

		url := args[0]

		u, err := user.Current() // Get user
		if err != nil {
			log.Fatal(err)
		}

		if url[:4] == "http" {
			// TODO Support Internet
			log.Fatal("暂不支持从网络获取")
		}

		log.Info("Creat cache folder")
		cDir := configdir.LocalCache("Yi")
		err = configdir.MakePath(path.Join(cDir, "./unzip/", filepath.Base(url)))
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Extract zip file")
		err = targz.Extract(url, path.Join(cDir, "./unzip/", filepath.Base(url)))
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Get version")
		err = exec.Command("sh", "-c", "chmod 0755 "+path.Join(cDir, "./unzip/", filepath.Base(url), "cangjie")+" -R").Run()
		if err != nil {
			log.Fatal(err)
		}
		info, err := t.NewSDKInfo(path.Join(cDir, "./unzip/", filepath.Base(url), "cangjie"))
		if err != nil {
			log.Fatal(err)
		}

		log.Info("创建目录")
		p := path.Join(u.HomeDir, ".Yi/Cangjie", fmt.Sprintf("%s-%s", info.Ver, uuid.NewV4()))
		err = configdir.MakePath(p)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Copy files to ", p)
		err = copyFolder(path.Join(cDir, "./unzip/", filepath.Base(url), "cangjie"), p)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Get version")
		c := exec.Command("sh", "-c", "chmod 0755 "+p+" -R")
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			log.Fatal(err)
		}

		err = sdk.GlobalSDKManger.AddSDK(p)
		if err != nil {
			log.Fatal(err)
		}

		defer func(url string) {
			log.Info("Delete cache dir")
			err = os.RemoveAll(path.Join(cDir, "./unzip/", filepath.Base(url)))
			if err != nil {
				log.Error(err)
			}
		}(url)
	},
}

func init() {
	Command.AddCommand(installCmd)
}
