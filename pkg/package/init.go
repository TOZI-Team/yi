package cj_package

import (
	"fmt"
	"golang.org/x/mod/semver"
	"regexp"
)

var packageR *regexp.Regexp
var packageVersionR *regexp.Regexp

// IsCjPackageName 判断是否为合规仓颉包名
func IsCjPackageName(name string) bool {
	if packageR == nil {
		packageR = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")
	}
	return packageR.MatchString(name)
}

// IsCjPackageVersion 判断是否为有效版本
func IsCjPackageVersion(version string) bool {
	if packageVersionR == nil {
		packageVersionR = regexp.MustCompile("^(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$")
	}
	return packageVersionR.MatchString(version)
}

func ToCJPMPackageVersion(version string) string {
	return fmt.Sprintf("%s.%s", semver.Major(version), semver.MajorMinor(version))
}
