package installer

import (
	"bytes"
	"compress/flate"
	"github.com/mholt/archiver"
	"github.com/softleader/slctl/pkg/slpath"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveInstaller_Install(t *testing.T) {
	home, err := ioutil.TempDir("", "sl_home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(home)
	b := bytes.NewBuffer(nil)
	hh := slpath.Home(home)
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}
	arcName := "test.zip"
	arcPath := filepath.Join(hh.String(), arcName)
	arcSrc := filepath.Join(hh.String(), "file.txt")
	if err := ioutil.WriteFile(arcSrc, []byte("test"), 0744); err != nil {
		t.Error(err)
		return
	}
	if err := z.Archive([]string{arcSrc}, arcPath); err != nil {
		t.Error(err)
		return
	}

	i, err := newArchiveInstaller(b, arcPath, hh)
	if err != nil {
		t.Error(err)
		return
	}

	downloaded, err := i.downloader.download()
	if err != nil {
		t.Error(err)
		return
	}

	dst := filepath.Join(hh.CacheArchives(), arcName)
	if downloaded != dst {
		t.Errorf("expected downloaded dir to be %s", dst)
	}

	if !isLocalReference(arcPath) {
		t.Errorf("expected downloaded dir to be a legal local reference: %s", dst)
	}
}