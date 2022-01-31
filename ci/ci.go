package ci

import (
	"github.com/neurafuse/tools-go/cloud/providers/gcloud/nodepools"
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/kubernetes/deployments"
	"github.com/neurafuse/tools-go/kubernetes/resources"
	"github.com/neurafuse/tools-go/kubernetes/services"
	"github.com/neurafuse/tools-go/kubernetes/volumes"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/random"
	"github.com/neurafuse/tools-go/releases"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var contextLocal string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var module string
	var modules []string = []string{"releases"}
	if len(cliArgs) < 2 {
		module = terminal.GetUserSelection("Which "+contextLocal+" module do you want to start?", modules, false, false)
	} else {
		module = cliArgs[1]
	}
	switch module {
	case modules[0]:
		releases.F.Router(releases.F{}, cliArgs)
	}
}

func (f F) Exists(namespace, context string) bool {
	var contextID string = f.GetContextID()
	if !volumes.F.Exists(volumes.F{}, namespace, contextID) {
		return false
	} else if !deployments.F.Exists(deployments.F{}, namespace, contextID) {
		return false
	} else if !services.F.Exists(services.F{}, namespace, contextID) {
		return false
	}
	return true
}

func (f F) Create(namespace, context, imageAddrs, accType, clusterIP string, volumesSpec, containerPorts [][]string) {
	if accType != "" {
		resources.Check(context, accType)
	}
	var contextID string = f.GetContextID()
	volumes.F.Create(volumes.F{}, namespace, contextID, f.getServiceCluster(context), volumesSpec)
	var repoAddrs string = config.Setting("get", "dev", "Spec.Containers.Registry.Address", "")
	if repoAddrs != "" {
		imageAddrs = repoAddrs + "/" + imageAddrs
	}
	deployments.F.Create(deployments.F{}, namespace, contextID, imageAddrs, f.getServiceCluster(context), accType, volumesSpec, containerPorts)
	services.F.Create(services.F{}, namespace, contextID, clusterIP, containerPorts)
}

func (f F) NodeScheduling(context string) string {
	go f.CreateNodePool(context)
	return "success"
}

func (f F) CreateNodePool(context string) {
	if infraConfig.F.ProviderIDIsActive(infraConfig.F{}, "gcloud") {
		nodepools.F.Create(nodepools.F{}, f.getServiceCluster(context), f.GetType(context, false))
	} else {
		errors.Check(nil, contextLocal, "Unable to create nodepool for selfhosted setup!", true, false, true)
	}
}

func (f F) DeleteNodePool(context string) {
	if infraConfig.F.ProviderIDIsActive(infraConfig.F{}, "gcloud") {
		nodepools.F.Delete(nodepools.F{}, f.getServiceCluster(context))
	} else {
		errors.Check(nil, contextLocal, "Unable to delete nodepool for selfhosted setup!", true, false, true)
	}
}

func (f F) RecreateNodePool(context string) {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Recreating nodepool "+context+"..", 0)
	f.DeleteNodePool(context)
	f.CreateNodePool(context)
}

func (f F) Update(namespace, context, imageAddrs, accType string, volumes, containerPorts [][]string) {
	if accType != "" {
		resources.Check(context, accType)
	}
	var contextID string = f.GetContextID()
	deployments.F.Delete(deployments.F{}, namespace, contextID)
	var repoAddrs string = config.Setting("get", "dev", "Spec.Containers.Registry.Address", "")
	if repoAddrs != "" {
		imageAddrs = repoAddrs + "/" + imageAddrs
	}
	deployments.F.Create(deployments.F{}, namespace, contextID, imageAddrs, f.getServiceCluster(context), accType, volumes, containerPorts)
}

func (f F) Delete(namespace, context string, volumesSpec [][]string) {
	var contextID string = f.GetContextID()
	deployments.F.Delete(deployments.F{}, namespace, contextID)
	volumes.F.Delete(volumes.F{}, namespace, contextID, volumesSpec)
	services.F.Delete(services.F{}, namespace, contextID)
}

func (f F) getServiceCluster(context string) string {
	var serviceCluster string = vars.OrganizationNameRepo
	if context != vars.NeuraKubeNameID {
		dedicated := config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".NodePools.Dedicated", "")
		if dedicated == "true" {
			serviceCluster = serviceCluster + "-" + context + "-" + f.GetType(context, false)
		} else {
			serviceCluster = serviceCluster + "-" + f.GetType(context, false)
		}
	}
	return serviceCluster
}

func (f F) GetContextID() string {
	return config.Setting("get", "project", "Metadata.ID", "")
}

func (f F) GetClusterIP(min, max int) string {
	var baseIP string = "10.24.0."
	return baseIP + strings.ToString(random.GetInt(min, max))
}

func (f F) GetInitWaitDuration(context string) int {
	var waitDuration int
	var accType string = f.GetType(context, false)
	if accType == "tpu" {
		waitDuration = 20
	} else {
		waitDuration = 8
	}
	return waitDuration
}

func (f F) GetType(context string, upperCase bool) string {
	var resType string = config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Type", "")
	if upperCase {
		resType = strings.ToUpper(resType)
	}
	return resType
}
