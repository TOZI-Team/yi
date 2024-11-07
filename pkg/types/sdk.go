package types

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

type SDKInfo struct {
	Ver  Version
	Path string
	Note string
}

type SDKRunOptions struct {
	shell        string
	linkToStdout bool
	linkToStderr bool
}

// CheckIsHave 检查是否真实存在
func (i SDKInfo) CheckIsHave() bool {
	if path.IsAbs(i.Path) && path.IsAbs(path.Join(i.Path, "bin/cjc")) && path.IsAbs(path.Join(i.Path, "tools/bin/")) {
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

// 构建项目
func (i SDKInfo) BuildProject(p string, options BuildOptions) (string, error) {
	command := options.MakeBackendShellArgs()
	err := i.RunCommand(command, p)
	if err != nil {
		return "", err
	}
	return path.Join(p, options.GetOutputPath()), nil
}

func (i SDKInfo) GetActivityEnvScript() []string {
	return []string{"source", fmt.Sprintf("\"%s\"", path.Join(i.Path, "./envsetup.sh"))}
}

func NewSDKInfo(path string) *SDKInfo {
	// TODO 根据 path 获取版本
	return &SDKInfo{Path: path, Ver: "0.53.13"}
}
