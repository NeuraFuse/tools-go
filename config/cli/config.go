package cli

import (
	"../../env"
	"../../filesystem"
	"../../logging"
	"../../objects"
	"../../readers/yaml"
	"../../runtime"
	"../../terminal"
	"../../users"
	"../../vars"
	"../../objects/strings"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var template *Default
var fileName string = context + format
var filePath string

var cliChecked bool = false

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
	template.APIVersion = vars.NeuraKubeVersion
	template.Kind = strings.Title(context)
}

func (f F) GetConfig() *Default {
	return template
}

func (f F) GetFilePath() string {
	if filePath == "" {
		return "users/" + fileName
	} else {
		return filePath
	}
}

func (f F) SetFilePath(basePath string) {
	filePath = basePath + "/" + fileName
}

func (f F) Configure() {
	selection := terminal.GetUserSelection("What would you like to configure?", []string{"Default user", "Default project", "Updates"}, false, false)
	if selection == "Default user" {
		f.templateDefaultUser()
	} else if selection == "Default project" {
		f.templateDefaultProject()
	} else if selection == "Updates" {
		f.SetUpdates()
	}
}

func (f F) SetDefaultProject(idActive string) bool {
	var set bool
	selectionDefaultProject := terminal.GetUserSelection("Do you want to set "+idActive+" as your default project?", []string{}, false, true)
	if selectionDefaultProject == "Yes" {
		f.setValue("Spec.Projects.DefaultID", idActive)
		set = true
	}
	return set
}

func (f F) HasUserNameDefault() bool {
	exists := true
	if template.Spec.Users.DefaultID == "" {
		exists = false
	}
	return exists
}

func (f F) templateDefaultUser() {
	if f.HasUserNameDefault() {
		logging.Log([]string{"", vars.EmojiUser, vars.EmojiWarning}, "You are about to overwrite your current default user: "+f.getValue("Users.Default"), 0)
	}
	f.setValue("Spec.Users.Default", terminal.GetUserSelection("Which user should be the default?", users.GetAllIDs(), false, false))
}

func (f F) templateDefaultProject() {
	if f.ValidConfigDefaultProject() {
		logging.Log([]string{"", vars.EmojiProject, vars.EmojiWarning}, "You are about to overwrite your current default project: "+f.getValue("Spec.Projects.DefaultID"), 0)
	}
	f.setValue("Spec.Projects.DefaultID", terminal.GetUserSelection("Which project should be the default?", f.GetAllProjectIDs(), false, false))
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
	logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiInfo}, "It is highly recommended to turn on auto. updates because "+envActive+" is still in an alpha development status.", 0)
	sel := terminal.GetUserSelection("Should "+envActive+" automatically update itself?", []string{}, false, true)
	if sel == "Yes" {
		f.setValue("Spec.Updates.Auto.Status", "active")
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiSuccess}, "Checking and applying new updates automatically.\n", 0)
	} else {
		f.setValue("Spec.Updates.Auto.Status", "disabled")
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiSuccess}, "Auto update is now turned off.\n", 0)
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