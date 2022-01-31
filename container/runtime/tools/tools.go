package tools

import (
	"github.com/neurafuse/tools-go/data/providers/crawler"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) DataAggregation(module, project, dataPath string) {
	switch module {
	case "gpt":
		dataPath = dataPath + "input.txt"
		if !filesystem.Exists(dataPath) {
			crawler.F.Router(crawler.F{}, project, dataPath)
		} else {
			logging.Log([]string{"", vars.EmojiData, vars.EmojiSuccess}, "Found existing training data.\n", 0)
		}
	}
}
