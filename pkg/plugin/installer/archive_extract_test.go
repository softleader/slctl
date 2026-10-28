package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type archiveEntry struct {
	name     string
	body     string
	mode     os.FileMode
	link     string
	hardlink bool
}

// Fixtures use independent standard-library writers, not the extraction library.
func writeTestArchive(t *testing.T, format string, entries ...archiveEntry) string {
	t.Helper()
	name := filepath.Join(t.TempDir(), "plugin."+format)
	out, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if format == "zip" {
		zw := zip.NewWriter(out)
		for _, entry := range entries {
			header := &zip.FileHeader{Name: entry.name, Method: zip.Deflate}
			header.SetMode(entry.mode)
			w, err := zw.CreateHeader(header)
			if err != nil {
				t.Fatal(err)
			}
			body := entry.body
			if entry.mode&os.ModeSymlink != 0 {
				body = entry.link
			}
			if _, err := io.WriteString(w, body); err != nil {
				t.Fatal(err)
			}
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		return name
	}
	var w io.Writer = out
	if format == "tar.gz" {
		gz := gzip.NewWriter(out)
		defer func() {
			if err := gz.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		w = gz
	}
	tw := tar.NewWriter(w)
	for _, entry := range entries {
		header := &tar.Header{Name: entry.name, Mode: int64(entry.mode.Perm()), Size: int64(len(entry.body)), Linkname: entry.link}
		switch {
		case entry.hardlink:
			header.Typeflag = tar.TypeLink
		case entry.mode&os.ModeSymlink != 0:
			header.Typeflag = tar.TypeSymlink
		case entry.mode.IsDir():
			header.Typeflag = tar.TypeDir
		case entry.mode&os.ModeNamedPipe != 0:
			header.Typeflag = tar.TypeFifo
		default:
			header.Typeflag = tar.TypeReg
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(tw, entry.body); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return name
}

func TestExtractPreservesFilesAndExecutableMode(t *testing.T) {
	for _, format := range []string{"zip", "tar", "tar.gz"} {
		t.Run(format, func(t *testing.T) {
			source := writeTestArchive(t, format,
				archiveEntry{name: "./", mode: os.ModeDir | 0755},
				archiveEntry{name: "./bin/run", body: "#!/bin/sh\necho plugin\n", mode: 0755},
				archiveEntry{name: "empty/", mode: os.ModeDir | 0755},
			)
			destination := filepath.Join(t.TempDir(), "extracted")
			if err := extract(source, destination); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(destination, "bin/run"))
			if err != nil || string(data) != "#!/bin/sh\necho plugin\n" {
				t.Fatalf("extracted content = %q, error = %v", data, err)
			}
			info, err := os.Stat(filepath.Join(destination, "bin/run"))
			if err != nil {
				t.Fatal(err)
			}
			if runtime.GOOS != "windows" && info.Mode().Perm() != 0755 {
				t.Errorf("executable mode = %o, want 755", info.Mode().Perm())
			}
			info, err = os.Stat(filepath.Join(destination, "empty"))
			if err != nil || !info.IsDir() {
				t.Fatalf("empty directory missing: %v", err)
			}
		})
	}
}

func TestExtractRejectsUnsafeEntries(t *testing.T) {
	for _, format := range []string{"zip", "tar", "tar.gz"} {
		for _, name := range []string{"../outside", "../extracted-sibling/outside", "/absolute", `..\outside`, `C:\outside`, "file:stream"} {
			t.Run(format+"/"+name, func(t *testing.T) {
				root := t.TempDir()
				outside := filepath.Join(root, "outside")
				if err := os.WriteFile(outside, []byte("untouched"), 0644); err != nil {
					t.Fatal(err)
				}
				source := writeTestArchive(t, format, archiveEntry{name: name, body: "malicious", mode: 0644})
				if err := extract(source, filepath.Join(root, "extracted")); err == nil {
					t.Fatal("unsafe archive entry was accepted")
				}
				data, err := os.ReadFile(outside)
				if err != nil || string(data) != "untouched" {
					t.Fatalf("file outside extraction directory changed: %q, %v", data, err)
				}
				if _, err := os.Stat(filepath.Join(root, "extracted-sibling/outside")); !os.IsNotExist(err) {
					t.Fatalf("archive escaped to sibling directory: %v", err)
				}
			})
		}
	}
}

func TestExtractRejectsLinksAndSpecialFiles(t *testing.T) {
	for _, tc := range []struct {
		name   string
		format string
		entry  archiveEntry
	}{
		{"zip symlink", "zip", archiveEntry{name: "link", mode: os.ModeSymlink | 0777, link: "../outside"}},
		{"tar symlink", "tar", archiveEntry{name: "link", mode: os.ModeSymlink | 0777, link: "../outside"}},
		{"tar hardlink", "tar", archiveEntry{name: "link", mode: 0644, link: "../outside", hardlink: true}},
		{"tar fifo", "tar", archiveEntry{name: "pipe", mode: os.ModeNamedPipe | 0644}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := writeTestArchive(t, tc.format, tc.entry)
			destination := t.TempDir()
			if err := extract(source, destination); err == nil {
				t.Fatal("link or special file was accepted")
			}
			if _, err := os.Lstat(filepath.Join(destination, tc.entry.name)); !os.IsNotExist(err) {
				t.Fatalf("unsafe entry was created: %v", err)
			}
		})
	}
}

func TestExtractDoesNotOverwriteExistingFiles(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			source := writeTestArchive(t, format,
				archiveEntry{name: "metadata.yaml", body: "first", mode: 0644},
				archiveEntry{name: "metadata.yaml", body: "second", mode: 0644},
			)
			destination := t.TempDir()
			if err := extract(source, destination); err == nil {
				t.Fatal("duplicate file was accepted")
			}
			data, err := os.ReadFile(filepath.Join(destination, "metadata.yaml"))
			if err != nil || string(data) != "first" {
				t.Fatalf("first entry overwritten: %q, %v", data, err)
			}
		})
	}
}

func TestExtractCannotFollowExistingSymlinkOutsideRoot(t *testing.T) {
	outside := t.TempDir()
	destination := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(destination, "linked")); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf("symlinks unavailable: %v", err)
		}
		t.Fatal(err)
	}
	source := writeTestArchive(t, "zip", archiveEntry{name: "linked/escaped", body: "malicious", mode: 0644})
	if err := extract(source, destination); err == nil {
		t.Fatal("extraction followed a symlink outside its root")
	}
	if _, err := os.Stat(filepath.Join(outside, "escaped")); !os.IsNotExist(err) {
		t.Fatalf("file created outside extraction root: %v", err)
	}
}

func TestExtractIdentifiesArchiveContents(t *testing.T) {
	source := writeTestArchive(t, "zip", archiveEntry{name: "metadata.yaml", body: "name: plugin\n", mode: 0644})
	unnamed := filepath.Join(t.TempDir(), "download")
	if err := os.Rename(source, unnamed); err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir()
	if err := extract(unnamed, destination); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(destination, "metadata.yaml"))
	if err != nil || string(data) != "name: plugin\n" {
		t.Fatalf("archive without extension was not extracted: %q, %v", data, err)
	}
}

func TestExtractRejectsInvalidArchive(t *testing.T) {
	source := filepath.Join(t.TempDir(), "invalid.zip")
	if err := os.WriteFile(source, []byte("not an archive"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := extract(source, t.TempDir()); err == nil {
		t.Fatal("invalid archive accepted")
	}
}
