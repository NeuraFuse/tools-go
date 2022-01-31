package users

import (
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var quest string = "What is your intention for the resource users?"
	var new string = "Create a new one"
	var sw string = "Switch to an existing one"
	var selOpts []string = []string{new, sw}
	var sel string = terminal.GetUserSelection(quest, selOpts, false, false)
	var update bool
	switch sel {
	case new:
		id.F.CreateNew(id.F{})
		update = true
	case sw:
		update = f.chooseExisting()
	}
	if update {
		f.setDefault()
	}
}

func (f F) chooseExisting() bool {
	var update bool
	if !f.multipleExisting() {
		logging.Log([]string{"\n", vars.EmojiAssistant, vars.EmojiInfo}, id.F.GetActive(id.F{})+", there are no other users created yet.", 0)
		var sel string = terminal.GetUserSelection("Do you want to create an additional one?", []string{}, false, true)
		if sel == "Yes" {
			id.F.CreateNew(id.F{})
			update = true
		}
	} else {
		var quest string = "To which user do you want to switch?"
		var opts []string = f.GetAllIDs()
		var userID string = terminal.GetUserSelection(quest, opts, false, false)
		id.F.SetActive(id.F{}, userID)
		update = true
	}
	return update
}

func (f F) multipleExisting() bool {
	var status bool
	var allIds []string = f.GetAllIDs()
	if len(allIds) > 1 {
		status = true
	}
	return status
}

func (f F) GetAllIDs() []string {
	var path string
	if env.F.CLI(env.F{}) {
		path = f.GetLocalUserBasePath()
	} else if env.F.API(env.F{}) {
		path = f.GetAPIUserBasePath()
	}
	return filesystem.Explorer("files", path, []string{}, []string{"hidden", ".yaml"})
}

func (f F) setDefault() {
	var defaultIDrecent string = config.Setting("get", "cli", "Spec.Users.DefaultID", "")
	var idActive string = id.F.GetActive(id.F{})
	if defaultIDrecent != idActive {
		var sel string = terminal.GetUserSelection("Do you want to set "+idActive+" as your default user?", []string{}, false, true)
		if sel == "Yes" {
			config.Setting("set", "cli", "Spec.Users.DefaultID", idActive)
		}
	} else {
		logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiInfo}, idActive+" is the default user.\n", 0)
	}
}

func (f F) Create(userID string) {
	var userPath string = f.getUserPath(userID)
	if !filesystem.Exists(userPath) {
		filesystem.CreateDir(userPath, false)
	}
}

func (f F) getUserPath(userID string) string {
	var userPath string
	if env.F.CLI(env.F{}) {
		userPath = f.getUserBasePath() + userID
	} else if env.F.API(env.F{}) {
		userPath = f.GetAPIUserBasePath() + userID
	}
	return userPath
}

func (f F) getUserBasePath() string {
	var basePath string
	if env.F.CLI(env.F{}) {
		basePath = f.GetLocalUserBasePath()
	} else if env.F.API(env.F{}) {
		basePath = f.GetAPIUserBasePath()
	}
	return basePath
}

func (f F) GetLocalUserBasePath() string {
	var basePath string
	basePath = runtime.F.GetOSInstallDir(runtime.F{}) + env.F.GetActive(env.F{}, false) + "/" + id.BasePath // TODO: Ref.
	return basePath
}

func (f F) GetAPIUserPath() string {
	var path string
	path = f.GetAPIUserBasePath() + id.F.GetActive(id.F{})
	return path
}

func (f F) GetAPIUserBasePath() string {
	var basePath string
	if env.F.Container(env.F{}) {
		basePath = id.BasePath
	} else {
		basePath = filesystem.GetWorkingDir(false) + id.BasePath
	}
	return basePath
}

func (f F) Exists(id string) bool {
	return strings.ArrayContains(f.GetAllIDs(), id)
}

func (f F) Existing() bool {
	var basePath string = f.getUserBasePath()
	if !filesystem.Exists(basePath) {
		filesystem.CreateDir(basePath, false)
		return false
	} else {
		if len(filesystem.Explorer("files", basePath, []string{}, []string{"lost+found"})) == 0 {
			return false
		} else {
			return true
		}
	}
}
