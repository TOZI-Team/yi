package backendInterface

import t "Yi/pkg/types"

type BackendBuildOptions interface {
	ToShellArgs() []string
	RewriteFromBuildOptions(*t.BuildOptions)
	//
	// GetOutputPath
	//  @Description: 获取构建生产物相对项目的位置
	//
	GetOutputPath() string
}

type BackendProjectConfig interface {
	GenerateFromProjectConfig()
	WriteConfigToDir(string) error
}
