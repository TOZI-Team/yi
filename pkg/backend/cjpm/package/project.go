package cjpmPackage

import cjpackage "yi/pkg/package"

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

	//TODO 完善
	return nil
}

func (b CJPMProjectBackend) MakeBuildArgs(options *cjpackage.BuildOptions, opt *cjpackage.BackendConfigOption) (*cjpackage.BuildResult, error) {
	//TODO 待完善
	return nil, nil
}
