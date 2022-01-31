package infrastructures

import (
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var infraIDDefault string = config.Setting("get", "user", "Spec.Defaults.Infrastructure.ID", "")
	var optSwitch string = "Switch default infrastructure (" + infraIDDefault + ")"
	var opts []string = []string{optSwitch}
	var quest string = "What is your intention?"
	var sel string = terminal.GetUserSelection(quest, opts, false, false)
	switch sel {
	case optSwitch:
		f.sw()
	}
}

func (f F) sw() {
	var ids []string = infraConfig.F.GetAllIDs(infraConfig.F{})
	if len(ids) > 1 {
		var id string = terminal.GetUserSelection("Which infrastructure should be the new default one?", ids, false, false)
		config.Setting("set", "user", "Spec.Defaults.Infrastructure.ID", id)
	} else {
		logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiInfo}, "The selected default is the only existing infrastructure.\n", 0)
	}
}
