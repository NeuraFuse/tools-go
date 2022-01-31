package gcloud

import (
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/io"

	// "github.com/neurafuse/tools-go/kubernetes/client/kubeconfig"
	// "github.com/neurafuse/tools-go/kubernetes/namespaces"
	"github.com/neurafuse/tools-go/cloud/providers/gcloud/clusters"
	gconfig "github.com/neurafuse/tools-go/cloud/providers/gcloud/config"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(action string, cliArgs []string) bool { // nodePool("list") , nodePool("create") , nodePool("delete")
	if !config.ValidSettings("infrastructure", vars.InfraProviderGcloud, true) {
		gconfig.F.SetConfigs(gconfig.F{})
	}
	f.checkAPIAvailability()
	logging.ProgressSpinner("start")
	var success bool
	if action == "inspect" {
		f.inspect()
		success = true
	} else if action == "create" {
		success = clusters.F.Create(clusters.F{})
	} else if action == "recreate" {
		clusters.F.Delete(clusters.F{})
		success = clusters.F.Create(clusters.F{})
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid action argument: "+cliArgs[1], true, true, true)
	}
	if !success {
		f.inspect()
	}
	return success
}

func (f F) inspect() {
	logging.Log([]string{"", vars.EmojiInfra, vars.EmojiInspect}, "Starting "+runtime.F.GetCallerInfo(runtime.F{}, true)+" inspection..", 0)
	clusters.F.Get(clusters.F{}, true)
}

func (f F) checkAPIAvailability() {
	var apiURL string = "status.cloud.google.com"
	if io.F.Reachable(io.F{}, apiURL) {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "gcloud is reachable.\n", 0)
	} else {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiWarning}, "Provider "+vars.InfraProviderGcloud+" is not reachable!", 0)
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There is probably an error with networking on your side or at "+vars.InfraProviderGcloud+"!", true, true, true)
	}
}
