package config

import (
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"

	//"github.com/neurafuse/tools-go/kubernetes/deployments/templates"

	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type F struct{}

var infraProviderActiveLogged bool

func (f F) SetConfigs() {
	if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, true) {
		if !infraProviderActiveLogged {
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfra}, "Active infrastructure provider: "+vars.InfraProviderGcloud+"\n", 0)
			infraProviderActiveLogged = true
		}
	} else {
		logging.Log([]string{"\n", vars.EmojiInfra, vars.EmojiWarning}, "Invalid and/or missing "+vars.InfraProviderGcloud+" settings in infrastructure config.\nStarting assistant to retrieve missing settings..\n", 0)
		infraConfig.F.SetProject(infraConfig.F{})
		infraConfig.F.SetMachineSpec(infraConfig.F{})
	}
	infraConfig.F.SetProviderIDActive(infraConfig.F{}, "gcloud")
}

func (f F) ClusterConfig(machineType, accType string) *containerpb.Cluster {
	cluster := &containerpb.Cluster{}
	cluster.Name = config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")
	cluster.Description = "Created by " + vars.OrganizationName + " from user " + usersID.F.GetActive(usersID.F{}) + "."
	cluster.NodePools = f.NodePoolConfig(vars.OrganizationNameRepo, machineType, accType)
	cluster.MasterAuth = f.GetMasterAuth()
	cluster.IpAllocationPolicy = f.GetIPAllocationPolicy()
	cluster.EnableTpu = true
	return cluster
}

func (f F) GetIPAllocationPolicy() *containerpb.IPAllocationPolicy {
	ipAP := &containerpb.IPAllocationPolicy{}
	ipAP.UseIpAliases = true
	return ipAP
}

func (f F) GetMasterAuth() *containerpb.MasterAuth {
	var masterAuth *containerpb.MasterAuth
	return masterAuth
}

func (f F) NodePoolConfigSingle(context, machineType, accType string) *containerpb.NodePool {
	nodePool := &containerpb.NodePool{}
	nodePool.Name = context
	nodePool.Config = f.NodeConfig(context, machineType, accType)
	nodePool.InitialNodeCount = 1
	nodePool.Autoscaling = f.NodePoolAutoscalig()
	return nodePool
}

func (f F) NodePoolConfig(context, machineType, accType string) []*containerpb.NodePool {
	nodePool := []*containerpb.NodePool{&containerpb.NodePool{}}
	nodePool[0].Name = context
	nodePool[0].Config = f.NodeConfig(context, machineType, accType)
	nodePool[0].InitialNodeCount = 1
	nodePool[0].Autoscaling = f.NodePoolAutoscalig()
	return nodePool
}

func (f F) NodePoolAutoscalig() *containerpb.NodePoolAutoscaling {
	nodePoolAutoscalig := &containerpb.NodePoolAutoscaling{}
	nodePoolAutoscalig.MinNodeCount = 1
	nodePoolAutoscalig.MaxNodeCount = 1
	nodePoolAutoscalig.Enabled = true
	return nodePoolAutoscalig
}

func (f F) NodeConfig(context, machineType, accType string) *containerpb.NodeConfig {
	nodeConfig := &containerpb.NodeConfig{}
	nodeConfig.MachineType = machineType
	var diskSizeGbKey string
	if accType != "" {
		diskSizeGbKey = "Spec.Gcloud.Accelerator." + strings.ToUpper(accType) + ".Node.DiskSizeGb"
	} else {
		diskSizeGbKey = "Spec.Cluster.Nodes.DiskSizeGb"
	}
	nodeConfig.DiskSizeGb = strings.ToInt32(config.Setting("get", "infrastructure", diskSizeGbKey, ""))
	nodeConfig.OauthScopes = []string{"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/devstorage.full_control",
	}
	nodeConfig.Labels = map[string]string{"service-cluster": context}
	if accType == "gpu" {
		nodeConfig.Accelerators = f.AcceleratorConfigGPU()
	}
	nodeConfig.DiskType = "pd-ssd"
	nodeConfig.Preemptible = true
	return nodeConfig
}

/*func (f F) NodeTaint(accelerator bool) []*containerpb.NodeTaint {
	nodeTaint := []*containerpb.NodeTaint{&containerpb.NodeTaint{}}
	if accelerator {
		nodeTaint[0].Key = templates.GPUKey
		nodeTaint[0].Value = "1"
		nodeTaint[0].Effect = containerpb.NodeTaint_Effect(int32(2))
	}
	return nodeTaint
}*/

func (f F) AcceleratorConfigGPU() []*containerpb.AcceleratorConfig {
	acceleratorConfig := []*containerpb.AcceleratorConfig{&containerpb.AcceleratorConfig{}}
	acceleratorConfig[0].AcceleratorType = config.Setting("get", "infrastructure", "Spec.Gcloud.Accelerator.GPU.Type", "")
	acceleratorConfig[0].AcceleratorCount = 1
	return acceleratorConfig
}
