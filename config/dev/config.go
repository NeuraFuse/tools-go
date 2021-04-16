package dev

import (
	"../../env"
	"../../filesystem"
	"../../objects"
	"../../objects/strings"
	"../../readers/yaml"
	"../../runtime"
	"../../users"
	"../../vars"
)

type F struct{}
var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var config *Default
var filePath string
var FileName string = context+format

func (f F) SetConfig() (*Default, string) {
	f.Exists()
	return f.GetConfig(), f.GetFilePath()
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) GetFilePath() string {
	if filePath == "" {
		return users.BasePath + "/" + users.GetIDActive() + "/" + FileName
	} else {
		return filePath
	}
}

func (f F) SetFilePath(basePath, userID string) {
	filePath = basePath + "/" + userID + "/" + FileName
}

func (f F) Exists() bool {
	f.setTemplate()
	filePath := f.GetFilePath()
	var exists bool = false
	if filesystem.Exists(filePath) {
		yaml.FileToStruct(filePath, &config)
		exists = true
	}
	return exists
}

func (f F) setTemplate() {
	config = &Default{}
	config.APIVersion = vars.NeuraKubeVersion
	config.Kind = strings.Title(context)
	config.Spec.Build.Neurakube.Version = vars.NeuraKubeVersion
	config.Spec.Build.Neuracli.Version = vars.NeuraCLIVersion
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(config, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}