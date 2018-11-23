package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type downloader interface {
	downloadTo(destination string) error
}

func newDownloader(source interface{}) (downloader, error) {
	switch src := source.(type) {
	case string:
		return urlDownloader{
			url: src,
		}, nil
	case io.ReadCloser:
		return rcDownloader{
			rc: src,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type '%s' of source to new downloader", src)
	}
}

type urlDownloader struct {
	url string
}

func (d urlDownloader) downloadTo(destination string) error {
	if _, err := os.Stat(destination); os.IsExist(err) {
		os.Remove(destination)
	}
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(d.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

type rcDownloader struct {
	rc io.ReadCloser
}

func (d rcDownloader) downloadTo(destination string) error {
	defer d.rc.Close()
	if _, err := os.Stat(destination); os.IsExist(err) {
		os.Remove(destination)
	}
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, out); err != nil {
		return err
	}
	return nil
}
