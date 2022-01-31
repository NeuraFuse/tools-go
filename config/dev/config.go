package dev

import (
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/logging/emoji"
	"github.com/neurafuse/tools-go/terminal"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var config *Default
var filePath string
var FileName string = context + format

func (f F) SetConfig() (*Default, string) {
	f.Exists()
	return f.GetConfig(), f.GetFilePath()
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) IsLogLevelActive(mode string) bool {
	var status bool
	if f.GetConfig().Spec.LogLevel == mode {
		status = true
	}
	return status
}

func (f F) GetFilePath() string {
	var filePath string
	filePath = usersID.BasePath + "/" + usersID.F.GetActive(usersID.F{}) + "/" + FileName
	if env.F.CLI(env.F{}) {
		filePath = runtime.F.GetOSInstallDir(runtime.F{}) + env.F.GetActive(env.F{}, false) + "/" + filePath
	}
	return filePath
}

func (f F) Exists() bool {
	f.setTemplate()
	filePath := f.GetFilePath()
	var exists bool
	if filesystem.Exists(filePath) {
		yaml.FileToStruct(filePath, &config)
		exists = true
	}
	return exists
}

func (f F) setTemplate() {
	config = &Default{}
	config.APIVersion = vars.NeuraCLIAPIVersion
	config.Kind = strings.Title(context)
	config.Spec.Build.Neurakube.Version = vars.NeuraKubeVersion
	config.Spec.Build.Neuracli.Version = vars.NeuraCLIVersion
}

func (f F) setValue(key string, value string) {
	yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(config, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}

func (f F) RequireAPIAddrsNonLocal(context string) {
	var configKey string = "Spec.API.Address"
	if f.getValue(configKey) == "local" {
		emoji.Println("", vars.EmojiAPI, vars.EmojiWaiting, "The "+context+" module is not available if you have activated the API localhost mode.")
		var sel string = terminal.GetUserSelection("Do you want to switch to API cluster mode?", []string{}, false, true)
		if sel == "Yes" {
			f.setValue(configKey, "cluster")
		} else {
			terminal.Exit(0, "")
		}
	}
}