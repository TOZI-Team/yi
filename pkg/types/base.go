package types

import (
	"golang.org/x/mod/semver"
)

type Version = string

// 获取主版本
func getMajor(version Version) string {
	return semver.Major(version)
}

func getMinorMinor(version Version) string {
	return semver.MajorMinor(version)
}

func IsSemVer(version string) bool {
	return semver.IsValid(version)
}
