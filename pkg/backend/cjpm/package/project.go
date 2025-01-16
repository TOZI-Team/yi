package cjpmPackage

import (
	"github.com/BurntSushi/toml"
	"os"
	"yi/internal/sdk"
	cjpackage "yi/pkg/package"
	t "yi/pkg/types"
)

const version = "0.1.0"

type CJPMProjectBackend struct {
}

func (b CJPMProjectBackend) BackendInfo() string {
	return "CJPM Backend Ver " + version
}

func (b CJPMProjectBackend) MakeConfig(p *cjpackage.Package, opt *cjpackage.BackendConfigOption) error {
	c := newCJPMConfigV2()
	c.Package.Name = p.GetName()
	c.Package.Version = cjpackage.ToCJPMPackageVersion(p.Config().Project.Version)
	c.Package.ComVer = p.Config().Project.CjcVersion
	c.Package.Description = ""

	for k, v := range p.Config().Dependencies.LocalPackage {
		if !v.IsLegacy {
			continue
		}
		c.Depend.CJPMLocalDepend[k] = CJPMLocalDepend{Path: v.Path}
	}

	for k, v := range p.Config().Dependencies.GitPackage {
		if !v.IsLegacy {
			continue
		}

		c.Depend.CJPMGitDepend[k] = CJPMGitDepend{URL: v.URL, Tag: v.Tag}
	}

	marshal, err := toml.Marshal(*c)
	if err != nil {
		return err
	}

	err = os.WriteFile("cjpm.toml", marshal, 0644)
	if err != nil {
		return err
	}
	//TODO 完善
	return nil
}

func (b CJPMProjectBackend) Build(options *cjpackage.BuildOptions) (*cjpackage.BuildResult, error) {
	ot := t.PackageConfigV0{}
	err := ot.LoadFromDir(options.Path)
	if err != nil {
		return nil, err
	}

	cv := ot.Base.ComVer
	byVersion, err := sdk.GlobalSDKManger.FindByVersion(cv)
	if err != nil {
		return nil, err
	}

	err = byVersion.RunCommand([]string{"cjpm", "build"}, options.Path)
	if err != nil {
		return nil, err
	}
	//TODO 待完善
	return &cjpackage.BuildResult{Success: true}, nil
}

func (b CJPMProjectBackend) Clean() error {
	return nil
}
