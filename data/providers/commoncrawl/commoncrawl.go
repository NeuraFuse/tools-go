package commoncrawl

import (
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(action string) {
	config := GetConfiguration()
	logging.Log([]string{"", vars.EmojiData, vars.EmojiDir}, " Creating folder: "+config.DataFolder, 0)
	filesystem.CreateDir(config.DataFolder, false)
	logging.Log([]string{"", vars.EmojiData, vars.EmojiDir}, "Creating folder: "+config.MatchFolder, 0)
	filesystem.CreateDir(config.MatchFolder, false)
	logging.Log([]string{"", vars.EmojiData, vars.EmojiInspect}, "Starting scan...", 0)
	logging.ProgressSpinner("start")
	if action == "scan" {
		Scan(config, "golang")
	}
	logging.ProgressSpinner("stop")
	//cc.download()
	//cc.extract()
}
