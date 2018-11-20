package getter

import (
	"bytes"
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
)

// Getter is an interface to support GET to the specified URL.
type Getter interface {
	//Get file content by url string
	Get(url string) (*bytes.Buffer, error)
}

// Constructor is the function for every getter which creates a specific instance
// according to the configuration
type Constructor func(URL, CertFile, KeyFile, CAFile string) (Getter, error)

// Provider represents any getter and the schemes that it supports.
//
// For example, an HTTP provider may provide one getter that handles both
// 'http' and 'https' schemes.
type Provider struct {
	Schemes []string
	New     Constructor
}

// Provides returns true if the given scheme is supported by this Provider.
func (p Provider) Provides(scheme string) bool {
	for _, i := range p.Schemes {
		if i == scheme {
			return true
		}
	}
	return false
}

// Providers is a collection of Provider objects.
type Providers []Provider

// ByScheme returns a Provider that handles the given scheme.
//
// If no provider handles this scheme, this will return an error.
func (p Providers) ByScheme(scheme string) (Constructor, error) {
	for _, pp := range p {
		if pp.Provides(scheme) {
			return pp.New, nil
		}
	}
	return nil, fmt.Errorf("scheme %q not supported", scheme)
}

// All finds all of the registered getters as a list of Provider instances.
// Currently the build-in http/https getter and the discovered
// plugins with downloader notations are collected.
func All(settings environment.EnvSettings) Providers {
	result := Providers{
		{
			Schemes: []string{"http", "https"},
			New:     newHTTPGetter,
		},
	}
	pluginDownloaders, _ := collectPlugins(settings)
	result = append(result, pluginDownloaders...)
	return result
}

// ByScheme returns a getter for the given scheme.
//
// If the scheme is not supported, this will return an error.
func ByScheme(scheme string, settings environment.EnvSettings) (Provider, error) {
	// Q: What do you call a scheme string who's the boss?
	// A: Bruce Schemestring, of course.
	a := All(settings)
	for _, p := range a {
		if p.Provides(scheme) {
			return p, nil
		}
	}
	return Provider{}, fmt.Errorf("scheme %q not supported", scheme)
}
