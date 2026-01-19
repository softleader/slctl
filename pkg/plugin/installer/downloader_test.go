package installer

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/softleader/slctl/pkg/paths"
)

func TestURLDownloader_Download(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "content")
	}))
	defer server.Close()

	tempHome, _ := os.MkdirTemp("", "sl-home-download")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.CacheArchives(), 0755)

	d := newURLDownloader(server.URL, hh, "test.zip")
	path, err := d.download()
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	expected := filepath.Join(hh.CacheArchives(), "test.zip")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestReaderDownloader_Download(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-reader")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.CacheArchives(), 0755)

	content := "reader content"
	var r io.Reader = strings.NewReader(content)
	
	d := newReaderDownloader(&r, hh, "test-reader.txt")
	path, err := d.download()
	if err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("expected %s, got %s", content, string(got))
	}
}

func TestReadCloserDownloader_Download(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-readcloser")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.CacheArchives(), 0755)

	content := "readcloser content"
	rc := io.NopCloser(strings.NewReader(content))
	
	d := newReadCloserDownloader(&rc, len(content), hh, "test-rc.txt")
	path, err := d.download()
	if err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("expected %s, got %s", content, string(got))
	}
}


