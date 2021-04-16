package kubernetes

import (
	"strings"

	"../errors"
	"../logging"
	"../nlp"
	"../objects"
	"../runtime"
	"../terminal"
	"../vars"
	"./deployments"
	"./namespaces"
)

type F struct{}

var actionTypes []string = []string{"get", "logs", "inspect", "create", "delete"}
var resourceTypes []string = []string{"pod", "deployment", "service", "volume", "namespace", "node", "daemonset"}

func (f F) Router(cliArgs []string, routeAssistant bool) { // TODO: Bugfix logs [namespace] ..
	var action string = ""
	var namespace string = ""
	if routeAssistant || len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which "+runtime.F.GetCallerInfo(runtime.F{}, true)+" action do you want to start?", actionTypes, false, false)
	} else {
		action = cliArgs[1]
	}
	lenCliArgs := len(cliArgs)
	namespaceMissing := false
	if len(cliArgs) <= 2 || len(cliArgs) < 4 {
		if action != "inspect" && action != "logs" && action != "get" {
			namespace = terminal.GetUserInput("Which namespace do you want to select for the action " + action + "?")
		} else if len(cliArgs) < 3 || (len(cliArgs) <= 3 && action != "pods" && action != "logs") {
			namespace = namespaces.Default
		} else {
			namespace = cliArgs[2]
		}
		lenCliArgs = lenCliArgs + 1
		namespaceMissing = true
	} else {
		namespace = cliArgs[2]
	}
	if action == "inspect" {
		f.inspect(namespace)
	}
	var resourceType string = ""
	var actionTypeMultiple bool = false
	if routeAssistant || lenCliArgs < 4 {
		if action == "create" || action == "delete" {
			resourceType = terminal.GetUserSelection("On which resource type do you want to call the action "+action+"?", resourceTypes, false, false)
		} else if action == "logs" {
			resourceType = "pods"
		} else {
			actionTypeMultiple = true
			resourceType = terminal.GetUserSelection("On which resource type do you want to call the action "+action+"?", resourceTypes, false, false)
		}
	} else {
		if action == "logs" {
			resourceType = "pods"
		} else {
			if namespaceMissing {
				lenCliArgs = lenCliArgs - 2
			} else {
				lenCliArgs = lenCliArgs - 1
			}
			resourceType = cliArgs[lenCliArgs]
			resourceType, actionTypeMultiple = nlp.ConvertToPlural(resourceType)
		}
	}
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInspect}, nlp.VerbToAction(strings.Title(action))+"ing "+resourceType+"..", 0)
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Namespace: "+namespace, 0)
	var resourceID string = ""
	if !actionTypeMultiple {
		index := 5
		if action == "logs" {
			index = 4
		}
		if routeAssistant || len(cliArgs) < index {
			deployments := deployments.F.GetList(deployments.F{}, namespace, false)
			if len(deployments) != 0 {
				resourceID = terminal.GetUserSelection("On which specific "+resourceType+" do you want to call the action "+action+"?", deployments, false, false)
			} else {
				logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "There are no existing deployments and therefore no pods to stream logs.", 0)
				terminal.Exit(0, "")
			}
		} else {
			resourceID = cliArgs[index-1]
		}
	}
	success := false
	if actionTypeMultiple {
		if resourceType == "volumes" {
			success, _ = objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), namespace, "pvc", true)
		} else {
			success, _ = objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), namespace, true)
		}
	} else {
		if action == "logs" {
			success, _ = objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), namespace, resourceID, "", false, 5)
		} else {
			success, _ = objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), namespace, resourceID)
		}
	}
	if !success {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Either the combination between your selected resourceType ("+resourceType+") and action ("+action+")\nis invalid or the operation is not supported yet by "+vars.NeuraCLIName+".", true, true, true)
	}
}

func (f F) inspect(namespace string) {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInspect}, "Starting "+runtime.F.GetCallerInfo(runtime.F{}, true)+" inspection..", 0)
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiInfo}, "Namespace: "+namespace, 0)
	action := "get"
	for _, resourceType := range resourceTypes {
		resourceType, _ = nlp.ConvertToPlural(resourceType)
		if resourceType == "volumes" {
			objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), namespace, "pvc", true)
		} else if resourceType == "namespaces" || resourceType == "nodes" || resourceType == "daemonsets" {
			objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), true)
		} else {
			objects.CallStructInterfaceFuncByName(ResourceTypes{}, strings.Title(resourceType), strings.Title(action), namespace, true)
		}
	}
	terminal.Exit(0, "")
}
