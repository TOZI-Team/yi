package types

import "bytes"

type BackendBuildOptions interface {
	ToShellArgs() []string
	RewriteFromBuildOptions(*BuildOptions)
	//
	// GetOutputPath
	//  @Description: 获取构建生产物相对项目的位置
	//
	GetOutputPath() string
}

// BackendProjectConfigV0 编译后端项目配置
type BackendProjectConfigV0 interface {
	GenerateFromProjectConfig(config PackageConfigV0)
	GenerateFromPackageConfig(config PackageConfigV0)
	// WriteConfigToDir
	//  @Description: 将后端设置写入磁盘。
	//  @param string 文件夹路径
	//  @return error
	//
	WriteConfigToDir(string) error
}

// BackendProjectConfigV1 后端项目配置
type BackendProjectConfigV1 interface {
	GenerateFromProjectConfig(config PackageConfigV1)
	GenerateFromPackageConfig(config PackageConfigV1)
	ToBytes() (*bytes.Buffer, error)
	LoadFromDir(string) error
	ToPackageConfig(p *PackageConfigV1)
}

type BackendProjectConfig BackendProjectConfigV1
