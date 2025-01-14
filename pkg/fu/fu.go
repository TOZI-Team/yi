package fu

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	devlog "yi/log"
)

type Fu struct {
	path string
}

const (
	EXECUTABLE = "executable" // 可执行文件
	STATIC     = "static"     // 静态链接库
	DYNAMIC    = "dynamic"    // 动态链接库
)

// Unpack 解压
func (f Fu) Unpack(dir string) error {
	err := extractTarGz(f.path, dir)
	if err != nil {
		return err
	}
	return nil
}

// LoadFromDisk 从磁盘获取包
func LoadFromDisk(path string) *Fu {
	return &Fu{path: path}
}

// LoadFromURL 从URL获取包
//
// Deprecated: 终止支持
func LoadFromURL(url string, name string, version string) (*Fu, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(r.Body)

	p := configdir.LocalCache("fuxo", "fu", name+"-"+version+".fu")
	err = getter.GetFile(p, url)
	if err != nil {
		return nil, err
	}
	return &Fu{path: p}, nil
}

func extractTarGz(src, dst string) error {
	// 打开 .tar.gz 文件
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)

	// 创建一个 gzip 读取器
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func(gzr *gzip.Reader) {
		err := gzr.Close()
		if err != nil {
			log.Error("Failed to close gzip reader")
		}
	}(gzr)

	// 创建一个 tar 读取器
	tr := tar.NewReader(gzr)

	// 遍历 tar 文件中的每个文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // 到达文件末尾
		}
		if err != nil {
			return err
		}

		// 构建目标文件的完整路径
		target := filepath.Join(dst, header.Name)

		// 根据文件类型进行处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 如果是目录，创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// 如果是普通文件，创建文件并写入内容
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				err := outFile.Close()
				if err != nil {
					return err
				}
				return err
			}
			err = outFile.Close()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %v", header.Typeflag)
		}
	}

	return nil
}

// ExtractTarGzFromURL 从 URL 下载并解压 .tar.gz 文件
func extractTarGzFromURL(url, dst string) error {
	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			devlog.DevLog.Errorf("Failed to close response body")
		}
	}(resp.Body)

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 创建 gzip 读取器
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func(gzr *gzip.Reader) {
		err := gzr.Close()
		if err != nil {
			devlog.DevLog.Errorf("Failed to close gzip reader: %s", err)
		}
	}(gzr)

	// 创建 tar 读取器
	tr := tar.NewReader(gzr)

	// 遍历 tar 文件中的每个文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // 到达文件末尾
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// 构建目标文件的完整路径
		target := filepath.Join(dst, header.Name)

		// 根据文件类型进行处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 如果是目录，创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// 如果是普通文件，创建文件并写入内容
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				err := outFile.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("failed to write file content: %w", err)
			}
			err = outFile.Close()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type: %v", header.Typeflag)
		}
	}

	return nil
}
