package kubernetes

import (
	"github.com/neurafuse/tools-go/kubernetes/daemonsets"
	"github.com/neurafuse/tools-go/kubernetes/deployments"
	"github.com/neurafuse/tools-go/kubernetes/namespaces"
	"github.com/neurafuse/tools-go/kubernetes/nodes"
	"github.com/neurafuse/tools-go/kubernetes/pods"
	"github.com/neurafuse/tools-go/kubernetes/services"
	"github.com/neurafuse/tools-go/kubernetes/volumes"
)

type ResourceTypes struct {
	Pods        pods.F
	Deployments deployments.F
	Services    services.F
	Volumes     volumes.F
	Namespaces  namespaces.F
	Nodes       nodes.F
	Daemonsets  daemonsets.F
}
