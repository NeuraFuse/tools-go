package user

import (
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	infraID "github.com/neurafuse/tools-go/infrastructures/id"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
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
	config.APIVersion = vars.NeuraKubeAPIVersion
	config.Kind = strings.Title(context)
}

func (f F) GetFilePath() string {
	var filePath string
	filePath = usersID.BasePath + "/" + usersID.F.GetActive(usersID.F{}) + "/" + fileName
	if env.F.CLI(env.F{}) {
		filePath = runtime.F.GetOSInstallDir(runtime.F{}) + env.F.GetActive(env.F{}, false) + "/" + filePath
	}
	return filePath
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(config, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}

func (f F) setValue(key string, value string) {
	yaml.StructToFile(f.GetFilePath(), objects.SetFieldValueFromStruct(config, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) SetDefaults() {
	var id string = terminal.GetUserSelection("What infrastructure should be the default?", infraConfig.F.GetAllIDs(infraConfig.F{}), true, false)
	f.setValue("Spec.Defaults.Infrastructure.ID", id)
	infraID.F.SetActive(infraID.F{}, id)
}

func (f F) SetDefaultInfraID() {
	var id string = f.getValue("Spec.Defaults.Infrastructure.ID")
	infraID.F.SetActive(infraID.F{}, id)
}
