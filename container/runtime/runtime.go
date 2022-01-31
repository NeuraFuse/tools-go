package container

import (
	"github.com/neurafuse/tools-go/api/client"
	"github.com/neurafuse/tools-go/apps/tensorflow/tensorboard"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/projects"
	"github.com/neurafuse/tools-go/random"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"

	//"../data/providers/twitter"
	"github.com/neurafuse/tools-go/container/runtime/env/python"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Starting container..\n", 0)
	f.init()
	var project string
	var assistant bool
	project, assistant = f.getProject(cliArgs)
	var module string = f.getModule(cliArgs, project, assistant)
	var codeFormat string = "py"
	var pathExec string = projects.F.GetContainerExternalExecPath(projects.F{}, project, codeFormat)
	var dataPath string = projects.F.GetContainerExternalDataPath(projects.F{}, project, codeFormat)
	f.apps(module)
	python.F.Router(python.F{}, project, module, pathExec, dataPath, f.GetServerSyncWaitMsg())
}

func (f F) init() {
	if env.F.Container(env.F{}) {
		// f.connectAPI() TODO: Debug
	}
}

func (f F) apps(module string) {
	if module != "modelserver" {
		tensorboard.F.Start(tensorboard.F{})
	}
}

func (f F) connectAPI() {
	usersID.F.SetActive(usersID.F{}, f.getID())
	projects.F.CheckAuth(projects.F{})
	client.F.Connect(client.F{}, vars.NeuraKubeNameID)
}

func (f F) getID() string {
	return "container-" + random.GetString(8)
}

func (f F) getProject(cliArgs []string) (string, bool) {
	var project string
	var assistant bool = true
	if len(cliArgs) > 1 {
		project = cliArgs[1]
		if project == "lightning-py" {
			assistant = false
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unsupported project argument: "+project, true, false, true)
		}
	}
	if assistant {
		project = terminal.GetUserSelection("Please choose a project", []string{"lightning-py"}, false, false)
	}
	return project, assistant
}

func (f F) getModule(cliArgs []string, project string, assistant bool) string {
	var module string
	if len(cliArgs) < 3 || assistant {
		module = terminal.GetUserSelection("Please choose a module for the project "+project+"", []string{"gpt", "modelserver"}, false, false)
	} else if len(cliArgs) >= 3 {
		module = cliArgs[2]
	}
	return module
}

func (f F) GetServerSyncWaitMsg() string {
	return "Waiting for project to be synced to the server.."
}
