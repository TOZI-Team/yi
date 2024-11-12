package tSdk

type Channel string

const (
	BETA Channel = "beta"
	LTS  Channel = "lts"
	DEV  Channel = "dev"
	NULL Channel = "null"
)

type SDKFullInfo struct {
	Name        string
	Version     string
	Channel     Channel
	Path        string
	ManagerByOS bool
}
