package tools

import (
	colly "github.com/gocolly/colly/v2"
	"github.com/neurafuse/tools-go/data/processors/nlp/knowledge"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Index(domain, domainPath, filter, savePath string, maxDepth int, domainonly, saveOnlyMerged, overwriteExisting bool) int {
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInspect}, "Starting indexing..", 0)
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInfo}, "Domain: "+domain, 0)
	url := f.getURL(domain + domainPath)
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInfo}, "URL: "+url, 0)
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInfo}, "Max url depth: "+strings.ToString(maxDepth), 0)
	langPref := "en-US"
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInfo}, "Preferred language: "+langPref, 0)
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiDir}, "Saving to path: "+savePath, 0)
	if filter != "" {
		logging.Log([]string{"\n", vars.EmojiCompression, vars.EmojiInfo}, "Active filter: "+filter, 0)
	}
	var c *colly.Collector
	if domainonly {
		c = colly.NewCollector(
			colly.MaxDepth(maxDepth),
			colly.AllowedDomains(domain),
			colly.Async(true),
		)
	} else {
		c = colly.NewCollector(
			colly.MaxDepth(maxDepth),
			colly.Async(true),
		)
	}
	logging.ProgressSpinner("start")
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 25})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", langPref)
		logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Visiting: "+r.URL.String(), 2)
	})
	pageID := 0
	resI := 0
	resIlog := 0
	var index string
	c.OnResponse(func(r *colly.Response) {
		logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Saving url: "+url+"\n", 2)
		if saveOnlyMerged {
			if overwriteExisting && pageID == 0 {
				filesystem.Delete(savePath, false)
			}
			contentNew := strings.BytesToString(r.Body)
			if filter == "knowledge" {
				index = index + knowledge.F.Live(knowledge.F{}, index, contentNew, "html")
			} else {
				index = index + contentNew
			}
		} else {
			filesystem.SaveByteArrayToFile(r.Body, savePath+strings.ToString(pageID)+"/index.html")
		}
		if resIlog == maxDepth*22 {
			logging.ProgressSpinner("stop")
			logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Already processed urls: "+strings.ToString(resI), 0)
			//logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Recent index size: "+strings.FloatToString(filesystem.GetSize(savePath, "gb"))+" GB", 0)
			logging.ProgressSpinner("start")
			resIlog = 0
		}
		pageID++
		resI++
		resIlog++
	})
	c.OnError(func(r *colly.Response, err error) {
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to process request: "+r.Request.URL.String()+"!", false, false, false)
	})
	c.Visit(url)
	c.Wait()
	filesystem.AppendStringToFile(savePath, index) // WIP: objects.CallStructInterfaceFuncByName(Packages{}, filter, "GPT", strings.BytesToString(r.Body))
	logging.ProgressSpinner("stop")
	content := filesystem.FileToString(savePath)
	filesystem.Delete(savePath, false)
	var maxLines int
	maxLines = strings.LinesCount(content)
	if filter == "knowledge" {
		logging.Log([]string{"", vars.EmojiCompression, vars.EmojiInfo}, "Line count after live filter: "+strings.ToString(maxLines)+"\n", 0)
		content, maxLines = knowledge.F.PostProcess(knowledge.F{}, content)
	}
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Final line count: "+strings.ToString(maxLines), 0)
	filesystem.AppendStringToFile(savePath, content)
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Processed "+strings.ToString(resI)+" urls in total.", 0)
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Total index size: "+strings.FloatToString(filesystem.GetSize(savePath, "gb"))+" GB", 0)
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiSuccess}, "Finished crawling.\n", 0)
	return maxLines
}

func (f F) getURL(domain string) string {
	var https string = "https://"
	/*var www string = "www."
	if !strings.Contains(domain, www) {
		domain = www + domain
	}*/
	if !strings.Contains(domain, https) {
		domain = https + domain
	}
	return domain
}
