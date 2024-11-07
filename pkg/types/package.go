package types

import (
	cjpmPackage "Yi/pkg/backend/cjpm/package"
	backendInterface "Yi/pkg/backend/interface"
	"bytes"
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type InitConfig struct {
	Path        string
	Name        string
	Authors     string
	Version     string
	Description string
	ComVer      string
}

func (i InitConfig) ToCJPMConfig() CJPMConfig {
	c := CJPMConfig{
		Package: CJPMPackageConfig{
			ComVer:      i.ComVer,
			Description: i.Description,
			Name:        i.Name,
			Version:     i.Version,
		},
	}
	return c
}

type OutputType = string

const (
	EXECUTABLE = "executable"
	STATIC     = "static"
	DYNAMIC    = "dynamic"
)

var DefaultConfig = InitConfig{
	Path:        "./",
	Name:        "exampleProject",
	Authors:     "You <you@example.com>",
	Version:     "1.0.0",
	Description: "lalalala",
}

type CJPMPackageConfig struct {
	ComVer         string     `toml:"cjc-version"`
	CompilarOption string     `toml:"compiler-option"`
	Description    string     `toml:"description"`
	Name           string     `toml:"name"`
	Version        string     `toml:"version"`
	ScrPath        string     `toml:"src-dir"`
	TargetPath     string     `toml:"target-dir"`
	OutputType     OutputType `toml:"output-type"`
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

type BuildOptions struct {
	ProjectPath string
	IsRelease   bool
	backendOpt  backendInterface.BackendBuildOptions
}

func (opts BuildOptions) SyncToBackendOpt() {
	opts.backendOpt.RewriteFromBuildOptions(&opts)
}

// MakeBackendShellArgs 获取编译后端对应command
func (opts BuildOptions) MakeBackendShellArgs() []string {
	opts.SyncToBackendOpt()
	return opts.backendOpt.ToShellArgs()
}

func NewBuildOptions() *BuildOptions {
	return &BuildOptions{backendOpt: cjpmPackage.NewCJPMBuildOptions()}
}

// SetBackendOptions 设置编译后端配置
func (opts BuildOptions) SetBackendOptions(bOptions backendInterface.BackendBuildOptions) {
	opts.backendOpt = bOptions
	opts.SyncToBackendOpt()
}

// GetBackendOptions 获取编译后端配置
func (opts BuildOptions) GetBackendOptions() *backendInterface.BackendBuildOptions {
	opts.SyncToBackendOpt()
	return &opts.backendOpt
}

// GetOutputPath 获取后端构建产物的相对位置
func (opts BuildOptions) GetOutputPath() string {
	return opts.backendOpt.GetOutputPath()
}
