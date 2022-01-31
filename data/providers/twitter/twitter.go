package twitter

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(module, dataPath string) {
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInspect}, "Starting twitter API..", 0)
}
