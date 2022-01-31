package cli

import (
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging/emoji"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var template *Default
var fileName string = context + format
var filePath string

var cliChecked bool

func (f F) SetConfig() (*Default, string) {
	f.exists()
	if !cliChecked {
		cliChecked = true
	}
	return f.GetConfig(), f.GetFilePath()
}

func (f F) exists() {
	f.setTemplate()
	if filesystem.Exists(f.GetFilePath()) {
		yaml.FileToStruct(f.GetFilePath(), &template)
	}
}

func (f F) setTemplate() {
	template = &Default{}
	template.APIVersion = vars.NeuraCLIAPIVersion
	template.Kind = strings.Title(context)
}

func (f F) GetConfig() *Default {
	return template
}

func (f F) GetFilePath() string {
	var filePath string
	filePath = "users/" + fileName
	if env.F.CLI(env.F{}) {
		filePath = runtime.F.GetOSInstallDir(runtime.F{}) + env.F.GetActive(env.F{}, false) + "/" + filePath
	}
	return filePath
}

func (f F) Configure() {
	var opts []string = []string{"Automatic updates"}
	var selection string = terminal.GetUserSelection("What would you like to configure?", opts, false, false)
	if selection == opts[0] {
		f.SetUpdates()
	}
}

func (f F) SetDefault(resourceType string, idActive string) bool {
	var set bool
	var selectionDefaultProject string = terminal.GetUserSelection("Do you want to set "+idActive+" as your default "+resourceType+"?", []string{}, false, true)
	if selectionDefaultProject == "Yes" {
		f.setValue("Spec."+strings.Title(resourceType)+"s.DefaultID", idActive)
		set = true
	}
	return set
}

func (f F) HasUserNameDefault() bool {
	var exists bool = true
	if template.Spec.Users.DefaultID == "" {
		exists = false
	}
	return exists
}

func (f F) GetAllProjectIDs() []string {
	ids := filesystem.Explorer("files", vars.ProjectsBasePath, []string{}, []string{"hidden", ".yaml"})
	return ids
}

func (f F) ValidConfigDefaultProject() bool {
	var valid bool
	fieldsPrefix := "Projects."
	fields := []string{"Default"}
	valid = objects.StructFieldValuesExisting(template, fieldsPrefix, fields, runtime.F.GetCallerInfo(runtime.F{}, true))
	return valid
}

func (f F) SetUpdates() {
	var envActive string = env.F.GetActive(env.F{}, true)
	emoji.Println("", vars.EmojiAssistant, vars.EmojiWarning, envActive+" is still in an alpha development status.")
	emoji.Println("", vars.EmojiAssistant, vars.EmojiInfo, "It is highly recommended to automatically install updates.")
	var sel string = terminal.GetUserSelection("Should "+envActive+" automatically update itself?", []string{}, false, true)
	if sel == "Yes" {
		f.setValue("Spec.Updates.Auto.Status", "active")
		emoji.Println("", vars.EmojiSettings, vars.EmojiSuccess, "Checking for and applying new updates automatically.\n")
	} else {
		f.setValue("Spec.Updates.Auto.Status", "disabled")
		emoji.Println("", vars.EmojiSettings, vars.EmojiSuccess, "Auto update is now turned off.\n")
	}
}

func (f F) reset(key string) {
	yaml.StructToFile(f.GetFilePath(), objects.SetFieldValueFromStruct(template, key, "", runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) setValue(key, value string) {
	yaml.StructToFile(f.GetFilePath(), objects.SetFieldValueFromStruct(template, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(template, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}
