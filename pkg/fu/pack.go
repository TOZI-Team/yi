package fu

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
)

import "github.com/gobwas/glob"

var defaultFuIncludeFile []glob.Glob
var defaultFuExcludeFile []glob.Glob = []glob.Glob{glob.MustCompile(".git/**"), glob.MustCompile("target/**"), glob.MustCompile("**/*.exe")}

// 将文件添加到tar归档中
func addFileToTar(tw *tar.Writer, path, baseDir string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("Error closing file: ", err)
		}
	}(file)

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// 创建tar头信息
	header, err := tar.FileInfoHeader(info, info.Name())
	header.Uname = "fuxo"
	header.Mode = int64(0755)
	header.Gname = "fuxo"
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(baseDir, path)
	if err != nil {
		return err
	}

	// 使用相对路径作为tar中的文件名
	header.Name = filepath.ToSlash(rel)

	// 写入头信息
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	// 写入文件内容
	_, err = io.Copy(tw, file)
	return err
}

// 创建.tar.gz压缩包
func createTarGz(outputPath string, sourceDir string, selector *FileSelector) error {
	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("Error closing file: ", err)
		}
	}(file)

	// 创建gzip写入器
	gw := gzip.NewWriter(file)
	defer func(gw *gzip.Writer) {
		err := gw.Close()
		if err != nil {
			log.Error("Error closing gzip writer: ", err)
		}
	}(gw)

	// 创建tar写入器
	tw := tar.NewWriter(gw)
	defer func(tw *tar.Writer) {
		err := tw.Close()
		if err != nil {
			log.Error("Error closing tar writer: ", err)
		}
	}(tw)

	// 遍历源目录
	return filepath.Walk(sourceDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录，只添加文件
		if !info.IsDir() {
			if defaultFileSelector.Match(path.Join(p, info.Name())) && selector.Match(path.Join(p, info.Name())) {
				if err := addFileToTar(tw, p, filepath.Dir(sourceDir)); err != nil {
					return err
				}
			}

		}

		return nil
	})
}

// NewFuFromDir 从文件夹创建包
// output 应当为文件夹
func NewFuFromDir(dir string, output string, selector *FileSelector) (*Fu, error) {
	f := new(Fu)
	//p, err := cjpackage.LoadPackageFromDir(dir, false)
	//if err != nil {
	//	return nil, err
	//}
	//f.

	// 创建压缩包
	err := createTarGz(path.Join(output, fmt.Sprintf("%s.fu", path.Base(dir))), dir, selector)
	if err != nil {
		return nil, err
	}
	f.path = dir

	return f, nil
}

// FileSelector 一个文件筛选器
type FileSelector struct {
	include []glob.Glob
	exclude []glob.Glob
}

func (s *FileSelector) AddInclude(i string) error {
	g, err := glob.Compile(i)
	if err != nil {
		return err
	}

	s.include = append(s.include, g)
	return nil
}

func (s *FileSelector) AddExclude(i string) error {
	g, err := glob.Compile(i)
	if err != nil {
		return err
	}
	s.exclude = append(s.exclude, g)
	return nil
}

func (s *FileSelector) Match(path string) bool {
	for _, g := range s.include {
		if g.Match(path) {
			return true
		}
	}

	for _, g := range s.exclude {
		if g.Match(path) {
			return false
		}
	}
	return true
}

func NewFileSelector() *FileSelector {
	return &FileSelector{}
}

var defaultFileSelector FileSelector

func init() {
	defaultFileSelector = FileSelector{include: defaultFuIncludeFile, exclude: defaultFuExcludeFile}
}
