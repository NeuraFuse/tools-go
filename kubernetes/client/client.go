package client

import (
	"../../../neurakube/infrastructure/providers/gcloud/clusters"
	"../../config"
	"../../errors"
	"../../runtime"
	"../../vars"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type F struct{}

func (f F) GetAuth() *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(f.GetRestConfig())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create a new Clientset for the given config!", false, false, true)
	return clientset
}

func (f F) GetRestConfig() *rest.Config {
	var restconfig *rest.Config
	var err error
	if config.ValidSettings("infrastructure", "kube", false) {
		restconfig, err = clientcmd.BuildConfigFromFlags("", config.Setting("get", "infrastructure", "Spec.Cluster.Auth.KubeConfigPath", ""))
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to BuildConfigFromFlags!", false, true, true)
	} else if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, true) {
		restconfig = clusters.F.GetAuthConfig(clusters.F{})
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Missing valid settings for kubernetes client auth!", true, false, true)
	}
	return restconfig
}
