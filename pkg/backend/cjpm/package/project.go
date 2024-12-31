package cjpmPackage

import cj_package "yi/pkg/package"

const version = "0.1.0"

type CJPMProjectBackend struct {
}

func (b CJPMProjectBackend) BackendInfo() string {
	return "CJPM Backend Ver " + version
}

func (b CJPMProjectBackend) MakeConfig(p *cj_package.Package, dir string) error {
	//TODO 完善
	return nil
}

func (b CJPMProjectBackend) MakeBuildArgs(options *cj_package.BuildOptions) (*cj_package.BuildResult, error) {
	//TODO 待完善
	return nil, nil
}
