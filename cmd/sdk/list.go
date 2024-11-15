package sdkCmd

import (
	"Yi/internal/sdk"
	"encoding/json"
	"github.com/scylladb/termtables"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var format string
var online bool

var listCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List Sdks",
	Run: func(cmd *cobra.Command, args []string) {
		if online {
			log.Fatal("Sorry,not support the arg: --online")
		}

		sdks := sdk.GlobalSDKManger.GetSDKs()
		switch format {
		case "json":
			data, err := json.Marshal(sdks)
			if err != nil {
				log.Fatal(err)
			}
			print(string(data))
		case "normal":
			t := termtables.CreateTable()
			t.AddHeaders("Version", "Path", "Note")
			for _, sdk := range *sdks {
				t.AddRow(sdk.Ver, sdk.Path, sdk.Note)
			}
			print(t.Render())
		}
	},
}

func init() {
	listCommand.Flags().StringVarP(&format, "format", "f", "normal", "Output format: normal|json")
	listCommand.Flags().BoolVar(&online, "online", false, "List online sdks")
	Command.AddCommand(listCommand)
}
