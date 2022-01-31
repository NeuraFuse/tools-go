package client

import (
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// _ "k8s.io/client-go/plugin/pkg/client/auth" // TODO: Load all known auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type F struct{}

func (f F) GetAuth() *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	var err error
	var success bool
	for ok := true; ok; ok = !success {
		clientset, err = kubernetes.NewForConfig(f.GetRestConfig())
		if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create a new Clientset for the given config!", false, false, true) {
			infraConfig.F.SetKubeConfig(infraConfig.F{})
		} else {
			success = true
		}
	}
	return clientset
}

func (f F) GetRestConfig() *rest.Config {
	var restconfig *rest.Config
	var kubeConfigPath string = infraConfig.F.GetInfraKubeAuthPath(infraConfig.F{}, true)
	var err error
	restconfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to build restconfig!", false, false, true) {
		if env.F.CLI(env.F{}) {
			infraConfig.F.SetKubeConfig(infraConfig.F{})
		}
	}
	return restconfig
}
