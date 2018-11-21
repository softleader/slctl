package plugin

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

const (
	PluginFileName = "plugin.yaml"
)

var Creators = []creator{
	golang{},
	java{},
	node{},
}

type creator interface {
	files(plugin *Metadata, pluginDir string) []file
}

func findCreator(lang string) creator {
	for _, c := range Creators {
		if reflect.TypeOf(c).Name() == lang {
			return c
		}
	}
	return nil
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

	pluginDir := filepath.Join(path, plugin.Name)
	if fi, err := os.Stat(pluginDir); err == nil && !fi.IsDir() {
		return pluginDir, fmt.Errorf("file %s already exists and is not a directory", pluginDir)
	}
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return pluginDir, err
	}

	creator := findCreator(lang)
	if creator == nil {
		return pluginDir, fmt.Errorf(`unsupported creating %s template.
You might need to run 'slctl plugin create langs'`, lang)
	}

	files := creator.files(plugin, pluginDir)
	files = append(files, marshal{
		path: filepath.Join(pluginDir, PluginFileName),
		in:   plugin,
	})

	for _, file := range files {
		if err := save(file); err != nil {
			return pluginDir, err
		}
	}
	return pluginDir, nil
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

type compile struct {
	path     string
	in       interface{}
	template string
}

func (u compile) filepath() string {
	return u.path
}

func (u compile) content() ([]byte, error) {
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Parse(u.template))
	if err := parsed.Execute(&buf, u.in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func save(file file) (err error) {
	out, err := file.content()
	if err != nil {
		return
	}
	return ioutil.WriteFile(file.filepath(), out, 0644)
}
