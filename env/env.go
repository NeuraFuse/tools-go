package env

import (
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var actionActive string

func (f F) GetVersion() string {
	var version string
	if f.CLI() {
		version = vars.NeuraCLIVersion
	} else if f.API() {
		version = vars.NeuraKubeVersion
	}
	return version
}

func (f F) GetContainerWorkingDir() string {
	return "/app/"
}

func (f F) Container() bool {
	workingDir := filesystem.GetWorkingDir(false)
	if workingDir == f.GetContainerWorkingDir() {
		return true
	} else {
		return false
	}
}

func (f F) GetContext(context string, title bool) string {
	if title {
		context = strings.Title(context)
	}
	return context
}

func (f F) GetAPIHTTPCertPath() string {
	return "server/http/certs/"
}

func (f F) CLI() bool {
	return f.IsFrameworkActive(f.GetID("neuracli"))
}

func (f F) API() bool {
	return f.IsFrameworkActive(f.GetID("neurakube"))
}

func (f F) Develop() bool {
	return f.IsActionActive(f.GetID("develop"))
}

func (f F) App() bool {
	return f.IsActionActive(f.GetID("app"))
}

func (f F) Inference() bool {
	return f.IsActionActive(f.GetID("inference"))
}

func (f F) GetID(name string) string {
	var id string
	switch name {
	case "neuracli":
		id = vars.NeuraCLINameID
	case "neurakube":
		id = vars.NeuraKubeNameID
	case "develop":
		id = "develop"
	case "app":
		id = "app"
	case "inference":
		id = "inference"
	}
	return id
}

func (f F) IsFrameworkActive(envName string) bool {
	if envName == vars.FrameworkEnvActive {
		return true
	} else {
		return false
	}
}

func (f F) IsActionActive(action string) bool {
	if actionActive == action {
		return true
	} else {
		return false
	}
}

func (f F) GetActive(caseTitle bool) string {
	var env string
	if caseTitle {
		env = strings.Title(vars.FrameworkEnvActive)
	} else {
		env = vars.FrameworkEnvActive
	}
	return env
}

func (f F) SetFramework(module string) {
	vars.FrameworkEnvActive = module
}

func (f F) SetAction(action string) {
	actionActive = action
}
