package project

import (
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	projectsID "github.com/neurafuse/tools-go/projects/id"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/logging"
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
	var filePath string = f.GetFilePath()
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
	if env.F.CLI(env.F{}) {
		filePath = filesystem.GetWorkingDir(false) + "." + vars.NeuraCLINameID + "/" + fileName
	} else if env.F.API(env.F{}) {
		var preemble string
		if env.F.Container(env.F{}) {
			preemble = infraConfig.F.GetContainerUserPath(infraConfig.F{})
		} else {
			preemble = filesystem.GetWorkingDir(false) + usersID.BasePath
		}
		filePath = preemble + usersID.F.GetActive(usersID.F{}) + "/projects/" + projectsID.GetActive() + "/" + fileName
	}
	return filePath
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) setValue(key string, value string) {
	yaml.StructToFile(f.GetFilePath(), objects.SetFieldValueFromStruct(config, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(config, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}

func (f F) GetRemoteURL(appID string) string {
	var remoteURL string
	var remoteURLConfigKey string = "Spec.Containers.Sync.RemoteURL"
	remoteURL = f.getValue(remoteURLConfigKey)
	if remoteURL == "" {
		logging.Log([]string{"", vars.EmojiRocket, vars.EmojiWaiting}, "Please provide information about the routing to the app "+appID+".", 0)
		remoteURL = terminal.GetUserInput("What is the remote URL (or IP) on which the app is reachable via network?")
		f.setValue(remoteURLConfigKey, remoteURL)
	}
	return remoteURL
}

func (f F) NetworkMode(mode string) bool {
	var active bool
	var configKey string = "Spec.Containers.Network.Mode"
	var modeActive string = f.getValue(configKey)
	if modeActive == "" {
		logging.Log([]string{"", vars.EmojiClient, vars.EmojiKubernetes}, "Please provide information about how "+vars.NeuraCLIName+" should connect to the container network.", 0)
		var networkModes []string = []string{"port-forward", "remote-url"}
		var sel string = terminal.GetUserSelection("Please select a container network mode:", networkModes, false, false)
		f.setValue(configKey, sel)
		modeActive = sel
	}
	if modeActive == mode {
		active = true
	}
	return active
}