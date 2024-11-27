package types

import (
	"bytes"
	"encoding/xml"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	yError "yi/pkg/types/error"
)

type InitConfig struct {
	Path        string
	Name        string
	Authors     string
	Version     string
	Description string
	ComVer      string
	SDK         *SDKInfo
	Output      OutputType
}

const (
	CacheFileName  = "project.lock"
	ConfigFileName = "project.yml"
)

//func (i InitConfig) ToCJPMConfig() cjpmPackage.CJPMConfig {
//	c := cjpmPackage.CJPMConfig{
//		Package: cjpmPackage.PackageConfigV0{
//			ComVer:      i.ComVer,
//			Description: i.Description,
//			Name:        i.Name,
//			Version:     i.Version,
//		},
//	}
//	return c
//}

// PackageConfigV0 描述项目配置
//
// Deprecated: 终止支持
type PackageConfigV0 struct {
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
	} `toml:"base"`
	cache struct {
		XMLName     xml.Name `xml:"Cache"`
		Version     string   `xml:"version,attr"`
		CompilerSet SDKInfo  `xml:"cacheCompiler"`
	}
	Scripts map[string]string `toml:"scripts"`
	backend BackendProjectConfigV0
}

func (c *PackageConfigV0) LoadFromDir(p string) error {
	// 设置 path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	c.Base.path = absPath

	// 加载 project.yml
	yamlFile, err := os.ReadFile(path.Join(p, ConfigFileName))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &c.Base)
	if err != nil {
		return err
	}

	// 加载 project.lock
	xmlFile, err := os.ReadFile(path.Join(p, CacheFileName))
	if err != nil {
		return err
	}
	err = xml.Unmarshal(xmlFile, &c.cache)
	if err != nil {
		return err
	}
	return nil
}

func (c *PackageConfigV0) GenerateFromInitConfig(config *InitConfig, backend BackendProjectConfigV0) {
	c.Base.path = config.Path
	//c.backend = backend
	c.Base.Name = config.Name
	c.Base.ComVer = config.SDK.Ver
	c.Base.Description = config.Description
	c.Base.Authors = []string{config.Authors}
	c.Base.Version = config.Version
	c.cache.CompilerSet = *config.SDK
	c.Base.OutputType = config.Output
	c.SetBackend(backend)
	//log.Info(c.Base)
}

// SetBackend 方法用于设置 PackageConfigV0 结构体中的 backend 字段
func (c *PackageConfigV0) SetBackend(config BackendProjectConfigV0) {
	// 将传入的 config 参数赋值给 c.backend 字段
	c.backend = config
}

// WriteToDisk 方法用于将 PackageConfigV0 结构体中的 Base 字段和 cache 字段的内容分别写入到 project.yml 和 project.lock 文件中
func (c *PackageConfigV0) WriteToDisk() error {
	// 将 c.Base 序列化为 YAML 格式
	data, err := yaml.Marshal(c.Base)
	// 如果序列化过程中发生错误，记录错误并返回
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	// 记录日志，显示项目配置将被写入的路径
	log.Info("项目配置写入： ", c.Base.path)
	// 将序列化后的数据写入到 project.yml 文件中
	err = os.WriteFile(path.Join(c.Base.path, ConfigFileName), data, 0644)
	// 如果写入过程中发生错误，记录错误并返回
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	// 将后端配置同步到磁盘
	c.SyncToBackendConfig()
	// 将后端配置写入到磁盘
	err = c.WriteBackendConfigToDisk()
	// 如果写入过程中发生错误，记录错误并返回
	if err != nil {
		return err
	}

	// 返回 nil 表示写入成功
	return nil
}

func (c *PackageConfigV0) SyncToBackendConfig() {
	c.backend.GenerateFromProjectConfig(*c)
}

// CheckCache 检查项目缓存是否有效
// 当缓存中的编译器无效时
func (c *PackageConfigV0) CheckCache() error {
	if !c.cache.CompilerSet.CheckIsHave() {
		return yError.NewNotFoundSDKErr(c.cache.CompilerSet.Path)
	}
	return nil
}

func (c *PackageConfigV0) WriteBackendConfigToDisk() error {
	err := c.backend.WriteConfigToDir(c.Base.path)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (c *PackageConfigV0) ResetCache(p string) error {
	sdk, err := NewSDKInfo(p)
	if err != nil {
		return err
	}
	if sdk.Ver != c.Base.ComVer {
		log.Warn("编译器版本与项目设置不匹配。")
	}
	c.cache.CompilerSet = *sdk
	return nil
}

// IsProjectDir
//
//	@Description: 判断该目录是否为绎项目
//	@param p
//	@return bool
func IsProjectDir(p string) bool {
	s, err := os.Stat(path.Join(p, "cjpm.toml"))
	if err != nil {
		return false
	}
	if s.IsDir() {
		return false
	}

	return true
}

func (c *PackageConfigV0) GetCacheSDK() *SDKInfo {
	return &c.cache.CompilerSet
}

// NewPackageConfigV0 创建新的 PackageConfigV0
// Deprecated:
// 终止支持
func NewPackageConfigV0() *PackageConfigV0 {
	c := new(PackageConfigV0)
	c.cache.Version = "1"
	return c
}

type PackageConfigV1 struct {
	ComVer         string     `toml:"cjc-version"`
	CompilerOption string     `toml:"compiler-option"`
	Description    string     `toml:"description"`
	Name           string     `toml:"name"`
	Version        string     `toml:"version"`
	ScrPath        string     `toml:"src-dir"`
	TargetPath     string     `toml:"target-dir"`
	OutputType     OutputType `toml:"output-type"`
	Backend        string     `toml:"backend"`
	path           string
	Authors        []string
	cache          struct {
		XMLName     xml.Name `xml:"Cache"`
		Version     string   `xml:"version,attr"`
		CompilerSet SDKInfo  `xml:"cacheCompiler"`
	}
	Scripts map[string]string `toml:"scripts"`
	backend BackendProjectConfigV1
}

type PackageConfig PackageConfigV1

func (c *PackageConfigV1) LoadFromDir(p string) error {
	// 设置 path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	c.path = absPath

	err = c.backend.LoadFromDir(c.path)
	if err != nil {
		return err
	}
	c.backend.ToPackageConfig(c)

	f, err := os.ReadFile(path.Join(absPath, "cjpm.toml"))
	if err != nil {
		return err
	}

	t := NewPackageConfigV1()
	_, err = toml.NewDecoder(bytes.NewReader(f)).Decode(t)
	if err != nil {
		return err
	}

	c.Scripts = t.Scripts

	// 加载 project.lock
	xmlFile, err := os.ReadFile(path.Join(p, CacheFileName))
	if err != nil {
		return err
	}
	err = xml.Unmarshal(xmlFile, &c.cache)
	if err != nil {
		return err
	}
	return nil
}

func (c *PackageConfigV1) LoadFromBackend() error {
	err := c.backend.LoadFromDir(c.path)
	return err
}

func (c *PackageConfigV1) GenerateFromInitConfig(config *InitConfig) {
	c.path = config.Path
	//c.backend = backend
	c.Name = config.Name
	c.ComVer = config.SDK.Ver
	c.Description = config.Description
	c.Authors = []string{config.Authors}
	c.Version = config.Version
	c.cache.CompilerSet = *config.SDK
	c.OutputType = config.Output
	//c.SetBackend(backend)
	//log.Info(c.Base)
}

// SetBackend 方法用于设置 PackageConfigV0 结构体中的 backend 字段
func (c *PackageConfigV1) SetBackend(config BackendProjectConfigV1) {
	// 将传入的 config 参数赋值给 c.backend 字段
	c.backend = config
}

// WriteToDisk 写入配置
func (c *PackageConfigV1) WriteToDisk() error {
	c.SyncToBackendConfig()

	buf, err := c.backend.ToBytes()
	if err != nil {
		return err
	}

	err = toml.NewEncoder(buf).Encode(map[string]map[string]string{"scripts": c.Scripts})
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(c.path, "./cjpm.toml"), buf.Bytes(), 0755)
	if err != nil {
		return err
	}

	// 将 c.cache 序列化为 XML 格式，并使用空格和制表符进行缩进
	data, err := xml.MarshalIndent(c.cache, " ", "   ")
	// 如果序列化过程中发生错误，记录错误并返回
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	// 将序列化后的数据写入到 project.lock 文件中
	err = os.WriteFile(path.Join(c.path, CacheFileName), data, 0644)
	// 如果写入过程中发生错误，记录错误并返回
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (c *PackageConfigV1) SyncToBackendConfig() {
	c.backend.GenerateFromProjectConfig(*c)
}

// CheckCache 检查项目缓存是否有效
// 当缓存中的编译器无效时
func (c *PackageConfigV1) CheckCache() error {
	if !c.cache.CompilerSet.CheckIsHave() {
		return yError.NewNotFoundSDKErr(c.cache.CompilerSet.Path)
	}
	return nil
}

func (c *PackageConfigV1) ResetCache(p string) error {
	sdk, err := NewSDKInfo(p)
	if err != nil {
		return err
	}
	if sdk.Ver != c.ComVer {
		log.Warn("编译器版本与项目设置不匹配。")
	}
	c.cache.CompilerSet = *sdk
	return nil
}

func (c *PackageConfigV1) GetCacheSDK() *SDKInfo {
	return &c.cache.CompilerSet
}

func NewPackageConfigV1() *PackageConfigV1 {
	c := new(PackageConfigV1)
	c.OutputType = EXECUTABLE
	c.Scripts = map[string]string{}
	return c
}

func (c *PackageConfigV1) FindScript(key string) string {
	for k, v := range c.Scripts {
		if k == key {
			return v
		}
	}
	return ""
}

type OutputType = string

const (
	EXECUTABLE = "executable" // 可执行文件
	STATIC     = "static"     // 静态链接库
	DYNAMIC    = "dynamic"    // 动态链接库
)

var DefaultInitConfig = InitConfig{
	Path:        "./",
	Name:        "exampleProject",
	Authors:     "You <you@example.com>",
	Version:     "1.0.0",
	Description: "A demo",
}

type BuildOptions struct {
	ProjectPath   string
	IsRelease     bool
	RunAfterBuild bool
	backendOpt    BackendBuildOptions
}

func (opts *BuildOptions) SyncToBackendOpt() error {
	if opts.backendOpt == nil {
		return yError.NewNoBackendError("Not found build backend")
	}
	opts.backendOpt.RewriteFromBuildOptions(opts)
	return nil
}

// MakeBackendShellArgs 获取编译后端对应command
func (opts *BuildOptions) MakeBackendShellArgs() ([]string, error) {
	err := opts.SyncToBackendOpt()
	if err != nil {
		return nil, err
	}
	return opts.backendOpt.ToShellArgs(), nil
}

func NewBuildOptions() *BuildOptions {
	return &BuildOptions{}
}

func (opts *BuildOptions) SetBackend(options BackendBuildOptions) {
	opts.backendOpt = options
}

// SetBackendOptions 设置编译后端配置
func (opts *BuildOptions) SetBackendOptions(bOptions BackendBuildOptions) error {
	opts.backendOpt = bOptions
	err := opts.SyncToBackendOpt()
	if err != nil {
		return err
	}
	return nil
}

// GetBackendOptions 获取编译后端配置
func (opts *BuildOptions) GetBackendOptions() *BackendBuildOptions {
	err := opts.SyncToBackendOpt()
	if err != nil {
		return nil
	}
	return &opts.backendOpt
}

// GetOutputPath 获取后端构建产物的相对位置
func (opts *BuildOptions) GetOutputPath() string {
	return opts.backendOpt.GetOutputPath()
}
