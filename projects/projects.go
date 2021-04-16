package projects

import (
	"../config"
	cliconfig "../config/cli"
	devconfigAssistant "../config/dev/assistant"
	"../crypto/jwt"
	"../env"
	"../filesystem"
	"../logging"
	"../objects/strings"
	"../runtime"
	"../terminal"
	"../users"
	"../vars"
)

type F struct{}

var IDActive string = ""

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var new string = "Create a new project"
	var sw string = "Switch to an existing project"
	sel := terminal.GetUserSelection("What do you want to do?", []string{new, sw}, false, false)
	switch sel {
	case new:
		f.createNew(false)
	case sw:
		IDActive = f.choose()
		config.Setting("set", "cli", "Spec.Projects.ActiveID", IDActive)
		cliconfig.F.SetDefaultProject(cliconfig.F{}, IDActive)
	}
}

func (f F) CheckConfigs() {
	if !f.Existing() {
		f.createNew(true)
	}
	f.loadExisting()
}

func (f F) Exists(id string) bool {
	return strings.ArrayContains(f.GetAllIDs(), id)
}

func (f F) Existing() bool {
	return !filesystem.DirIsEmpty(vars.ProjectsBasePath)
}

func (f F) setPaths() { // TODO: Refactor users.GetIDActive()
	users.ProjectPathActive = users.BasePath + "/" + users.GetIDActive() + "/" + config.Setting("get", "cli", "Spec.Projects.ActiveID", "") + "/"
}

func (f F) loadExisting() {
	f.prepareSettings()
	f.CheckAuth()
	config.Setting("init", "dev", "Spec.", "")
	f.setPaths()
}

func (f F) prepareSettings() {
	f.setIDActive()
	if !config.ValidSettings("cli", "updates", false) {
		cliconfig.F.SetUpdates(cliconfig.F{})
	}
}

func (f F) setIDActive() {
	var idActive string
	var kind string
	if config.ValidSettings("cli", "projects", false) {
		idActive = config.Setting("get", "cli", "Spec.Projects.DefaultID", "")
		kind = "default "
	} else {
		idActive = f.choose()
		set := cliconfig.F.SetDefaultProject(cliconfig.F{}, idActive)
		if set {
			kind = "default "
		}
	}
	config.Setting("set", "cli", "Spec.Projects.ActiveID", idActive)
	logging.Log([]string{"", vars.EmojiSettings, vars.EmojiProject}, "Using "+kind+"project: "+idActive+"\n", 0)
}

func (f F) choose() string {
	return terminal.GetUserSelection("Which project do you want to open?", f.GetAllIDs(), false, false)
}

func (f F) CheckAuth() {
	if !config.ValidSettings("user", "Spec.login", false) {
		config.Setting("set", "user", "Spec.Auth.JWT.SigningKey", jwt.GenerateSigningKey())
	}
}

func (f F) GetAllIDs() []string {
	return filesystem.Explorer("files", vars.ProjectsBasePath, []string{}, []string{"hidden", ".yaml", "infrastructure"})
}

func (f F) createNew(firstProject bool) {
	if firstProject {
		logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "In order to continue you have to create your first project.", 0)
		logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "All your settings are saved as .yaml at: "+filesystem.GetWorkingDir()+vars.NeuraCLINameRepo+"/"+users.BasePath, 0)
		logging.Log([]string{"", vars.EmojiProject, vars.EmojiEdit}, "After the creation you can easy edit them with any editor.", 0)
	}
	config.Setting("reset", "cli", "Spec.Projects.DefaultID", "")
	IDActive = terminal.GetUserInput("Choose a name for your new project:")
	vars.ProjectDirActive = vars.ProjectsBasePath + IDActive
	config.Setting("set", "cli", "Spec.Projects.ActiveID", IDActive)
	f.setPaths()
	//vars.ProjectPathYaml = vars.ProjectDirActive + "/config.yaml"
	filesystem.CreateDir(vars.ProjectDirActive, false)
	config.Setting("set", "project", "Metadata.ID", IDActive)
	config.Setting("set", "project", "Metadata.Name", IDActive)
	//workingDir := terminal.UserSelectionFiles("Where is your project located locally?", "directories", filesystem.GetUserHomeDir(), []string{}, []string{"hidden"}, false, true)
	logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "Please provide the absolute path to your project working directory.", 0)
	logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "E.g. /home/"+runtime.F.GetOSUsername(runtime.F{})+"/projects/"+IDActive, 0)
	workingDir := terminal.GetUserInput("Project working directory path:")
	config.Setting("set", "project", "Spec.WorkingDir", workingDir)
	cliconfig.F.SetDefaultProject(cliconfig.F{}, IDActive)
	cliconfig.F.SetUpdates(cliconfig.F{})
	devconfigAssistant.Create()
}

func (f F) Create(id string) {
	projectPath := vars.ProjectsBasePath + "/" + id
	if !filesystem.Exists(projectPath) {
		filesystem.CreateDir(projectPath, false)
	}
}

func (f F) GetWorkingDir() string {
	var wd string = config.Setting("get", "dev", "Spec.WorkingDir", "")
	if !strings.HasSuffix(wd, "/") {
		wd = wd + "/"
	}
	return wd
}

func (f F) GetExternalExecPath(project string) string {
	var dir string = project + "/" + project
	if !env.F.Container(env.F{}) {
		dir = "../" + dir
	}
	return dir
}

func (f F) GetExternalDataPath() string {
	var dir string = f.GetWorkingDir() + "/data/training/"
	if !env.F.Container(env.F{}) {
		dir = "../" + dir
	}
	return dir
}
