package plugin

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var registeredCreators = []creator{
	golang{},
	java{},
	nodejs{},
}
var Creators = func() (m map[string]creator) {
	m = make(map[string]creator, len(registeredCreators))
	for _, c := range registeredCreators {
		m[reflect.TypeOf(c).Name()] = c
	}
	return
}()

type creator interface {
	command(plugin *Metadata) string
	files(plugin *Metadata, pluginDir string) []file
}

func Create(lang string, plugin *Metadata, dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return path, err
	}

	if fi, err := os.Stat(path); err != nil {
		return path, err
	} else if !fi.IsDir() {
		return path, fmt.Errorf("no such directory %s", path)
	}

	pdir := filepath.Join(path, plugin.Name)
	if err := mkdir(pdir); err != nil {
		return pdir, err
	}

	creator, found := Creators[lang]
	if !found {
		return pdir, fmt.Errorf(`unsupported creating %s template.
You might need to run 'slctl plugin create langs'`, lang)
	}

	plugin.Command = creator.command(plugin)
	files := creator.files(plugin, pdir)
	files = append(files, marshal{
		path: filepath.Join(pdir, MetadataFileName),
		in:   plugin,
	})

	for _, file := range files {
		if err := save(file); err != nil {
			return pdir, err
		}
	}
	return pdir, nil
}

func mkdir(path string) error {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		return fmt.Errorf("file %s already exists and is not a directory", path)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

type file interface {
	filepath() string
	content() ([]byte, error)
}

type marshal struct {
	path string
	in   interface{}
}

func (u marshal) filepath() string {
	return u.path
}

func (u marshal) content() ([]byte, error) {
	return yaml.Marshal(u.in)
}

type tpl struct {
	path     string
	in       interface{}
	template string
}

func (u tpl) filepath() string {
	return u.path
}

func (u tpl) content() ([]byte, error) {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Funcs(funcMap).Parse(u.template))
	if err := parsed.Execute(&buf, u.in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func save(file file) (err error) {
	mkdir(filepath.Dir(file.filepath())) // ensure parent exist
	out, err := file.content()
	if err != nil {
		return
	}
	return ioutil.WriteFile(file.filepath(), out, 0755)
}
