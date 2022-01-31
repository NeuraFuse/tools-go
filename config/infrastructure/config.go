package infrastructure

import (
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/exec"
	"github.com/neurafuse/tools-go/filesystem"
	infraID "github.com/neurafuse/tools-go/infrastructures/id"
	kubeID "github.com/neurafuse/tools-go/kubernetes/client/id"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var config *Default
var filePath string
var fileName string = context + format

func (f F) SetConfig() (*Default, string) {
	f.exists()
	return f.GetConfig(), f.GetPath(true)
}

func (f F) exists() {
	f.setTemplate()
	filePath = f.GetPath(true)
	if filesystem.Exists(filePath) {
		yaml.FileToStruct(filePath, &config)
	}
}

func (f F) setTemplate() {
	config = &Default{}
	config.APIVersion = vars.NeuraKubeAPIVersion
	config.Kind = strings.Title(context)
}

func (f F) GetBasePath() string {
	var path string
	var preemble string
	if env.F.CLI(env.F{}) {
		preemble = runtime.F.GetOSInstallDir(runtime.F{}) + vars.NeuraCLINameID + "/" + usersID.BasePath // TODO: Ref.
	} else if env.F.API(env.F{}) {
		if env.F.Container(env.F{}) {
			preemble = f.GetContainerUserPath()
		} else {
			preemble = filesystem.GetWorkingDir(false) + usersID.BasePath
		}
	}
	path = preemble + usersID.F.GetActive(usersID.F{}) + "/" + context + "/"
	return path
}

func (f F) GetPath(withFileName bool) string {
	var path string
	path = f.GetBasePath() + infraID.F.GetActive(infraID.F{}) + "/"
	if withFileName {
		path = path + fileName
	}
	return path
}

func (f F) GetAllIDs() []string {
	return filesystem.Explorer("files", f.GetBasePath(), []string{}, []string{"hidden", ".yaml", "infrastructure"})
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) GetContainerUserPath() string {
	return env.F.GetContainerWorkingDir(env.F{}) + usersID.BasePath
}

func (f F) SetProject() {
	var sel string = terminal.GetUserSelection("What is the "+vars.InfraProviderGcloud+" projectID?", []string{"Example: my-ai-project"}, true, false)
	if sel != "" {
		f.setValue("Spec.Gcloud.ProjectID", sel)
	}
	sel = terminal.GetUserSelection("What is the "+vars.InfraProviderGcloud+" zone?", []string{"Default: us-central1-a"}, true, false)
	if sel != "" {
		f.setValue("Spec.Gcloud.Zone", sel)
	}
}

var providerIDActive string // TODO: Ref.

func (f F) ProviderIDIsActive(id string) bool {
	var status bool
	if id == f.GetProviderIDActive() {
		status = true
	}
	return status
}

func (f F) GetProviderIDActive() string {
	var id string
	id = providerIDActive
	if id == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There is no infrastructure provider set!", true, true, true)
	}
	return id
}

func (f F) SetProviderIDActive(id string) {
	providerIDActive = id
}

func (f F) SetCluster() {
	fieldPrefix := "Spec.Cluster."
	var sel string = terminal.GetUserSelection("What should be the name of the cluster?", []string{"Default: cluster-ai-1"}, true, false)
	if sel != "" {
		f.setValue(fieldPrefix+"Name", sel)
	}
	logging.Log([]string{"\n", vars.EmojiSettings, vars.EmojiKubernetes}, "Please choose the nodes disk size.\nWe recommend about 30 GB for medium sized setups.", 0)
	sel = terminal.GetUserInput("Which disk size should be configured for the cluster nodes? [GB]")
	if sel != "" {
		f.setValue(fieldPrefix+"Nodes.DiskSizeGb", sel)
	}
	f.clusterSelfDeletion()
}

func (f F) clusterSelfDeletion() {
	logging.Log([]string{"\n", vars.EmojiSettings, vars.EmojiKubernetes}, "The self deletion feature enables the automatic release of all hardware related resources by deleting the cluster after n hours of inactivity.", 0)
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "This can be useful to save development costs but keep in mind that all persistent cluster storage will also be deleted with the cluster itself.", 0)
	var sel string = terminal.GetUserSelection("Do you want to enable the cluster self deletion?", []string{}, false, true)
	if sel == "Yes" {
		f.setValue("Spec.Cluster.SelfDeletion.Active", "true")
		var timeDuration string = terminal.GetUserInput("After which number of hours of inactivity should the self deletion occur?")
		f.setValue("Spec.Cluster.SelfDeletion.TimeDurationHours", timeDuration)
	} else {
		f.setValue("Spec.Cluster.SelfDeletion.Active", "false")
	}
}

func (f F) SetMachineSpec() {
	var sel string = terminal.GetUserSelection("Which machine type should be the default?", []string{"Default: e2-standard-4"}, true, false)
	if sel != "" {
		f.setValue("Spec.Gcloud.MachineType", sel)
	}
}

func (f F) SetModuleSpec(context string) {
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Please configure your "+context+" environment.", 0)
	var accType string = terminal.GetUserSelection("Which type of accelerator do you want to choose for "+context+"?", []string{"tpu", "gpu"}, false, false)
	contextID := strings.Title(context)
	f.setValue("Spec."+contextID+".Type", accType)
	var sel string = terminal.GetUserSelection("Which IDE do you want to choose?", []string{"vscode", "other"}, false, false)
	if sel != "" {
		f.setValue("Spec."+contextID+".Environment.IDE", sel)
	}
	sel = terminal.GetUserSelection("Which ML framework do you want to choose?", vars.MLSupportedFrameworks, false, false)
	if sel != "" {
		f.setValue("Spec."+contextID+".Environment.Framework", sel)
	}
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Please choose a volume size for the persistant storage of the accelerators.\nWe recommend about 40 GB for medium sized setups.", 0)
	sel = terminal.GetUserSelection("How big should be the volume size of the accelerator? [GB]", []string{"40", "80", "120"}, true, false)
	if sel != "" {
		f.setValue("Spec."+contextID+".VolumeSizeGB", sel)
	}
	sel = terminal.GetUserSelection("Do you want to enable the accelerator self deletion after n hours of inactivity?", []string{}, false, true)
	if sel == "Yes" {
		f.setValue("Spec."+contextID+".SelfDeletion.Active", "true")
		var timeDuration string = terminal.GetUserInput("After which number of hours of inactivity should the self deletion occur?")
		f.setValue("Spec."+contextID+".SelfDeletion.TimeDurationHours", timeDuration)
	} else {
		f.setValue("Spec."+contextID+".SelfDeletion.Active", "false")
	}
	sel = terminal.GetUserSelection("Should the "+context+" environment have a dedicated node pool (takes more time to spin up)?", []string{}, false, true)
	if sel == "Yes" {
		f.setValue("Spec."+contextID+".NodePools.Dedicated", "true")
	} else {
		f.setValue("Spec."+contextID+".NodePools.Dedicated", "false")
	}
}

func (f F) SetGcloudAccelerator(context, accType string) {
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Please configure your "+vars.InfraProviderGcloud+" "+accType+" accelerator.", 0)
	var nodesDiskSizeGb string = terminal.GetUserSelection("Which disk size should the kubernetes node have that is attached to the accelerator? [GB]", []string{"Default: 70"}, true, false)
	f.setValue("Spec.Gcloud.Accelerator."+strings.ToUpper(accType)+".Node.DiskSizeGb", nodesDiskSizeGb)
	if accType == "tpu" {
		var success bool
		var zone string = f.getValue("Spec.Gcloud.Zone")
		for ok := true; ok; ok = !success {
			var version string = terminal.GetUserSelection("Which "+accType+" version do you want to choose?", []string{"v3", "v2"}, false, false)
			f.setValue("Spec.Gcloud.Accelerator.TPU.Version", version)
			if version == "v2" {
				var cores string = terminal.GetUserSelection("Which number of cores should the TPU have?", []string{"8", "32", "128", "256", "512"}, false, false)
				if cores == "8" {
					var preemptible string = terminal.GetUserSelection("Should the TPU be preemptible for cost saving?", []string{}, false, true)
					if preemptible == "Yes" {
						f.setValue("Spec.Gcloud.Accelerator.TPU.Version", "preemptible-"+version)
					}
					if zone == "us-central1-a" {
						logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiWarning}, "Invalid choice: There is no TPU "+version+"-"+cores+" in your selected zone "+zone+".", 0)
						logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiWarning}, "Please retry with a valid TPU configuration.", 0)
					} else {
						f.setValue("Spec.Gcloud.Accelerator.TPU.Cores", cores)
						success = true
					}
				} else {
					f.setValue("Spec.Gcloud.Accelerator.TPU.Cores", cores)
					success = true
				}
			} else if version == "v3" {
				var preemptible string = terminal.GetUserSelection("Should the TPU be preemptible for cost saving?", []string{}, false, true)
				if preemptible == "Yes" {
					f.setValue("Spec.Gcloud.Accelerator.TPU.Version", "preemptible-"+version)
				}
				f.setValue("Spec.Gcloud.Accelerator.TPU.Cores", "8")
				success = true
			}
		}
		var framework string = f.getValue("Spec." + strings.Title(context) + ".Environment.Framework")
		var tfVersionOptions []string
		if framework == "pytorch" {
			tfVersionOptions = append(tfVersionOptions, "Default: pytorch-1.7")
		} else if framework == "tensorflow" {
			tfVersionOptions = append(tfVersionOptions, "Default: 2.3")
		}
		var tfVersion string = terminal.GetUserSelection("Which TF version should be the default for the TPUs?", tfVersionOptions, true, false)
		f.setValue("Spec.Gcloud.Accelerator.TPU.TF.Version", tfVersion)
		var machineType string = terminal.GetUserSelection("Which machine type (kubernetes node) should be the default for TPUs?", []string{"Default: n1-standard-8"}, true, false)
		f.setValue("Spec.Gcloud.Accelerator.TPU.MachineType", machineType)
	} else if accType == "gpu" {
		var machineType string = terminal.GetUserSelection("Which machine type (kubernetes node) should be the default for GPUs?", []string{"Default: n1-standard-8"}, true, false)
		f.setValue("Spec.Gcloud.Accelerator.GPU.MachineType", machineType)
		gpuTypeOptions := []string{"Default: nvidia-tesla-v100", "nvidia-tesla-t4", "nvidia-tesla-p100", "nvidia-tesla-p4", "nvidia-tesla-k80"}
		var gpuType string = terminal.GetUserSelection("Which GPU model do you want to choose?", gpuTypeOptions, true, false)
		f.setValue("Spec.Gcloud.Accelerator.GPU.Type", gpuType)
	}
}

func (f F) SetNeuraKubeSpec() {
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "Please choose a volume size for the persistant storage of "+vars.NeuraKubeName+".\nWe recommend about 5 GB for medium sized setups.", 0)
	var sel string = terminal.GetUserInput("How big should be the volume size of " + vars.NeuraKubeName + "? [GB]")
	if sel != "" {
		f.setValue("Spec.Neurakube.VolumeSizeGB", sel)
	}
}

func (f F) SetKubeConfig() {
	var clusterName string = f.getValue("Spec.Cluster.ID")
	logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiKubernetes}, "You have to set up the authentication (kubeconfig) for cluster "+clusterName+".\n", 0)
	f.setKubeConfigAuth()
}

func (f F) setKubeConfigAuth() {
	var pathDefault string = filesystem.GetUserHomeDir() + "/.kube/config"
	var optDefaultLoc string = "I already have one placed at " + pathDefault
	var optAssistant string = "How to get it?"
	var optCustomLoc string = "At a custom location"
	var userOpts []string
	userOpts = []string{optDefaultLoc, optAssistant, optCustomLoc}
	var selLoc string = terminal.GetUserSelection("Where is the file located?", userOpts, false, false)
	var filePath string
	var clusterName string = f.getValue("Spec.Cluster.ID")
	if selLoc == optDefaultLoc {
		filePath = pathDefault
	} else if selLoc == optAssistant {
		filePath = pathDefault
		var gcloudZone string = f.getValue("Spec.Gcloud.Zone")
		if f.ProviderIDIsActive("gcloud") {
			var gcloudAlias string = "gcloud"
			var gcloudArgs []string = strings.Split("container clusters get-credentials "+clusterName+" --zone="+gcloudZone, " ")
			var gcloudCmd string = gcloudAlias + " " + strings.Join(gcloudArgs, " ")
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfo}, "The auto. fetching of the kubeconfig independent of the "+vars.InfraProviderGcloud+" SDK is not yet available.", 0)
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfo}, "Therefore you have to initially get the cluster credentials manually via the "+vars.InfraProviderGcloud+" CLI:\n", 0)
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfo}, gcloudCmd+"\n", 0)
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfo}, "The file will be placed at: "+pathDefault+"config", 0)
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfo}, "Further information: https://cloud.google.com/sdk/gcloud/reference/container/clusters/get-credentials", 0)
			var sel string = terminal.GetUserSelection("Do you want "+env.F.GetActive(env.F{}, true)+" to execute the "+vars.InfraProviderGcloud+" command for you?", []string{}, false, true)
			if sel == "Yes" {
				var err error = exec.WithLiveLogs(gcloudAlias, gcloudArgs, true)
				if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to execute glcoud command: "+gcloudCmd, false, true, true) {
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Fetched kubeconfig from glcoud.", 0)
				}
			} else {
				var sel string = terminal.GetUserSelection("Do you have the kubeconfig file now in place (check at "+pathDefault+")", []string{}, false, true)
				if sel != "Yes" {
					logging.Log([]string{"", vars.EmojiSettings, vars.EmojiKubernetes}, "Okay, then try it again with a new run.", 0)
					terminal.Exit(0, "")
				}
			}
		} else if f.ProviderIDIsActive("selfhosted") {
			logging.Log([]string{"", vars.EmojiSettings, vars.EmojiKubernetes}, "There are no routines available yet for selfhosted infrastructures.", 0)
		}
	} else if selLoc == optCustomLoc {
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiKubernetes}, "Okay, then please manually copy your kubeconfig file to your infrastructure auth. located at: "+pathDefault, 0)
		filePath = terminal.GetUserInput("Where is the kubeconfig located [filePath]?")
	}
	if selLoc == optDefaultLoc || selLoc == optAssistant {
		var fileName string = filesystem.GetFileNameFromDirPath(filePath)
		var projectAuthPath string = f.GetInfraKubeAuthPath(false) + fileName
		if filesystem.Exists(projectAuthPath) {
			filesystem.Delete(projectAuthPath, false)
		}
		filesystem.Copy(filePath, projectAuthPath, false)
	}
}

var clusterRecentlyDeleted bool

func (f F) SetClusterRecentlyDeleted(deleted bool) {
	clusterRecentlyDeleted = deleted
}

func (f F) GetClusterRecentlyDeleted() bool {
	return clusterRecentlyDeleted
}

func (f F) GetInfraKubePath() string {
	return f.GetPath(false) + "clusters/"
}

func (f F) GetInfraKubeAuthPath(withFileName bool) string {
	var path string
	path = f.GetInfraKubePath() + kubeID.F.GetActive(kubeID.F{}) + "/auth/"
	if withFileName {
		path = path + "config"
	}
	return path
}

func (f F) GetInfraGcloudPath() string {
	return f.GetPath(false) + "public-clouds/gcloud/"
}

func (f F) GetInfraGcloudAuthPath() string {
	return f.GetInfraGcloudPath() + "auth/service-account.json"
}

func (f F) setValue(key string, value string) {
	yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(config, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}
