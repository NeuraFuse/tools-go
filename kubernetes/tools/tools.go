package tools

import (
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/kubernetes/deployments"
	"github.com/neurafuse/tools-go/kubernetes/namespaces"
	"github.com/neurafuse/tools-go/kubernetes/pods"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) GetContainerID(namespace, appID string) int {
	var containerID int
	if len(pods.F.GetContainers(pods.F{}, namespace, appID)) > 1 {
		var containerName string = terminal.GetUserSelection("To which container pod do you want to connect?", pods.F.GetContainerNamesList(pods.F{}, namespace, appID), false, false)
		containerID = pods.F.GetContainerIDByName(pods.F{}, namespace, appID, containerName)
	}
	return containerID
}

func (f F) GetDeploymentNamespace(appID string) string {
	var namespace string
	var namespaceLocs []string
	var err error
	err, namespaceLocs = f.GetDeploymentNamespaces(appID)
	if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to automatically locate the deployment "+appID+" in the cluster!", false, true, true) {
		namespace = terminal.GetUserInput("In which cluster namespace is your app " + appID + " deployed?")
	} else {
		if len(namespaceLocs) == 1 {
			namespace = namespaceLocs[0]
		} else {
			logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "There are multiple namespaces that are containing a deployment named "+appID+".", 0)
			namespace = terminal.GetUserSelection("Please choose the right namespace", namespaceLocs, false, false)
		}
	}
	return namespace
}

func (f F) GetDeploymentNamespaces(id string) (error, []string) {
	var targetNamespaceIDs []string
	var ns []string
	var err error
	err, ns = namespaces.F.Get(namespaces.F{}, false)
	for _, namespaceID := range ns {
		var deployments []string = deployments.F.Get(deployments.F{}, namespaceID, false)
		for _, deploymentID := range deployments {
			if deploymentID == id {
				targetNamespaceIDs = append(targetNamespaceIDs, namespaceID)
				break
			}
		}
	}
	if len(targetNamespaceIDs) == 0 {
		err = errors.New("Unable to find namespace for deployment " + id + "!")
	}
	return err, targetNamespaceIDs
}
