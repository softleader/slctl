package plugin

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config/token"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/strcase"
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
	exec(plugin *Metadata) Commands
	hook(plugin *Metadata) Commands
	files(plugin *Metadata, pluginDir string) []file
}

func Create(lang string, plugin *Metadata, path string) (string, error) {
	if path = strings.TrimSpace(path); path != "" {
		if expanded, err := homedir.Expand(path); err != nil {
			path = expanded
		}
	} else {
		if wd, err := os.Getwd(); err != nil {
			return "", err
		} else {
			path = filepath.Join(wd, plugin.Name)
		}
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return path, err
	}

	if err := paths.EnsureDirectory(logrus.StandardLogger(), path); err != nil {
		return path, err
	}

	creator, found := Creators[lang]
	if !found {
		return path, fmt.Errorf(`unsupported creating %s template.
You might need to run 'slctl plugin create langs'`, lang)
	}
	plugin.Exec = creator.exec(plugin)
	plugin.Hook = creator.hook(plugin)
	plugin.GitHub.Scopes = token.Scopes
	plugin.IgnoreGlobalFlags = false
	files := creator.files(plugin, path)
	files = append(files, marshal{
		path: filepath.Join(path, MetadataFileName),
		in:   plugin,
	})

	for _, file := range files {
		if err := save(file); err != nil {
			return path, err
		}
	}
	return path, nil
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
		"title":      strings.Title,
		"camel":      strcase.ToCamel,
		"lowerCamel": strcase.ToLowerCamel,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
	}
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Funcs(funcMap).Parse(u.template))
	if err := parsed.Execute(&buf, u.in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func save(file file) (err error) {
	paths.EnsureDirectory(logrus.StandardLogger(), filepath.Dir(file.filepath())) // ensure parent exist
	out, err := file.content()
	if err != nil {
		return
	}
	return ioutil.WriteFile(file.filepath(), out, 0755)
}
