package config

import (
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/runtime"

	//"github.com/neurafuse/tools-go/objects/strings"
	cliConfig "github.com/neurafuse/tools-go/config/cli"
	devConfig "github.com/neurafuse/tools-go/config/dev"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	projectConfig "github.com/neurafuse/tools-go/config/project"
	serverConfig "github.com/neurafuse/tools-go/config/server"
	userConfig "github.com/neurafuse/tools-go/config/user"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/yaml"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

func Setting(action, configID, configKey string, configValue string) string {
	if usersID.F.ActiveIsSet(usersID.F{}) {
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiDev}, action+" "+configID+" "+configKey+" "+configValue+"\n", 2)
	}
	var config interface{}
	var packageName string = configID
	var filePath string
	if configID == "cli" {
		cliConfig.F.SetConfig(cliConfig.F{})
		config = cliConfig.F.GetConfig(cliConfig.F{})
		filePath = cliConfig.F.GetFilePath(cliConfig.F{})
	} else if configID == "dev" {
		devConfig.F.SetConfig(devConfig.F{})
		config = devConfig.F.GetConfig(devConfig.F{})
		filePath = devConfig.F.GetFilePath(devConfig.F{})
	} else if configID == "user" {
		userConfig.F.SetConfig(userConfig.F{})
		config = userConfig.F.GetConfig(userConfig.F{})
		filePath = userConfig.F.GetFilePath(userConfig.F{})
	} else if configID == "infrastructure" {
		infraConfig.F.SetConfig(infraConfig.F{})
		config = infraConfig.F.GetConfig(infraConfig.F{})
		filePath = infraConfig.F.GetPath(infraConfig.F{}, true)
	} else if configID == "server" {
		serverConfig.F.SetConfig(serverConfig.F{})
		config = serverConfig.F.GetConfig(serverConfig.F{})
		filePath = serverConfig.F.GetFilePath(serverConfig.F{})
	} else if configID == "project" {
		projectConfig.F.SetConfig(projectConfig.F{})
		config = projectConfig.F.GetConfig(projectConfig.F{})
		filePath = projectConfig.F.GetFilePath(projectConfig.F{})
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid configID "+configID+"!", true, true, true)
	}
	var result string
	if action == "set" {
		setValue(config, filePath, configKey, configValue, packageName)
	} else if action == "get" {
		result = getValue(config, configKey, packageName)
	} else if action == "reset" {
		yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, configKey, "", packageName))
	} else if action == "init" {
	}
	return result
}

func ValidSettings(configID, configType string, logInvalid bool) bool {
	var valid bool = true
	var config interface{}
	var packageName string = configID
	var fields []string
	var fieldsPrefix string = ""
	if configID == "cli" {
		cliConfig.F.SetConfig(cliConfig.F{})
		config = cliConfig.F.GetConfig(cliConfig.F{})
		if configType == "users" {
			fieldsPrefix = "Spec.Users."
			fields = []string{"DefaultID"}
		} else if configType == "updates" {
			fieldsPrefix = "Spec.Updates."
			fields = []string{"Auto.Status"}
		}
	} else if configID == "user" {
		userConfig.F.SetConfig(userConfig.F{})
		config = userConfig.F.GetConfig(userConfig.F{})
		if configType == "login" {
			fieldsPrefix = "Spec."
			fields = []string{"Auth.JWT.SigningKey"}
		}
		if configType == "defaults/infra" {
			fieldsPrefix = "Spec.Defaults.Infrastructure."
			fields = []string{"ID"}
		}
	} else if configID == "infrastructure" {
		infraConfig.F.SetConfig(infraConfig.F{})
		config = infraConfig.F.GetConfig(infraConfig.F{})
		if configType == "cluster" {
			fieldsPrefix = "Spec.Cluster."
			fields = []string{"ID", "SelfDeletion.Active", "Nodes.DiskSizeGb"}
		} else if configType == vars.InfraProviderGcloud {
			fieldsPrefix = "Spec.Gcloud."
			fields = []string{"ProjectID", "Zone", "MachineType"}
			var filePath string = infraConfig.F.GetInfraGcloudAuthPath(infraConfig.F{})
			if !filesystem.Exists(filePath) {
				valid = false
			}
		} else if configType == vars.InfraProviderGcloud+"/accelerator/tpu" {
			fieldsPrefix = "Spec.Gcloud.Accelerator.TPU."
			fields = []string{"Version", "Cores", "MachineType", "TF.Version", "Node.DiskSizeGb"}
		} else if configType == vars.InfraProviderGcloud+"/accelerator/gpu" {
			fieldsPrefix = "Spec.Gcloud.Accelerator.GPU."
			fields = []string{"MachineType", "Type", "Node.DiskSizeGb"}
		} else if configType == "remote" || configType == "app" || configType == "inference" {
			fieldsPrefix = "Spec." + strings.Title(configType) + "."
			fields = []string{"Type", "Environment.IDE", "Environment.Framework", "VolumeSizeGB", "SelfDeletion.Active", "NodePools.Dedicated"}
		} else if configType == vars.NeuraKubeNameID {
			fieldsPrefix = "Spec.Neurakube."
			fields = []string{"VolumeSizeGB"}
		}
	} else if configID == "server" {
		serverConfig.F.SetConfig(serverConfig.F{})
		config = serverConfig.F.GetConfig(serverConfig.F{})
		if configType == "useradmin" {
			fieldsPrefix = "Spec.Users."
			fields = []string{"Admin"}
		}
	} else if configID == "project" {
		projectConfig.F.SetConfig(projectConfig.F{})
		config = projectConfig.F.GetConfig(projectConfig.F{})
		if configType == "containers/registry" {
			fieldsPrefix = "Spec.Containers.Registry."
			fields = []string{"Address", "Auth.Username", "Auth.Password"}
		} else if configType == "containers/sync" {
			fieldsPrefix = "Spec.Containers.Sync.PathMappings."
			fields = []string{"LocalAppRoot", "LocalIDERoot", "ContainerAppRoot", "ContainerAppDataRoot"}
		} else if configType == "app" {
			fieldsPrefix = "Spec.App."
			fields = []string{"Kind"}
		} else if configType == "infra/cluster" {
			fieldsPrefix = "Spec.Infrastructure.Cluster."
			fields = []string{"ID"}
		}
	} else if configID == "dev" {
		devConfig.F.SetConfig(devConfig.F{})
		config = devConfig.F.GetConfig(devConfig.F{})
		if configType == "containers" {
			fieldsPrefix = "Spec.Containers.Registry."
			fields = []string{"Address"}
		} else if configType == "api" {
			fieldsPrefix = "Spec.API."
			fields = []string{"Address"}
		} else {
			fieldsPrefix = "Spec."
			fields = []string{"Status", "Name", "WorkingDir", "LogLevel"}
		}
	}
	if !valid {
		return false
	} else {
		valid = objects.StructFieldValuesExisting(config, fieldsPrefix, fields, packageName)
		if !valid {
			if logInvalid {
				logging.Log([]string{"", vars.EmojiSettings, vars.EmojiWarning}, "Incomplete or invalid "+configType+" settings ("+configID+".yaml)!", 0)
				logging.Log([]string{"", vars.EmojiSettings, vars.EmojiAssistant}, "Starting assistant..", 0)
			}
		}
		return valid
	}
}

func APILocationCluster() bool {
	var cluster bool
	if DevConfigActive() {
		var configKind string = "dev"
		var configKey string = "Spec.API.Address"
		var apiAddrs string = Setting("get", configKind, configKey, "")
		if apiAddrs == "cluster" {
			cluster = true
		} else if apiAddrs == "local" {
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiAPI}, "API location: "+apiAddrs, 1)
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid setting for "+configKind+" config key: "+configKey, true, true, true)
		}
	} else {
		cluster = true
	}
	return cluster
}

func DevConfigActive() bool {
	var active bool
	if Setting("get", "dev", "Spec.Status", "") == "active" {
		active = true
	}
	return active
}

func GetResourceTypes() []string {
	var users string = "Users"
	var projects string = "Projects"
	var infra string = "Infrastructures"
	var types []string = []string{users, projects, infra}
	if DevConfigActive() {
		var dev string = "Development"
		types = append(types, dev)
	}
	return types
}

func setValue(config interface{}, filePath, key, value, packageName string) {
	yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, key, value, packageName))
}

func getValue(config interface{}, key, packageName string) string {
	return objects.GetFieldValueFromStruct(config, key, packageName)
}
