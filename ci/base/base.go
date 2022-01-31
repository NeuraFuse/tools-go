package base

import (
	"github.com/neurafuse/tools-go/apps/python/debugpy"
	"github.com/neurafuse/tools-go/apps/python/flask"
	"github.com/neurafuse/tools-go/apps/tensorflow/tensorboard"
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/kubernetes/namespaces"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	usersID "github.com/neurafuse/tools-go/users/id"
)

type F struct{}

func (f F) GetVolumeSizeGB(context string) string {
	return config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".VolumeSizeGB", "") + "Gi"
}

func (f F) GetVolumes(context string) [][]string {
	return [][]string{{f.GetVolumeMountPath(), f.GetVolumeSizeGB(context)}}
}

func (f F) GetVolumeMountPath() string {
	var projectName string = config.Setting("get", "project", "Metadata.Name", "")
	var containerAppDataRoot string = config.Setting("get", "project", "Spec.Containers.Sync.PathMappings.ContainerAppDataRoot", "")
	return env.F.GetContainerWorkingDir(env.F{}) + "/" + projectName + "/" + containerAppDataRoot
}

func (f F) GetNamespace() string {
	return namespaces.Default + "-" + usersID.F.GetActive(usersID.F{})
}

func (f F) GetResType(context string) string {
	return f.GetEnvFramework(context, false) + " " + context + " environment [" + f.GetType(context, false) + "]"
}

func (f F) GetEnvFramework(context string, titleCase bool) string {
	var envFramework string = config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Environment.Framework", "")
	if titleCase {
		envFramework = strings.Title(envFramework)
	}
	return envFramework
}

func (f F) GetResources(context string) string {
	return config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Type", "")
}

func (f F) GetContainerPorts(context string) [][]string {
	var apps []string
	var appsBase []string = []string{"debugpy", "tensorboard"}
	apps = appsBase
	switch context {
	case "inference":
		apps = []string{"flask"}
	}
	return f.GetContainerPortsForApps(apps)
}

func (f F) GetContainerPortsForApps(apps []string) [][]string {
	var containerPorts [][]string
	for _, app := range apps {
		switch app {
		case "debugpy":
			{
				containerPorts = append(containerPorts, debugpy.GetContainerPorts())
			}
		case "tensorboard":
			{
				containerPorts = append(containerPorts, tensorboard.F.GetContainerPorts(tensorboard.F{}))
			}
		case "flask":
			{
				containerPorts = append(containerPorts, flask.GetContainerPorts())
			}
		default:
			{
				errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There are no registred ports for the "+app+"!", true, true, true)
			}
		}
	}
	return containerPorts
}

func (f F) GetType(context string, upperCase bool) string {
	var resType string = config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Type", "")
	if upperCase {
		resType = strings.ToUpper(resType)
	}
	return resType
}
