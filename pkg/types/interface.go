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
	WriteConfigToDir(string) error
}
