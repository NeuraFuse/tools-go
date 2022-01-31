package vscode

import (
	"github.com/neurafuse/tools-go/apps/ide/vscode/templates"
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/json"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) CreateConfig(remoteURL, port string) {
	logging.Log([]string{"", vars.EmojiSettings, vars.EmojiSuccess}, "Creating VSCode remote debugging launch config..", 0)
	var launchJSONPath string = f.getProjectConfigPath() + "launch.json"
	var launchJSONBackupPath string = f.getProjectConfigPath() + "launch_backup.json"
	if filesystem.Exists(launchJSONPath) {
		filesystem.RenameFile(launchJSONPath, launchJSONBackupPath)
		filesystem.Delete(launchJSONPath, false)
	}
	configInt := f.getConfigInterface(remoteURL, port)
	json.StructToFile(launchJSONPath, configInt)
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInfo}, "You can now start debugging your remote environment by pressing F5.", 0)
	logging.Log([]string{"", vars.EmojiInfo, vars.EmojiInspect}, "To view the live logs switch in VSCode from the perspective TERMINAL to DEBUG CONSOLE.\n", 0)
}

func (f F) getConfigInterface(remoteURL, port string) interface{} {
	var configInt interface{} = json.FileToStruct(f.getConfigTemplatePath(), &templates.Default{})
	configInt.(*templates.Default).Configurations[0].Host = remoteURL
	configInt.(*templates.Default).Configurations[0].Port = strings.ToInt(port)
	//mlEnv := config.Setting("get", "infrastructure", "Spec.Remote.Environment.Framework", "")
	var localIDERoot string = config.Setting("get", "project", "Spec.Containers.Sync.PathMappings.LocalIDERoot", "")
	configInt.(*templates.Default).Configurations[0].PathMappings[0].LocalRoot = localIDERoot
	var localAppRoot string = config.Setting("get", "project", "Spec.Containers.Sync.PathMappings.LocalAppRoot", "")
	configInt.(*templates.Default).Configurations[0].PathMappings[0].RemoteRoot = localAppRoot
	return configInt
}

func (f F) getProjectConfigPath() string {
	var projectConfigPath string = filesystem.GetWorkingDir(false)
	if config.DevConfigActive() {
		projectConfigPath = projectConfigPath + "../"
	}
	projectConfigPath = projectConfigPath + ".vscode/"
	return projectConfigPath
}

func (f F) getConfigTemplatePath() string {
	var configTemplatePath string = "../tools-go@" + vars.ToolsGoVersion + "/apps/ide/vscode/templates/default.json"
	return configTemplatePath
}
