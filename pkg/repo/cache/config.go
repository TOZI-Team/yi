package cache

import (
	"fmt"
	"os"
	"path"
	"yi/pkg/fu"
	cjpackage "yi/pkg/package"
	t "yi/pkg/types"
)

type Config struct {
	name string
}

func (c *Config) Name() string {
	return c.name
}

const READYSTATE = 0
const UNPACKSTATE = 1
const NOTHINESATE = 2

func (c *Config) CacheState(name string, version string) int {
	wd, err := os.Getwd()
	if err != nil {
		return NOTHINESATE
	}

	d := path.Join(wd, ".fuxo", "cache", c.Name(), fmt.Sprintf("%s-%s", name, version))
	if s, err := os.Stat(d); err == nil && s.IsDir() {
		return READYSTATE
	}

	d = path.Join(wd, ".fuxo", "fu", c.Name(), fmt.Sprintf("%s-%s.fu", name, version))
	if s, err := os.Stat(d); err == nil && !s.IsDir() {
		return UNPACKSTATE
	}

	return NOTHINESATE
}

func (c *Config) MakeCache(name string, version string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %v", err)
	}

	f := fu.LoadFromDisk(path.Join(wd, ".fuxo", "cache", c.Name(), fmt.Sprintf("%s-%s", name, version)))
	if _, err := os.Stat(c.GetCache(name, version)); !os.IsNotExist(err) {
		err := os.RemoveAll(c.GetCache(name, version))
		if err != nil {
			return err
		}
	}
	err = f.Unpack(c.GetCache(name, version))
	return err
}

func (c *Config) GenerateCache(name string, version string) error {
	p := c.GetCache(name, version)
	if p == "" {
		return fmt.Errorf("no cache found for %s", name)
	}

	dir, err := cjpackage.LoadPackageFromDir(p, false)
	if err != nil {
		return err
	}

	err = dir.MakeBackendConfig(&cjpackage.BackendConfigOption{OutputType: t.STATIC})
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetCache(name string, version string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return path.Join(wd, ".fuxo", "cache", c.Name(), fmt.Sprintf("%s-%s", name, version))
}

func NewConfig(name string) *Config {
	return &Config{name: name}
}
