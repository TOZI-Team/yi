package fuxo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	devlog "yi/log"
	"yi/pkg/repo/dl"
	"yi/pkg/repo/index"
)

type RepoConfig struct {
	Download string `toml:"download"`
	Index    string `toml:"index"`
	API      string `toml:"api"`
}

func (c RepoConfig) Name() string {
	return fmt.Sprintf("%s-%s", c.Index[0:6], hash(c.Download))
}

func hash(s string) string {
	sha := sha256.Sum224([]byte(s))
	return hex.EncodeToString(sha[:])
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
			if _, err := os.Stat(p); err != nil {
				if os.IsNotExist(err) {
					return nil
				} else {
					return err
				}
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

//func (c RepoConfig) GetCache() *cache.Config {
//	return cache.NewConfig(c.Name())
//}

func (c RepoConfig) GetIndex() *index.Index {
	return index.NewIndex(&c.Index, index.GitIndex, c.Name())
}

func (c *Config) GetDefaultDl() *dl.Dl {
	config := c.Repos["fuxo"]
	return config.GetDL()
}

func (c RepoConfig) HavePackageCache(name, ver string) bool {
	p, err := c.GetFuPath(name, ver)
	if err != nil {
		return false
	}

	if _, err := os.Stat(p); err != nil {
		return false
	}

	return true
}

func (c RepoConfig) GetFuPath(name, ver string) (string, error) {
	wd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := path.Join(wd, ".fuxo", "fu", c.Name(), fmt.Sprintf("%s-%s.fu", name, ver))
	return p, nil
}

//func (c RepoConfig) GetFuCache(name, ver string) (*fu.Fu, error) {
//	if !c.HavePackageCache(name, ver) {
//		return nil, fmt.Errorf("package cache not exist")
//	}
//
//	p, err := c.GetFuPath(name, ver)
//	if err != nil {
//		return nil, err
//	}
//	return fu.LoadFromDisk(p), nil
//}

func (c *Config) GetDefaultIndex() *index.Index {
	return c.Repos["fuxo"].GetIndex()
}

func (c *Config) GetDefault() RepoConfig {
	return c.Repos["fuxo"]
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
	Name, Ver, Repo string
}

func (s SimplePackageMeta) Download() error {
	if s.Repo == "" {
		err := DownloadPackageToCache(s)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not suport custom repo now")
	}

	return nil
}

// FindAllDep
//
// Note: 可能无法有效处理循环依赖问题
func FindAllDep(name string, version string) ([]SimplePackageMeta, error) {
	deps := make([]SimplePackageMeta, 0)

	deps = append(deps, SimplePackageMeta{name, version, ""})
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

// DownloadPackageToCache 下载包至缓存文件夹
func DownloadPackageToCache(meta SimplePackageMeta) error {
	// 若已缓存，则直接返回
	if globalConfig.GetDefault().HavePackageCache(meta.Name, meta.Ver) {
		return nil
	}

	d := GlobalConfig().GetDefaultDl().GetDlUrl(meta.Name, meta.Ver)
	r, err := http.Get(d)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			devlog.DevLog.Errorf("close body err: %v", err)
		}
	}(r.Body)

	p, err := GlobalConfig().GetDefault().GetFuPath(meta.Name, meta.Ver)
	if err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			devlog.DevLog.Error(err)
		}
	}(f)

	_, err = io.Copy(f, r.Body)

	return nil
}
