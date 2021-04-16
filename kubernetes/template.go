package kubernetes

import (
	"./pods"
	"./deployments"
	"./services"
	"./volumes"
	"./namespaces"
	"./nodes"
	"./daemonsets"
)

type ResourceTypes struct {
	Pods pods.F
	Deployments deployments.F
	Services services.F
	Volumes volumes.F
	Namespaces namespaces.F
	Nodes nodes.F
	Daemonsets daemonsets.F
}