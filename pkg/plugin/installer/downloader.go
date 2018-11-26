package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/slpath"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type downloader interface {
	download() (string, error)
}

// the final dst will be home.CacheArchives() + dstDir
func newDownloader(source interface{}, home slpath.Home, dstDir string) (downloader, error) {
	dst := filepath.Join(home.CacheArchives(), dstDir)
	switch src := source.(type) {
	case string:
		return urlDownloader{
			dst: dst,
			url: src,
		}, nil
	case io.ReadCloser:
		return rcDownloader{
			dst: dst,
			rc:  src,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type '%s' of source to new downloader", src)
	}
}

type urlDownloader struct {
	dst string
	url string
}

func (d urlDownloader) download() (string, error) {
	if _, err := os.Stat(d.dst); os.IsExist(err) {
		os.Remove(d.dst)
	}
	out, err := os.Create(d.dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	resp, err := http.Get(d.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}
	return d.dst, nil
}

type rcDownloader struct {
	dst string
	rc  io.ReadCloser
}

func (d rcDownloader) download() (string, error) {
	defer d.rc.Close()
	if _, err := os.Stat(d.dst); os.IsExist(err) {
		os.Remove(d.dst)
	}
	out, err := os.Create(d.dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = io.Copy(out, d.rc); err != nil {
		return "", err
	}
	return d.dst, nil
}
