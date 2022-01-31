package web

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/nlp/filter"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var lineSep string = "\n"

func (f F) Router(index, contentNew, format string) string {
	contentNew = filter.Code(contentNew, format)
	contentNew = filter.TrimSpaceMultiline(contentNew, lineSep)
	contentNew = filter.UnrealWords(contentNew)
	return contentNew
}

func (f F) Post(contentNew string) string {
	logging.Log([]string{"", vars.EmojiCompression, vars.EmojiInfo}, "Performing post processing filters..", 0)
	logging.ProgressSpinner("start")
	contentNew = filter.EmptyLines(contentNew)
	logging.ProgressSpinner("stop")
	return contentNew
}
