package cjpmPackage

import (
	"strconv"
	t "yi/pkg/types"
)

type CJPMBuildOptions struct {
	EnableIncremental bool
	EnableVerbose     bool
	EnableDebugTarget bool
	EnableCoverage    bool
	EnableLintCheck   bool
	EnableMock        bool
	EnableSkipScript  bool
	Target            string
	Member            string
	CfgPath           string
	JobNum            int8
}

func NewCJPMBuildOptions() *CJPMBuildOptions {
	return &CJPMBuildOptions{}
}

func (opt *CJPMBuildOptions) ToShellArgs() []string {
	args := []string{"cjpm", "build"}
	if opt.EnableIncremental {
		args = append(args, "--incremental")
	}
	if opt.EnableVerbose {
		args = append(args, "--verbose")
	}
	if opt.EnableDebugTarget {
		args = append(args, "-g")
	}
	if opt.EnableCoverage {
		args = append(args, "-coverage")
	}
	if opt.EnableLintCheck {
		args = append(args, "--lint")
	}
	if opt.EnableMock {
		args = append(args, "--mock")
	}
	if opt.EnableSkipScript {
		args = append(args, "--skip-script")
	}
	if opt.JobNum > 0 {
		args = append(args, "--jobs")
		args = append(args, strconv.Itoa(int(opt.JobNum)))
	}
	return args
}

func (opt *CJPMBuildOptions) GetOutputPath() string {
	if !opt.EnableDebugTarget {
		return "target/release/bin/main"
	} else {
		return "target/debug/bin/main"
	}
}

// RewriteFromBuildOptions
//
//	@Description: 从 BuildOptions 覆写配置
func (opt *CJPMBuildOptions) RewriteFromBuildOptions(options *t.BuildOptions) {
	if options.IsRelease {
		opt.EnableDebugTarget = false
	} else {
		opt.EnableDebugTarget = true
	}
}
