package assistant

import (
	dev ".."
	"../../../../tools-go/logging"
	"../../../../tools-go/objects"
	"../../../../tools-go/readers/yaml"
	"../../../../tools-go/runtime"
	"../../../../tools-go/terminal"
	"../../../../tools-go/vars"
	"../../../../tools-go/env"
)

func Create() {
	envActive := env.F.GetActive(env.F{}, true)
	var opts []string = []string{"As a tool for my own developments", "I also want to develop "+envActive+" itself"}
	sel := terminal.GetUserSelection("How do you want to use "+envActive+"?", opts, false, false)
	if sel == opts[1] {
		setStatusActive()
		setLogLevel()
		setAPI()
		setDocker()
	}
}

func setStatusActive() {
	sel := terminal.GetUserSelection("Do you want to just create the devconfig for later use or do you want to activate it directly after configuration?", []string{"Default: active", "Disabled"}, false, false)
	setValue("Status", sel)
}

func setLogLevel() {
	logging.Log([]string{"\n", vars.EmojiDev, vars.EmojiInfo}, "Please choose your desired log level. You can choose between user level (default blank), info and debug.", 0)
	sel := terminal.GetUserSelection("Which log level do you want to choose?", []string{"Default: blank", "info", "debug"}, false, false)
	if sel != "" {
		setValue("LogLevel", sel)
	}
}

func setAPI() {
	logging.Log([]string{"\n", vars.EmojiDev, vars.EmojiInfo}, "You can deploy "+vars.NeuraKubeName+" locally or within a cluster.", 0)
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiSpeed}, "Sometimes it is useful to develop locally for faster prototyping.", 0)
	sel := terminal.GetUserSelection("Do you want "+vars.NeuraCLINameRepo+" to connect to a locally deployed "+vars.NeuraKubeName+"?", []string{"Default: cluster", "localhost"}, false, false)
	if sel != "" {
		setValue("API.Address", sel)
	}
}

func setDocker() {
	logging.Log([]string{"\n", vars.EmojiDev, vars.EmojiInfo}, "Please choose a default docker repository. Choose default blank to just use the official prebuilt "+vars.OrganizationName+" docker images.", 0)
	sel := terminal.GetUserSelection("What should be the default docker repository (address)?", []string{"Default: blank", "gcr.io/djw-ai/services"}, false, false)
	if sel != "" {
		setValue("Docker.RepoAddress", sel)
		sel := terminal.GetUserInput("What is the account username of your custom docker repository?")
		if sel != "" {
			setValue("Docker.User.Name", sel)
		}
		sel = terminal.GetUserInput("What is the account password of your custom docker repository?")
		if sel != "" {
			setValue("Docker.User.Password", sel)
		}
	}
}

func setValue(key string, value string) {
	yaml.StructToFile(dev.F.GetFilePath(dev.F{}), objects.SetFieldValueFromStruct(dev.F.GetConfig(dev.F{}), key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func getValue(key string) string {
	return objects.GetFieldValueFromStruct(dev.F.GetConfig(dev.F{}), key, runtime.F.GetCallerInfo(runtime.F{}, true))
}
