package dependencies

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/updater/golang"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) CheckBuild() bool {
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiProcess}, "Checking build dependencies..", 0)
	var updated bool
	if golang.F.Check(golang.F{}) {
		updated = true
	}
	return updated
}
