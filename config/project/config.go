package project

import (
	"../../env"
	"../../filesystem"
	"../../logging"
	"../../objects"
	"../../objects/strings"
	"../../readers/yaml"
	"../../runtime"
	"../../terminal"
	"../../users"
	"../../vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var config *Default
var fileName string = context + format

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
	return users.ProjectPathActive + fileName
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) SetSpec() {
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "There are missing or invalid project settings.", 0)
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "Please configure them.", 0)
	logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "Please provide the absolute path to your project working directory.", 0)
	logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "E.g. /home/"+runtime.F.GetOSUsername(runtime.F{})+"/projects/test-1", 0)
	workingDir := terminal.GetUserInput("Project working directory path:")
	f.setValue("Spec.WorkingDir", workingDir)
}

func (f F) setValue(key string, value string) {
	yaml.StructToFile(f.GetFilePath(), objects.SetFieldValueFromStruct(config, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}