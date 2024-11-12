package cjpmPackage

import (
	t "Yi/pkg/types"
	"bytes"
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type PackageConfig struct {
	ComVer         string       `toml:"cjc-version"`
	CompilerOption string       `toml:"compiler-option"`
	Description    string       `toml:"description"`
	Name           string       `toml:"name"`
	Version        string       `toml:"version"`
	ScrPath        string       `toml:"src-dir"`
	TargetPath     string       `toml:"target-dir"`
	OutputType     t.OutputType `toml:"output-type"`
}

type CJPMDepend struct {
	GitURL    string `toml:"git"`
	GitBranch string `toml:"branch"`
	Path      string `toml:"path"`
	Version   string `toml:"version"`
}

type CJPMConfig struct {
	Package PackageConfig         `toml:"package"`
	Depends map[string]CJPMDepend `toml:"depends"`
}

// Deprecated: Use t.ProjectConfig
func ReadCJPMConfig(p string) (CJPMConfig, error) {
	s, err := os.Stat(p)
	if err != nil {
		return CJPMConfig{}, err
	}
	if s.IsDir() {
		p = path.Join(p, "./cjpm.toml")
	}

	var config CJPMConfig
	_, err = toml.DecodeFile(p, &config)
	if err != nil {
		return CJPMConfig{}, err
	}
	return config, nil
}

func (c *CJPMConfig) WriteConfigToDir(p string) error {
	buf := bytes.NewBuffer([]byte{})

	err := toml.NewEncoder(buf).Encode(&c) // 将对象编码为TOML
	if err != nil {
		return err
	}

	s, err := os.Stat(p) // 如果传参为目录，则自动添加文件名
	if err == nil && s.IsDir() {
		p = path.Join(p, "./cjpm.toml")
	}

	err = os.WriteFile(p, buf.Bytes(), 0666) //写入文件
	if err != nil {
		return err
	}
	return nil
}

// WriteToConfig 将配置写人文件
// 当传入目录时，自动添加文件名
func (c *CJPMConfig) WriteToConfig(p string) error {
	return c.WriteConfigToDir(p)
}

func (c *CJPMConfig) GenerateFromPackageConfig(config t.PackageConfig) {
	c.Package.Name = config.Base.Name
	c.Package.Version = config.Base.Version
	c.Package.ComVer = config.Base.ComVer
	c.Package.Description = config.Base.Description
	c.Package.ScrPath = config.Base.ScrPath
	c.Package.TargetPath = config.Base.TargetPath
	c.Package.OutputType = config.Base.OutputType
	c.Package.CompilerOption = config.Base.CompilerOption
}
func (c *CJPMConfig) GenerateFromProjectConfig(config t.PackageConfig) {
	c.GenerateFromPackageConfig(config)
}

func NewCJPMConfig() *CJPMConfig {
	return &CJPMConfig{Package: PackageConfig{OutputType: t.EXECUTABLE}}
}
