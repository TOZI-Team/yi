package types

type BackendBuildOptions interface {
	ToShellArgs() []string
	RewriteFromBuildOptions(*BuildOptions)
	//
	// GetOutputPath
	//  @Description: 获取构建生产物相对项目的位置
	//
	GetOutputPath() string
}

// BackendProjectConfig 编译后端项目配置
type BackendProjectConfig interface {
	GenerateFromProjectConfig(config PackageConfig)
	GenerateFromPackageConfig(config PackageConfig)
	//
	// WriteConfigToDir
	//  @Description: 将后端设置写入磁盘。
	//  @param string 文件夹路径
	//  @return error
	//
	WriteConfigToDir(string) error
}
