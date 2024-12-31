package index

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/kirsle/configdir"
	"os"
	"path"
	"strings"
)

type indexType = int8

// GetPackageIndexPath
// Get the path of package
func GetPackageIndexPath(name string) []string {
	switch len(name) {
	case 1:
		return []string{"1"}
	case 2:
		return []string{"2"}
	case 3:
		return []string{"3", string(name[0])}
	default:
		return []string{string(name[0:2]), string(name[2:3])}
	}
}

const (
	GitIndex   = iota
	LocalIndex = iota
)

type Index struct {
	url  *string
	typ  indexType
	name string
}

// UpdateCache
//
//	@Description: 更新缓存
//	@receiver i
//	@return error
func (i Index) UpdateCache() error {
	if i.typ == LocalIndex {
		return nil
	}
	p := configdir.LocalCache("fuxo", "index", i.name)
	if r, err := git.PlainOpen(p); err != nil {

		err := os.RemoveAll(p)
		if err != nil {
		}

		_, err = git.PlainClone(configdir.LocalCache(p), false, &git.CloneOptions{URL: *i.url, Progress: os.Stdout, SingleBranch: true})
		if err != nil {
			return err
		}

	} else {
		err := r.Fetch(&git.FetchOptions{Force: true, RemoteURL: *i.url})
		if err != nil {
			return err
		}
	}

	return nil
}

// FindPackage 查找某个包
func (i Index) FindPackage(name string) (*PackageIndexes, error) {
	is := PackageIndexes{}
	p := path.Join(append(GetPackageIndexPath(name), name)...)
	if i.typ == GitIndex {
		f, err := os.ReadFile(path.Join(configdir.LocalCache("fuxo", "index", i.name), p))
		if err != nil {
			return nil, fmt.Errorf("do not found this package")
		}
		//strings.Split(string(f),"\n")
		for _, v := range strings.Split(string(f), "\n") {
			if v != "" {
				p, err := NewPackageIndexFromJson([]byte(v))
				if err != nil {
					continue
				}
				is = append(is, *p)
			}
		}
	}
	//TODO Support Locale
	return &is, fmt.Errorf("not support this index type")
}

func NewIndex(url *string, typ indexType, name string) Index {
	return Index{url: url, typ: typ, name: name}
}

type PackageDepend struct {
	Name     string `json:"name"`
	ReqVer   string `json:"req"`
	Registry string `json:"registry"`
}

type PackageIndex struct {
	Name             string          `json:"name"`
	Version          string          `json:"vers"`
	Depends          []PackageDepend `json:"deps"`
	Checksum         string          `json:"cksum"`
	Yanked           bool            `json:"yanked"`
	CompilerVersions []string        `json:"cjcv"`
}

// ToPackageDepend 从索引转换为依赖
func (p PackageIndex) ToPackageDepend() *PackageDepend {
	d := new(PackageDepend)
	d.Name = p.Name
	d.ReqVer = p.Version
	d.Registry = ""
	return d
}

func NewPackageIndexFromJson(b []byte) (*PackageIndex, error) {
	p := new(PackageIndex)
	err := json.Unmarshal(b, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

type PackageIndexes = []PackageIndex

// FindVersion
//
//	@Description: 查找索引表中是否有符合版本要求的项
//	@param s 对应索引表
//	@param ver 版本
func FindVersion(s PackageIndexes, ver string, useYanked bool) (PackageIndex, error) {
	for _, v := range s {
		if v.Yanked && !useYanked {
			continue
		}
		if v.Version == ver {
			return v, nil
		}
	}
	return PackageIndex{}, fmt.Errorf("version %s not found", ver)
}
