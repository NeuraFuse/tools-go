package projects

import (
	"strings"

	"github.com/neurafuse/tools-go/config"
	cliconfig "github.com/neurafuse/tools-go/config/cli"
	devconfigAssistant "github.com/neurafuse/tools-go/config/dev/assistant"
	projectConfig "github.com/neurafuse/tools-go/config/project"
	"github.com/neurafuse/tools-go/crypto/jwt"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/users"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var create bool
	if !f.exists() {
		create = true
	}
	var quest string = "What is your intention for the resource projects?"
	var new string = "Create a new project at this location"
	var selOpts []string = []string{new}
	var sel string = terminal.GetUserSelection(quest, selOpts, false, false)
	switch sel {
	case new:
		create = true
	}
	if create {
		f.createNew()
	}
}

func (f F) CheckConfigs() {
	if !f.exists() {
		f.createNew()
	}
	f.loadExisting()
}

func (f F) exists() bool {
	return filesystem.Exists(vars.ProjectsBasePath + "project.yaml") // TODO: Ref.
}

func (f F) loadExisting() {
	f.prepareSettings()
	f.CheckAuth()
	config.Setting("init", "dev", "Spec.", "")
}

func (f F) prepareSettings() {
	if !config.ValidSettings("cli", "updates", false) {
		cliconfig.F.SetUpdates(cliconfig.F{})
	}
}

func (f F) CheckAuth() {
	if !config.ValidSettings("user", "Spec.login", false) {
		config.Setting("set", "user", "Spec.Auth.JWT.SigningKey", jwt.GenerateSigningKey())
	}
}

func (f F) GetAllIDs() []string {
	return filesystem.Explorer("files", vars.ProjectsBasePath, []string{}, []string{"hidden", ".yaml", "infrastructure"})
}

func (f F) createNew() {
	vars.ProjectDirActive = vars.ProjectsBasePath // TODO: Ref.
	var configFilePath string = projectConfig.F.GetFilePath(projectConfig.F{})
	if filesystem.Exists(configFilePath) {
		logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiInfo}, "There is aleady a project created at this working directory.", 0)
		var sel string = terminal.GetUserSelection("Do you want to overwrite it (delete config & create new)?", []string{}, false, true)
		if sel == "Yes" {
			filesystem.Delete(configFilePath, false)
		} else {
			f.userInfoConfigEdit(configFilePath)
			terminal.Exit(0, "")
		}
	}
	var wdDirName string = filesystem.GetLastFolderFromPath(filesystem.GetWorkingDir(false))
	var id string = wdDirName
	if config.DevConfigActive() {
		id = terminal.GetUserSelection("Choose a name for your new project (has to be folder name of project)", []string{wdDirName}, true, false)
	}
	//vars.ProjectPathYaml = vars.ProjectDirActive + "/config.yaml"
	filesystem.CreateDir(vars.ProjectDirActive, false)
	config.Setting("init", "project", "", "")
	config.Setting("set", "project", "Metadata.ID", id)
	config.Setting("set", "project", "Metadata.Name", id)
	f.userInfoConfigEdit(configFilePath)
	cliconfig.F.SetUpdates(cliconfig.F{})
	devconfigAssistant.Create()
}

func (f F) userInfoConfigEdit(configFilePath string) {
	logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiInfo}, "All your settings are saved as .yaml files on different locations.", 0)
	logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiProject}, "Project related settings are located at the current working directory:\n"+configFilePath+"\n", 0)
	var localUserBasePath string = users.F.GetLocalUserBasePath(users.F{})
	var logMsg string = "Cli, user, developer and infrastructure related settings are located on a per user basis at the OS install directory:\n" + localUserBasePath + "\n"
	logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiUser}, logMsg, 0)
	logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiEdit}, "After the creation you can easy edit them with any editor or via the "+vars.NeuraCLIName+" assistant.", 0)
}

func (f F) GetContainerExternalExecPath(project string, codeFormat string) string {
	var fileName string = strings.TrimSuffix(project, "-"+codeFormat)
	var dir string = f.GetContainerProjectBasePath(project, codeFormat) + fileName + "." + codeFormat
	return dir
}

func (f F) GetContainerProjectBasePath(project string, codeFormat string) string {
	var dir string = project + "/"
	if !env.F.Container(env.F{}) {
		dir = "../" + dir
	}
	return dir
}

func (f F) GetContainerExternalDataPath(project string, codeFormat string) string {
	return f.GetContainerProjectBasePath(project, codeFormat) + "data/training/input/"
}
