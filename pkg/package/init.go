package cj_package

import "regexp"

var PackageR = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")

// IsCjPackageName 判断是否为合规仓颉包名
func IsCjPackageName(name string) bool {
	return PackageR.MatchString(name)
}
