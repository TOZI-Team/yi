package cjpmPackage

import (
	t "Yi/pkg/types"
	"bytes"
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type CJPMPackageConfig struct {
	ComVer         string       `toml:"cjc-version"`
	CompilarOption string       `toml:"compiler-option"`
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
	Package CJPMPackageConfig     `toml:"package"`
	Depends map[string]CJPMDepend `toml:"depends"`
}

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

// WriteToConfig 将配置写人文件
// 当传入目录时，自动添加文件名
func (c CJPMConfig) WriteToConfig(p string) error {
	buf := bytes.NewBuffer([]byte{})

	err := toml.NewEncoder(buf).Encode(&c) // 将对象编码为TOML
	if err != nil {
		return err
	}

	s, err := os.Stat(p) // 如果传参为目录，则自动添加文件名
	if err != nil {
		return err
	}
	if s.IsDir() {
		p = path.Join(p, "./cjpm.toml")
	}

	err = os.WriteFile(p, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
