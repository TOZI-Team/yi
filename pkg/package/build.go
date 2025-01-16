package cj_package

//func Build(options t.BuildOptions,path string,sdk *t.SDKInfo) (string, error) {
//	args := options.MakeBackendShellArgs()
//	sdk.
//}

type BuildType uint8

const (
	STATIC = iota
	EXEC   = iota
)

type BuildOptions struct {
	IsRelease     bool
	RunAfterBuild bool
	BuildType     BuildType
	ShowOutput    bool
	SDKVersion    string
	Path          string
}

type BuildResult struct {
	Success bool
	Bins    []string
}
