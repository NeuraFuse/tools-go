package data

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	pack := cliArgs[2]
	input := cliArgs[3]
	logging.Log([]string{"\n", vars.EmojiData, vars.EmojiInfo}, "Context "+input+" with package "+pack+"..", 0)
	switch pack {
	case "cc":
		pack = "commoncrawl"
	}
	objects.CallStructInterfaceFuncByName(Packages{}, strings.Title(pack), "Router", input)
}
