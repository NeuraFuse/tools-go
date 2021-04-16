package server

import (
	"../../env"
	"../../filesystem"
	"../../readers/yaml"
	"../../runtime"
	"../../users"
	"../../vars"
	"../../objects/strings"
)

type F struct{}
var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var config *Default
var fileName string = context+format

func (f F) SetConfig() (*Default, string) {
	f.exists()
	return f.GetConfig(), f.GetFilePath()
}

func (f F) exists() {
	f.setTemplate()
	filePath := f.GetFilePath()
	if filesystem.Exists(filePath) {
		yaml.FileToStruct(filePath, &config)
	}
}

func (f F) setTemplate() {
	config = &Default{}
	config.APIVersion = vars.NeuraKubeVersion
	config.Kind = strings.Title(context)
}

func (f F) GetFilePath() string {
	return users.BasePath + "/" + fileName
}

func (f F) GetConfig() *Default {
	return config
}