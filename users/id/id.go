package id

import (
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var BasePath string = "users/"
var idActive string

func (f F) ActiveIsSet() bool {
	var status bool
	if idActive != "" {
		status = true
	}
	return status
}

func (f F) CreateNew() {
	var quest string = "Type in a username:"
	var userID string = terminal.GetUserInput(quest)
	f.SetActive(userID)
	var activeUserPath string = f.getActiveUserPath(userID)
	filesystem.CreateDir(activeUserPath, false)
}

func (f F) SetActive(id string) {
	idActive = id
	vars.ProjectsBasePath = filesystem.GetWorkingDir(false) + "." + vars.NeuraCLINameID + "/"
}

func (f F) GetActive() string {
	if idActive == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There is no active user ID set!", true, true, true)
	}
	return idActive
}

func (f F) getActiveUserPath(id string) string {
	var filePath string = BasePath + "/" + id
	if env.F.CLI(env.F{}) {
		filePath = runtime.F.GetOSInstallDir(runtime.F{}) + env.F.GetActive(env.F{}, false) + "/" + filePath
	}
	return filePath
}

func (f F) GetAllIDs() []string {
	return filesystem.Explorer("files", BasePath, []string{}, []string{"hidden", ".yaml"})
}
