// +build !wasm

package app

import (
	"net/http"
	"strings"
)

// ResourceProvider is the interface that describes a provider for resources.
//
// App resources are the mandatory resources required to run a PWA. They are
// generated by the Handler and are accessible from the root path. Eg:
//  "/app-worker.js"
//  "/manifest.webmanifest"
//  "/wasm_exec.js"
//
// Static resources are the resources used by the PWA such as the web assembly
// binary, styles, scripts, or images. They can be located on a localhost or a
// remote bucket. In order to avoid confusion with PWA required resources,
// static resources URL paths are always prefixed by "/web". Eg:
//  "/web/app.wasm"
//  "/web/main.css"
//  "/web/background.jpg"
//
// If the resource provider is an http.Handler, the handler is used to serve
// static resources requests.
type ResourceProvider interface {
	// The path to the root directory where app resources are accessible.
	AppResources() string

	// The path or URL where the web directory that contains static resources is
	// located.
	StaticResources() string

	// The URL of the app.wasm file. This must match the pattern:
	//  StaticResources/web/WASM_FILE.
	AppWASM() string

	// The URL of the robots.txt file. This must match the pattern:
	//  StaticResources/web/robots.txt.
	RobotsTxt() string

	// The URL of the ads.txt file. This must match the pattern:
	//  StaticResources/web/ads.txt.
	AdsTxt() string
}

// LocalDir returns a resource provider that serves static resources from a
// local directory located at the given path.
func LocalDir(path string) ResourceProvider {
	return localDir{
		Handler: http.StripPrefix("/web/", http.FileServer(http.Dir(path))),
		path:    path,
	}
}

type localDir struct {
	http.Handler
	path string
}

func (d localDir) AppResources() string {
	return ""
}

func (d localDir) StaticResources() string {
	return ""
}

func (d localDir) AppWASM() string {
	return "/web/app.wasm"
}

func (d localDir) RobotsTxt() string {
	return "/web/robots.txt"
}

func (d localDir) AdsTxt() string {
	return "/web/ads.txt"
}

// RemoteBucket returns a resource provider that provides resources from a
// remote bucket such as Amazon S3 or Google Cloud Storage.
func RemoteBucket(url string) ResourceProvider {
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, "/web")

	return remoteBucket{
		url: url,
	}
}

type remoteBucket struct {
	url string
}

func (b remoteBucket) AppResources() string {
	return ""
}

func (b remoteBucket) StaticResources() string {
	return b.url
}

func (b remoteBucket) AppWASM() string {
	return b.StaticResources() + "/web/app.wasm"
}

func (b remoteBucket) RobotsTxt() string {
	return b.StaticResources() + "/web/robots.txt"
}

func (b remoteBucket) AdsTxt() string {
	return b.StaticResources() + "/web/ads.txt"
}

// GitHubPages returns a resource provider that provides resources from GitHub
// pages. This provider must only be used to generate static websites with the
// GenerateStaticWebsite function.
func GitHubPages(repoName string) ResourceProvider {
	if !strings.HasPrefix(repoName, "/") {
		repoName = "/" + repoName
	}

	return gitHubPages{repo: repoName}
}

type gitHubPages struct {
	repo string
}

func (g gitHubPages) AppResources() string {
	return g.repo
}

func (g gitHubPages) StaticResources() string {
	return g.repo
}

func (g gitHubPages) AppWASM() string {
	return g.StaticResources() + "/web/app.wasm"
}

func (g gitHubPages) RobotsTxt() string {
	return g.StaticResources() + "/web/robots.txt"
}

func (g gitHubPages) AdsTxt() string {
	return g.StaticResources() + "/web/ads.txt"
}
