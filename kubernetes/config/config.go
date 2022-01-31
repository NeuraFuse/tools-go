package config

import (
	"github.com/neurafuse/tools-go/cloud/providers/gcloud/clusters"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/kubernetes/namespaces"
	"github.com/neurafuse/tools-go/runtime"
)

type F struct{}

func (f F) CheckKubeAuth() bool {
	var setKubeConfig bool
	if filesystem.Exists(infraConfig.F.GetInfraKubeAuthPath(infraConfig.F{}, true)) {
		var err error
		err, _ = namespaces.F.Get(namespaces.F{}, false)
		if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to authenticate with current kubeconfig at the cluster!", false, false, true) {
			setKubeConfig = true
		}
	} else {
		setKubeConfig = true
	}
	var devConfigSkip bool
	devConfigSkip = f.CheckResources()
	if setKubeConfig && !devConfigSkip {
		infraConfig.F.SetKubeConfig(infraConfig.F{})
	}
	return !devConfigSkip
}

func (f F) CheckResources() bool {
	var devConfigSkip bool
	if infraConfig.F.ProviderIDIsActive(infraConfig.F{}, "gcloud") {
		var exists bool
		exists, _ = clusters.F.Exists(clusters.F{})
		if !exists {
			devConfigSkip = clusters.F.ResourceMissing(clusters.F{})
		}
	}
	return devConfigSkip
}
