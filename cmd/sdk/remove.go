package sdkCmd

import (
	"Yi/internal/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
)

var isDeletable bool

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Removes a SDK",
	Aliases: []string{"r"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer sdk.WriteGlobal()

		p, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatal(err.Error())
		}

		if isDeletable {
			if runtime.GOOS == "linux" {
				if filepath.Dir(p) == "/opt" {
					log.Fatal("该目录由系统管理，「绎」暂不支持直接删除。")
				}
			}
			log.Warn("Your Cangjie SDK will be deleted.")
		}

		s, err := sdk.GlobalSDKManger.FindByPath(p)
		if err != nil {
			log.Fatal(err.Error())
		}
		if isDeletable {
			err := os.RemoveAll(p)
			if err != nil {
				return
			}
		}
		sdk.GlobalSDKManger.RemoveSDK(s)
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&isDeletable, "delete", "d", false, "Delete the sdk from disk")

	Command.AddCommand(removeCmd)
}
