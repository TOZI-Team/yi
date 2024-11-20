package types

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var versionMatch = regexp.MustCompile("[0-9]+\\.[0-9]+\\.[0-9]+")

type SDKInfo struct {
	Ver  Version `xml:"version" yaml:"Version" json:"version"`
	Path string  `xml:"path,attr" yaml:"Path" json:"path"`
	Note string  `xml:"note" yaml:"Note" json:"note"`
}

type SDKRunOptions struct {
	shell        string
	linkToStdout bool
	linkToStderr bool
}

// CheckIsHave 检查是否真实存在
func (i SDKInfo) CheckIsHave() bool {
	// TODO fix
	if path.IsAbs(i.Path) && path.IsAbs(path.Join(i.Path, "bin/")) && path.IsAbs(path.Join(i.Path, "tools/bin/")) {
		return true
	}
	return false
}

// RunCommand 在环境下运行命令
func (i SDKInfo) RunCommand(command []string, p string) error {
	//command = append([]string{i.GetActivityEnvScript()}, command...)
	log.Info("Run command: ", "bash", "-c", strings.Join([]string{strings.Join(i.GetActivityEnvScript(), " "), strings.Join(command, " ")}, "&&"))
	cmd := exec.Command("bash", "-c", strings.Join([]string{strings.Join(i.GetActivityEnvScript(), " "), strings.Join(command, " ")}, "&&"))
	cmd.Dir = p
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// BuildProject 构建项目
func (i SDKInfo) BuildProject(p string, options BuildOptions) (string, error) {
	command, err := options.MakeBackendShellArgs()
	if err != nil {
		return "", err
	}
	err = i.RunCommand(command, p)
	if err != nil {
		return "", err
	}
	return path.Join(p, options.GetOutputPath()), nil
}

func (i SDKInfo) GetActivityEnvScript() []string {
	s := []string{"source", fmt.Sprintf("\"%s\"", path.Join(i.Path, "./envsetup.sh"))}
	return s
}

func NewSDKInfo(p string) (*SDKInfo, error) {
	c := exec.Command(path.Join(p, "./bin/cjc"), "-v")
	output, err := c.Output()
	log.Debugf("Cjc version output: \n===\n%s\n===", string(output))
	if err != nil {
		return nil, err
	}
	output = versionMatch.Find(output)
	return &SDKInfo{Path: p, Ver: Version(output)}, nil
}
