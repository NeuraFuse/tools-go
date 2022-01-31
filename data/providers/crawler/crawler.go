package crawler

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/vars"

	//"github.com/neurafuse/tools-go/errors"
	//"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/data/providers/crawler/tools"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/io"
	"github.com/neurafuse/tools-go/objects/strings"
)

type F struct{}

func (f F) Router(module, dataPath string) {
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInspect}, "Starting crawler..", 0)
	var maxDepth int
	if env.F.Container(env.F{}) {
		maxDepth = 2
	} else {
		maxDepth = 2
	}
	domain, domainPath, domainonly := f.getCrawlerSettings()
	if f.crawler(domainPath) {
		f.indexTarget(domain, domainPath, dataPath, maxDepth, domainonly)
	} else {
		f.download(domain, domainPath, dataPath)
	}
}

func (f F) getCrawlerSettings() (string, string, bool) {
	domain, domainPath := f.getTarget("bing")
	var domainonly bool
	return domain, domainPath, domainonly
}

func (f F) getTarget(id string) (string, string) {
	var domain string = id + ".com"
	var domainPath string
	switch id {
	case "google":
		domainPath = "/search?q=golang&lr=lang_en"
	case "duckduckgo":
		domainPath = "/?q=golang&ia=web"
	case "bing":
		domainPath = "/?q=golang"
	case "distill":
		domain = id + ".pub"
	case "tinyshakespeare":
		domain = "raw.githubusercontent.com"
		domainPath = "/karpathy/char-rnn/master/data/tinyshakespeare/input.txt"
	}
	return domain, domainPath
}

func (f F) crawler(domainPath string) bool {
	var index bool = true
	if strings.Contains(domainPath, ".txt") {
		index = false
	}
	return index
}

func (f F) indexTarget(domain, domainPath, dataPath string, maxDepth int, domainonly bool) int {
	var maxLinesAvailable int
	maxLinesAvailable = tools.F.Index(tools.F{}, domain, domainPath, "knowledge", dataPath, maxDepth, domainonly, true, true)
	return maxLinesAvailable
}

func (f F) download(domain, domainPath, dataPath string) {
	io.F.DownloadFile(io.F{}, dataPath, domain+domainPath)
}
