package dependencies

import (
	"../logging"
	"../updater/golang"
	"../vars"
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
