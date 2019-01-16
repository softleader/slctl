package installer

import (
	"github.com/softleader/slctl/pkg/paths"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type downloader interface {
	download() (string, error)
}

func newURLDownloader(url string, home paths.Home, dstDir string) *urlDownloader {
	return &urlDownloader{
		dst: filepath.Join(home.CacheArchives(), dstDir),
		url: url,
	}
}

type urlDownloader struct {
	dst string
	url string
}

func (d *urlDownloader) download() (string, error) {
	if _, err := os.Stat(d.dst); os.IsExist(err) {
		os.Remove(d.dst)
	}
	if err := os.MkdirAll(filepath.Dir(d.dst), 0755); err != nil {
		return "", err
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

	bar := pb.New64(resp.ContentLength).SetUnits(pb.U_BYTES).Start()
	defer bar.Finish()
	if _, err = io.Copy(out, bar.NewProxyReader(resp.Body)); err != nil {
		return "", err
	}
	return d.dst, nil
}

func newReaderDownloader(r *io.Reader, home paths.Home, dstDir string) *readerDownloader {
	return &readerDownloader{
		dst: filepath.Join(home.CacheArchives(), dstDir),
		r:   r,
	}
}

type readerDownloader struct {
	dst string
	r   *io.Reader
}

func (d *readerDownloader) download() (string, error) {
	if _, err := os.Stat(d.dst); os.IsExist(err) {
		os.Remove(d.dst)
	}
	if err := os.MkdirAll(filepath.Dir(d.dst), 0755); err != nil {
		return "", err
	}
	out, err := os.Create(d.dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = io.Copy(out, *d.r); err != nil {
		return "", err
	}
	return d.dst, nil
}

func newReadCloserDownloader(rc *io.ReadCloser, home paths.Home, dstDir string) *readCloserDownloader {
	return &readCloserDownloader{
		dst: filepath.Join(home.CacheArchives(), dstDir),
		rc:  rc,
	}
}

type readCloserDownloader struct {
	dst string
	rc  *io.ReadCloser
}

func (d *readCloserDownloader) download() (string, error) {
	defer (*d.rc).Close()
	if _, err := os.Stat(d.dst); os.IsExist(err) {
		os.Remove(d.dst)
	}
	if err := os.MkdirAll(filepath.Dir(d.dst), 0755); err != nil {
		return "", err
	}
	out, err := os.Create(d.dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = io.Copy(out, *d.rc); err != nil {
		return "", err
	}
	return d.dst, nil
}
