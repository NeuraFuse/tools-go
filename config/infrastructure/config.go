package infrastructure

import (
	"../../env"
	"../../filesystem"
	"../../logging"
	"../../objects"
	"../../objects/strings"
	"../../readers/yaml"
	"../../runtime"
	"../../terminal"
	"../../users"
	"../../vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
var format string = ".yaml"

var config *Default
var filePath string
var fileName string = context + format

func (f F) SetConfig() (*Default, string) {
	f.exists()
	return f.GetConfig(), f.GetFilePath()
}

func (f F) exists() {
	f.setTemplate()
	filePath = f.GetFilePath()
	if filesystem.Exists(filePath) {
		yaml.FileToStruct(filePath, &config)
	}
}

func (f F) setTemplate() {
	config = &Default{}
	config.APIVersion = vars.NeuraKubeVersion
	config.Kind = strings.Title(context)
}

func (f F) GetFilePath() string {
	return users.BasePath + "/" + users.GetIDActive() + "/" + context + "/" + fileName
}

func (f F) GetConfig() *Default {
	return config
}

func (f F) SetProject() {
	sel := terminal.GetUserSelection("What is the gcloud projectID?", []string{"Example: my-ai-project"}, true, false)
	if sel != "" {
		f.setValue("Spec.Gcloud.ProjectID", sel)
	}
	sel = terminal.GetUserSelection("What is the gcloud zone?", []string{"Default: us-central1-a"}, true, false)
	if sel != "" {
		f.setValue("Spec.Gcloud.Zone", sel)
	}
}

func (f F) SetAuth() {
	sel := terminal.UserSelectionFiles("Which gcloud service account json should be the default for your projects?", "files", f.getProjectAuthPath(), []string{"json"}, []string{"hidden"}, false, true)
	if sel != "" {
		f.setValue("Spec.Gcloud.Auth.ServiceAccountJSONPath", sel)
	}
}

func (f F) SetCluster() {
	fieldPrefix := "Spec.Cluster."
	sel := terminal.GetUserSelection("What should be the name of the cluster?", []string{"Default: cluster-ai-1"}, true, false)
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
	sel := terminal.GetUserSelection("Do you want to enable the cluster self deletion?", []string{}, false, true)
	if sel == "Yes" {
		f.setValue("Spec.Cluster.SelfDeletion.Active", "true")
		timeDuration := terminal.GetUserInput("After which number of hours of inactivity should the self deletion occur?")
		f.setValue("Spec.Cluster.SelfDeletion.TimeDurationHours", timeDuration)
	} else {
		f.setValue("Spec.Cluster.SelfDeletion.Active", "false")
	}
}

func (f F) SetMachineSpec() {
	sel := terminal.GetUserSelection("Which machine type should be the default?", []string{"Default: e2-standard-4"}, true, false)
	if sel != "" {
		f.setValue("Spec.Gcloud.MachineType", sel)
	}
}

func (f F) SetModuleSpec(context string) {
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Please configure your "+context+" environment.", 0)
	accType := terminal.GetUserSelection("Which type of accelerator do you want to choose for "+context+"?", []string{"tpu", "gpu"}, false, false)
	contextID := strings.Title(context)
	f.setValue("Spec."+contextID+".Type", accType)
	sel := terminal.GetUserSelection("Which IDE do you want to choose?", []string{"vscode", "other"}, false, false)
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
		timeDuration := terminal.GetUserInput("After which number of hours of inactivity should the self deletion occur?")
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
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Please configure your gcloud "+accType+" accelerator.", 0)
	nodesDiskSizeGb := terminal.GetUserSelection("Which disk size should the kubernetes node have that is attached to the accelerator? [GB]", []string{"Default: 70"}, true, false)
	f.setValue("Spec.Gcloud.Accelerator."+strings.ToUpper(accType)+".Node.DiskSizeGb", nodesDiskSizeGb)
	if accType == "tpu" {
		var success bool
		zone := f.getValue("Spec.Gcloud.Zone")
		for ok := true; ok; ok = !success {
			version := terminal.GetUserSelection("Which "+accType+" version do you want to choose?", []string{"v3", "v2"}, false, false)
			f.setValue("Spec.Gcloud.Accelerator.TPU.Version", version)
			if version == "v2" {
				cores := terminal.GetUserSelection("Which number of cores should the TPU have?", []string{"8", "32", "128", "256", "512"}, false, false)
				if cores == "8" {
					preemptible := terminal.GetUserSelection("Should the TPU be preemptible for cost saving?", []string{}, false, true)
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
				preemptible := terminal.GetUserSelection("Should the TPU be preemptible for cost saving?", []string{}, false, true)
				if preemptible == "Yes" {
					f.setValue("Spec.Gcloud.Accelerator.TPU.Version", "preemptible-"+version)
				}
				f.setValue("Spec.Gcloud.Accelerator.TPU.Cores", "8")
				success = true
			}
		}
		framework := f.getValue("Spec." + strings.Title(context) + ".Environment.Framework")
		var tfVersionOptions []string
		if framework == "pytorch" {
			tfVersionOptions = append(tfVersionOptions, "Default: pytorch-1.7")
		} else if framework == "tensorflow" {
			tfVersionOptions = append(tfVersionOptions, "Default: 2.3")
		}
		tfVersion := terminal.GetUserSelection("Which TF version should be the default for the TPUs?", tfVersionOptions, true, false)
		f.setValue("Spec.Gcloud.Accelerator.TPU.TF.Version", tfVersion)
		machineType := terminal.GetUserSelection("Which machine type (kubernetes node) should be the default for TPUs?", []string{"Default: n1-standard-8"}, true, false)
		f.setValue("Spec.Gcloud.Accelerator.TPU.MachineType", machineType)
	} else if accType == "gpu" {
		machineType := terminal.GetUserSelection("Which machine type (kubernetes node) should be the default for GPUs?", []string{"Default: n1-standard-8"}, true, false)
		f.setValue("Spec.Gcloud.Accelerator.GPU.MachineType", machineType)
		gpuTypeOptions := []string{"Default: nvidia-tesla-v100", "nvidia-tesla-t4", "nvidia-tesla-p100", "nvidia-tesla-p4", "nvidia-tesla-k80"}
		gpuType := terminal.GetUserSelection("Which GPU model do you want to choose?", gpuTypeOptions, true, false)
		f.setValue("Spec.Gcloud.Accelerator.GPU.Type", gpuType)
	}
}

func (f F) SetNeuraKubeSpec() {
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "Please choose a volume size for the persistant storage of "+vars.NeuraKubeName+".\nWe recommend about 5 GB for medium sized setups.", 0)
	sel := terminal.GetUserInput("How big should be the volume size of " + vars.NeuraKubeName + "? [GB]")
	if sel != "" {
		f.setValue("Spec.Neurakube.VolumeSizeGB", sel)
	}
}

func (f F) validKubeConfig() bool {
	var valid bool
	fieldsPrefix := "Spec.Cluster."
	fields := []string{"KubeConfigPath"}
	valid = objects.StructFieldValuesExisting(config, fieldsPrefix, fields, runtime.F.GetCallerInfo(runtime.F{}, true))
	if !filesystem.Exists(objects.GetFieldValueFromStruct(config, fieldsPrefix+fields[0], runtime.F.GetCallerInfo(runtime.F{}, true))) {
		valid = false
	}
	return valid
}

func (f F) SetKubeConfig() {
	if f.validKubeConfig() {
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfra}, "Using infrastructure provider: selfhosted (kubeconfig)\n", 0)
	} else {
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiInfra}, "Invalid and/or missing kubeconfigs in infrastructure config.\nStarting assistant to retrieve missing settings..\n", 0)
		f.setKubeConfigAuth()
	}
	vars.InfraProviderActive = vars.InfraProviderSelfHosted
}

func (f F) setKubeConfigAuth() {
	logging.Log([]string{"", vars.EmojiSettings, vars.EmojiKubernetes}, "Please configure the kubernetes cluster authentication (kubeconfig).", 0)
	locHomeDir := terminal.GetUserSelection("Do you want to import your kubeconfig file from your OS user home directory ("+filesystem.GetUserHomeDir()+")?", []string{}, false, true)
	var path string
	if locHomeDir == "Yes" {
		path = filesystem.GetUserHomeDir() + "/.kube/"
	} else {
		path = f.getProjectAuthPath()
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiKubernetes}, "Okay, then please manually copy your kubeconfig file to your project located at: "+path, 0)
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiKubernetes}, "If you have inserted the kubeconfig file at this location you can select it here:", 0)
	}
	projectAuthPath := terminal.UserSelectionFiles("Which kubeconfig file should be the default for your projects?", "files", path, []string{}, []string{"hidden"}, false, true)
	if locHomeDir == "Yes" {
		fileName := filesystem.GetFileNameFromDirPath(path)
		projectAuthPath = f.getProjectAuthPath() + fileName
		filesystem.Copy(path, projectAuthPath, false)
	}
	f.setValue("Spec.Cluster.Auth.KubeConfigPath", projectAuthPath)
}

func (f F) getProjectAuthPath() string {
	return users.GetIDActive() + vars.ProjectInfrastructureAuthPath
}

func (f F) setValue(key string, value string) {
	yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, key, value, runtime.F.GetCallerInfo(runtime.F{}, true)))
}

func (f F) getValue(key string) string {
	return objects.GetFieldValueFromStruct(config, key, runtime.F.GetCallerInfo(runtime.F{}, true))
}
