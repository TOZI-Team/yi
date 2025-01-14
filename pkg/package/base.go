package cj_package

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path"
	devlog "yi/log"
	fuxo "yi/pkg/repo"
	"yi/pkg/repo/index"
	"yi/pkg/types"
)

type PackageConfig struct {
	Project struct {
		Name         string   `toml:"name"`
		Version      string   `toml:"version"`
		Authors      []string `toml:"authors"`
		CjcVersion   string   `toml:"cjcv,omitempty"`
		BuildScript  string   `toml:"build,omitempty"`
		LinkBins     []string `toml:"link,omitempty"`
		ExcludeFiles []string `toml:"exclude,omitempty"`
		IncludeFiles []string `toml:"include,omitempty"`
		AblePublish  bool     `toml:"publish,omitempty"`
		Description  string   `toml:"description,omitempty"`
		Readme       string   `toml:"readme,omitempty"`
		Categories   []string `toml:"categories,omitempty"`
	} `toml:"project"`
	Dependencies struct {
		SimplePackage map[string]string `toml:"-"`
		OnlinePackage map[string]struct {
			Version        string            `toml:"version"`
			PackageFeature map[string]string `toml:"feature,omitempty"`
			IsLegacy       bool              `toml:"legacy,omitempty"`
		} `toml:"-"`
		LocalPackage map[string]struct {
			Path     string `toml:"path"`
			IsLegacy bool   `toml:"legacy,omitempty"`
		} `toml:"-"`
		GitPackage map[string]struct {
			URL      string `toml:"git"`
			Branch   string `toml:"branch"`
			Commit   string `toml:"commit"`
			Tag      string `toml:"tag"`
			IsLegacy bool   `toml:"legacy,omitempty"`
		} `toml:"-"`
	} `toml:"dependencies"`
}

type GitCacheLessName struct {
	URL      string `toml:"git"`
	Branch   string `toml:"branch"`
	CommitID string `toml:"commit"`
}

type GitCache struct {
	Name string
	GitCacheLessName
}

type OnlineCacheLessName struct {
	Version string `toml:"version"`
}

type OnlineCache struct {
	Name string
	OnlineCacheLessName
}

// PackageCache
//
// @Description: 描述包缓存
type PackageCache struct {
	OnlinePackage map[string]OnlineCacheLessName `toml:"-"`
	GitPackage    map[string]GitCacheLessName    `toml:"-"`
}

// loadCacheFromFile
//
//	@Description: 从文件加载缓存
func loadCacheFromFile(path string) (*PackageCache, error) {
	c := new(PackageCache)
	_, err := toml.DecodeFile(path, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *PackageCache) AddCacheGit(g GitCache) {
	if _, ok := c.OnlinePackage[g.Name]; !ok {
		delete(c.OnlinePackage, g.Name)
	}
	c.GitPackage[g.Name] = g.GitCacheLessName
}

func (c *PackageCache) AddOnlineCache(o OnlineCache) {
	if _, ok := c.GitPackage[o.Name]; !ok {
		delete(c.GitPackage, o.Name)
	}
	c.OnlinePackage[o.Name] = o.OnlineCacheLessName
}

func (c *PackageCache) RemoveCache(name string) error {
	if _, ok := c.GitPackage[name]; ok {
		delete(c.GitPackage, name)
	} else if _, ok := c.OnlinePackage[name]; ok {
		delete(c.OnlinePackage, name)
	} else {
		return fmt.Errorf("no cache for %s", name)
	}
	return nil
}

// Package 封装符合 fuxo 规范的包
type Package struct {
	config *PackageConfig
	cache  *PackageCache
}

// GetDepCachePath 获取指定依赖的存放路径
func (p *Package) GetDepCachePath(name string) (string, error) {
	for k, v := range p.config.Dependencies.OnlinePackage {
		if v.IsLegacy {
			continue
		}

		if k == name {
			v, ok := p.Cache().OnlinePackage[k]
			if !ok {
				return "", fmt.Errorf("must update project first")
			}

			if !fuxo.GlobalConfig().GetDefault().HavePackageCache(name, v.Version) {
				return "", fmt.Errorf("dependency %s not found in cache", k)
			}
		}
	}

	for k, v := range p.Config().Dependencies.LocalPackage {
		if k == name {
			return v.Path, nil
		}
	}

	//TODO 支持git依赖

	return "", fmt.Errorf("not found %s", name)
}

// MakeCache 生成 Cache
func (p *Package) MakeCache(overwrite bool) error {
	if overwrite {
		p.cache = new(PackageCache)
	}

	// 处理在线包
	for k, v := range p.config.Dependencies.OnlinePackage {
		if v.IsLegacy {
			continue
		}

		if _, ok := p.cache.OnlinePackage[k]; !ok {
			findPS, err := fuxo.GlobalConfig().GetDefaultIndex().FindPackage(k)
			if err != nil {
				devlog.DevLog.Warnf("find package %s failed", k)
				return err
			}

			version, err := index.FindVersion(findPS, v.Version, false)
			if err != nil {
				devlog.DevLog.Warnf("find package %s failed", k)
				return err
			}

			p.cache.AddOnlineCache(OnlineCache{Name: k, OnlineCacheLessName: OnlineCacheLessName{Version: version.Version}})
		}

	}

	// 处理简写包
	for k, v := range p.config.Dependencies.SimplePackage {
		if _, ok := p.cache.OnlinePackage[k]; !ok {
			findPS, err := fuxo.GlobalConfig().GetDefaultIndex().FindPackage(k)
			if err != nil {
				return err
			}

			version, err := index.FindVersion(findPS, v, false)
			if err != nil {
				return err
			}

			p.cache.AddOnlineCache(OnlineCache{Name: k, OnlineCacheLessName: OnlineCacheLessName{Version: version.Version}})
		}
	}

	// 处理 git 包
	for k, _ := range p.config.Dependencies.GitPackage {
		//TODO 支持git包
		return fmt.Errorf("not support git package:%s", k)
	}
	return nil
}

func (p *Package) AllOnlineDependencies() ([]fuxo.SimplePackageMeta, error) {
	m := make([]fuxo.SimplePackageMeta, 0)
	// 处理在线包
	for k, v := range p.config.Dependencies.OnlinePackage {
		if v.IsLegacy {
			continue
		}

		if c, ok := p.cache.OnlinePackage[k]; ok {
			m = append(m, fuxo.SimplePackageMeta{Name: k, Ver: c.Version})
		} else {
			return nil, fmt.Errorf("must make cache first")
		}
	}

	// 处理简写包
	for k, _ := range p.config.Dependencies.SimplePackage {
		if c, ok := p.cache.OnlinePackage[k]; ok {
			m = append(m, fuxo.SimplePackageMeta{Name: k, Ver: c.Version})
		} else {
			return nil, fmt.Errorf("must make cache first")
		}
	}

	return m, nil
}

func (p *Package) Config() *PackageConfig {
	return p.config
}

func (p *Package) Cache() *PackageCache {
	return p.cache
}

func (p *Package) GetName() string {
	return p.config.Project.Name
}

func (p *Package) WriteToDisk(pt string, cache bool) error {
	d, err := toml.Marshal(*p.config)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(pt, "fuxo.toml"), d, 0644)
	if err != nil {
		return err
	}

	if cache {
		d, err := toml.Marshal(*p.cache)
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(pt, "fuxo.lock"), d, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadPackageFromDir(dp string, loadCache bool) (*Package, error) {
	p := new(Package)

	_, err := toml.DecodeFile(path.Join(dp, "fuxo.toml"), &p.config)
	if err != nil {
		return nil, err
	}
	if loadCache {
		if _, err := os.Stat(path.Join(dp, "fuxo.lock")); err == nil {
			p.cache, err = loadCacheFromFile(path.Join(dp, "fuxo.lock"))
			if err != nil {
				return nil, err
			}
		}
	}

	return p, nil
}

type BackendConfigOption struct {
	OutputType    string
	ProjectDir    string
	StaticDepends map[string]string
	UseSDK        *types.SDKInfo
}

func InitConfigToProjectConfig(name string, version string, cjcv string) *Package {
	p := new(Package)

	p.config = new(PackageConfig)
	p.config.Project.Name = name
	p.config.Project.Version = version
	p.config.Project.CjcVersion = cjcv
	p.config.Project.Authors = []string{"example <example@example.com>"}

	return p
}
