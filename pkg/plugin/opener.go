package plugin

import (
	"github.com/skratchdot/open-golang/open"
)

// Open 開啟 plugin 的目錄或遠端網址
func (p *Plugin) Open(app string) error {
	src := src(p)
	if len(app) != 0 {
		return open.StartWith(src, app)
	}
	return open.Start(src)
}

// OpenAndWait 開啟 plugin 的目錄或遠端網址, 並等待 Open 命令結束
func (p *Plugin) OpenAndWait(app string) error {
	src := src(p)
	if len(app) != 0 {
		return open.RunWith(src, app)
	}
	return open.Run(src)
}

func src(p *Plugin) (source string) {
	if p.FromGitHub() {
		return "https://" + p.Source
	}
	return p.Dir
}
