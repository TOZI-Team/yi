package sdk

import (
	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	t "yi/pkg/types"
	yError "yi/pkg/types/error"
)

//var defaultMatch = regexp.MustCompile("\\[default]")

type Manager struct {
	Sdks    []t.SDKInfo `yaml:"Sdks"`
	Default string      `yaml:"default"`
}

func NewSDKManager() *Manager {
	sdks := new([]t.SDKInfo)
	return &Manager{Sdks: *sdks}
}

func (m *Manager) GetSDKs() *[]t.SDKInfo {
	return &m.Sdks
}

func (m *Manager) AddSDK(p string) error {
	sdk, err := t.NewSDKInfo(p)
	if err != nil {
		return err
	}
	m.Sdks = append(m.Sdks, *sdk)
	return nil
}

func (m *Manager) LoadFromFile(p string) error {
	s, err := os.Stat(p)
	if err != nil {
		return err
	}
	if s.IsDir() {
		p = path.Join(p, "compilers")
	}

	f, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) WriteToDisk(f string) error {
	s, err := os.Stat(f)
	if err == nil && s.IsDir() {
		f = path.Join(f, "compilers")
	}

	out, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = os.WriteFile(f, out, 0664)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Size() int {
	return len(m.Sdks)
}

func (m *Manager) FindByVersion(v string) (*t.SDKInfo, error) {
	for _, sdk := range m.Sdks {
		if sdk.Ver == v {
			return &sdk, nil
		}
	}
	return nil, yError.NewNotFoundSDKErr("Not find compiler what version is " + v)
}

func (m *Manager) FindByPath(p string) (int, error) {
	if !path.IsAbs(p) {
		return -1, yError.NewNotFoundSDKErr("Not absolute path")
	}
	for i, sdk := range m.Sdks {
		if sdk.Path == p {
			return i, nil
		}
	}
	return -1, yError.NewNotFoundSDKErr("Not find compiler what path is " + p)
}

func (m *Manager) GetDefault() *t.SDKInfo {
	//if s := viper.GetString("compiler-path"); s != "" {
	//	sdk, err := t.NewSDKInfo(s)
	//	if err == nil {
	//		return sdk
	//	}
	//}
	if s := os.Getenv("Yi_Compiler_PATH"); s != "" {
		sdk, err := t.NewSDKInfo(s)
		if err == nil {
			return sdk
		}
	}

	for _, sdk := range m.Sdks {
		if m.Default == sdk.Path {
			return &sdk
		}
	}
	return &m.Sdks[0]
}

func (m *Manager) SetDefault(p string) {
	m.Default = p
}

func (m *Manager) RemoveSDK(i int) {
	m.Sdks = append(m.Sdks[:i], m.Sdks[i+1:]...)
}

var GlobalSDKManger *Manager

func init() {
	GlobalSDKManger = NewSDKManager()
	cd := configdir.LocalConfig("yi")
	err := configdir.MakePath(cd)
	if err != nil {
		log.Fatalf("A error found in Init(): %s", err.Error())
	}
	if _, err := os.Open(path.Join(cd, "./compilers")); err == nil {
		err := GlobalSDKManger.LoadFromFile(path.Join(cd, "./compilers"))
		if err != nil {
			log.Fatalf("A error found in Init(): %s", err.Error())
		}
	} else {
		err := GlobalSDKManger.WriteToDisk(path.Join(cd, "./compilers"))
		if err != nil {
			log.Fatalf("A error found in Init(): %s", err.Error())
		}
	}
}

func WriteGlobal() {
	cd := configdir.LocalConfig("yi")

	log.Info("Write Global SDKs to config")
	err := GlobalSDKManger.WriteToDisk(path.Join(cd, "./compilers"))
	if err != nil {
		log.Fatalf("A error found in Init(): %s", err.Error())
	}
}
