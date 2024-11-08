package types

import (
	yError "Yi/pkg/types/error"
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	SDK         *SDKInfo
}

//func (i InitConfig) ToCJPMConfig() cjpmPackage.CJPMConfig {
//	c := cjpmPackage.CJPMConfig{
//		Package: cjpmPackage.PackageConfig{
//			ComVer:      i.ComVer,
//			Description: i.Description,
//			Name:        i.Name,
//			Version:     i.Version,
//		},
//	}
//	return c
//}

type PackageConfig struct {
	Base struct {
		ComVer         string     `yaml:"cjc-version"`
		CompilerOption string     `yaml:"compiler-option"`
		Description    string     `yaml:"description"`
		Name           string     `yaml:"name"`
		Version        string     `yaml:"version"`
		ScrPath        string     `yaml:"src-dir"`
		TargetPath     string     `yaml:"target-dir"`
		OutputType     OutputType `yaml:"output-type"`
		Backend        string     `yaml:"backend"`
		path           string
		Authors        []string
	}
	cache struct {
		XMLName     xml.Name `xml:"Cache"`
		Version     string   `xml:"version,attr"`
		CompilerSet SDKInfo  `xml:"cacheCompiler"`
	}
	backend BackendProjectConfig
}

func (c PackageConfig) LoadFromDir(p string) error {
	// 加载 project.yml
	yamlFile, err := os.ReadFile(path.Join(p, "./project.yml"))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &c.Base)
	if err != nil {
		return err
	}

	// 加载 project.lock
	xmlFile, err := os.ReadFile(path.Join(p, "./package.lock"))
	if err != nil {
		return err
	}
	err = xml.Unmarshal(xmlFile, &c.cache)
	if err != nil {
		return err
	}
	return nil
}

func (c PackageConfig) GenerateFromInitConfig(config *InitConfig, backend BackendProjectConfig) {
	c.Base.path = config.Path
	//c.backend = backend
	c.Base.Name = config.Name
	c.Base.ComVer = config.SDK.Ver
	c.Base.Description = config.Description
	c.Base.Authors = []string{config.Authors}
	c.Base.Version = config.Version
	c.cache.CompilerSet = *config.SDK

	c.SetBackend(backend)
	log.Info(c.Base)
}

func (c PackageConfig) SetBackend(config BackendProjectConfig) {
	c.backend = config
}

func (c PackageConfig) WriteToDisk() error {
	log.Info(c.Base)
	data, err := yaml.Marshal(c.Base) // 写入 project.yml
	if err != nil {
		return err
	}
	log.Info("项目配置写入： ", c.Base.path)
	err = os.WriteFile(path.Join(c.Base.path, "./project.yml"), data, 0644)
	if err != nil {
		return err
	}

	data, err = xml.Marshal(c.cache) // 写入 project.lock
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(c.Base.path, "project.lock"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c PackageConfig) SyncToBackendConfig() {
	c.backend.GenerateFromProjectConfig(c)
}

// CheckCache 检查项目缓存是否有效
// 当缓存中的编译器无效时
func (c PackageConfig) CheckCache() error {
	if !c.cache.CompilerSet.CheckIsHave() {
		return yError.NewNotFoundSDKErr(c.cache.CompilerSet.Path)
	}
	return nil
}

func NewPackageConfig() *PackageConfig {
	c := new(PackageConfig)
	c.cache.Version = "1"
	return c
}

type OutputType = string

const (
	EXECUTABLE = "executable"
	STATIC     = "static"
	DYNAMIC    = "dynamic"
)

var DefaultInitConfig = InitConfig{
	Path:        "./",
	Name:        "exampleProject",
	Authors:     "You <you@example.com>",
	Version:     "1.0.0",
	Description: "lalalala",
}

type BuildOptions struct {
	ProjectPath string
	IsRelease   bool
	backendOpt  BackendBuildOptions
}

func (opts BuildOptions) SyncToBackendOpt() error {
	if opts.backendOpt == nil {
		return yError.NewNoBackendError("Not found build backend")
	}
	opts.backendOpt.RewriteFromBuildOptions(&opts)
	return nil
}

// MakeBackendShellArgs 获取编译后端对应command
func (opts BuildOptions) MakeBackendShellArgs() ([]string, error) {
	err := opts.SyncToBackendOpt()
	if err != nil {
		return nil, err
	}
	return opts.backendOpt.ToShellArgs(), nil
}

func NewBuildOptions() *BuildOptions {
	return &BuildOptions{}
}

func (opts BuildOptions) SetBackend(options BackendBuildOptions) {
	opts.backendOpt = options
}

// SetBackendOptions 设置编译后端配置
func (opts BuildOptions) SetBackendOptions(bOptions BackendBuildOptions) error {
	opts.backendOpt = bOptions
	err := opts.SyncToBackendOpt()
	if err != nil {
		return err
	}
	return nil
}

// GetBackendOptions 获取编译后端配置
func (opts BuildOptions) GetBackendOptions() *BackendBuildOptions {
	err := opts.SyncToBackendOpt()
	if err != nil {
		return nil
	}
	return &opts.backendOpt
}

// GetOutputPath 获取后端构建产物的相对位置
func (opts BuildOptions) GetOutputPath() string {
	return opts.backendOpt.GetOutputPath()
}
