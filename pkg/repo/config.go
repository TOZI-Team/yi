package fuxo

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"yi/pkg/repo/dl"
	"yi/pkg/repo/index"
)

type RepoConfig struct {
	Download string `toml:"download"`
	Index    string `toml:"index"`
	API      string `toml:"api"`
}

// GenerateRepoConfig
//
//	@Description: 从仓库信息文件生成 RepoConfig 。
//	@param url 指向仓库信息文件的链接
func GenerateRepoConfig(url string) (*RepoConfig, error) {
	f, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	a, err := io.ReadAll(f.Body)
	if err != nil {
		return nil, err
	}
	return GenerateRepoConfigFromRepoInfo(&a)
}

// GenerateRepoConfigFromRepoInfo
//
//	@Description: 从 Fuxo 仓库信息文件内容创建本地配置
//	@param data Fuxo 仓库信息文件内容
func GenerateRepoConfigFromRepoInfo(data *[]byte) (*RepoConfig, error) {
	c := struct {
		Base struct {
			API   string `json:"api"`
			Index string `json:"index"`
		} `json:"base"`
	}{}
	err := json.Unmarshal(*data, &c)
	if err != nil {
		return nil, err
	}
	r := new(RepoConfig)
	r.Index = c.Base.Index
	r.API = c.Base.API
	r.Download = c.Base.API

	return new(RepoConfig), nil
}

// Config 仓库配置
type Config struct {
	Repos map[string]RepoConfig `toml:"-"`
}

// Save
//
//	@Description: 保存配置至磁盘
//	@param p 保存文件名，为空时使用用户配置文件夹。
func (c *Config) Save(p string) error {
	if p == "" {
		p = configdir.LocalConfig("fuxo")
		err := configdir.MakePath(p)
		if err != nil {
			return err
		}
		p = path.Join(p, "config.toml")
	}

	marshal, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(p, marshal, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Load
//
//	@Description:  从磁盘加载配置
//	@param p 配置路径，默认采用 用户配置 > 系统配置。
func (c *Config) Load(p string) error {
	if p == "" {
		p = path.Join(configdir.LocalConfig("fuxo"), "config.toml")
		if _, err := os.Stat(p); os.IsNotExist(err) {
			p = path.Join(configdir.SystemConfig("fuxo")[0], "config.toml")
			if _, err := os.Stat(p); os.IsNotExist(err) {
				return err
			}
		}
	}

	f, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(f, c)
	if err != nil {
		return err
	}
	return nil
}

func (c RepoConfig) GetDL() *dl.Dl {
	return dl.NewDl(&c.Download)
}

func (c RepoConfig) GetIndex(name string) *index.Index {
	return index.NewIndex(&c.Index, index.GitIndex, name)
}

func (c *Config) GetDefaultDl() *dl.Dl {
	return c.Repos["fuxo"].GetDL()
}

func (c *Config) GetDefaultIndex() *index.Index {
	return c.Repos["fuxo"].GetIndex("fuxo")
}

var globalConfig *Config

func init() {
	globalConfig = new(Config)
	err := globalConfig.Load("")
	if err != nil {
		log.Fatal(err)
	}
}

func GlobalConfig() *Config {
	return globalConfig
}

type SimplePackageMeta struct {
	Name, Ver string
}

// FindAllDep
//
// Note: 可能无法有效处理循环依赖问题
func FindAllDep(name string, version string) ([]SimplePackageMeta, error) {
	deps := make([]SimplePackageMeta, 0)

	deps = append(deps, SimplePackageMeta{name, version})
	for i := 0; len(deps) <= 1; i++ {
		pkg := deps[i]

		is, err := GlobalConfig().GetDefaultIndex().FindPackage(pkg.Name)
		if err != nil {
			return nil, err
		}

		findVer, err := index.FindVersion(is, pkg.Ver, false)
		if err != nil {
			return nil, err
		}

		for _, j := range findVer.Depends {
			deps = append(deps, SimplePackageMeta{Name: j.Name, Ver: j.ReqVer})
		}
	}

	return deps, nil
}
