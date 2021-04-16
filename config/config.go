package config

import (
	"../../tools-go/errors"
	"../../tools-go/filesystem"
	"../../tools-go/logging"
	"../../tools-go/objects"
	"../../tools-go/runtime"

	//"../../tools-go/objects/strings"
	"../../tools-go/objects/strings"
	"../../tools-go/readers/yaml"
	"../../tools-go/vars"
	cli "./cli"
	dev "./dev"
	infra "./infrastructure"
	project "./project"
	server "./server"
	user "./user"
)

func Setting(action, configID, key string, value string) string {
	/*success, refVal := objects.CallStructInterfaceFuncByName(Packages{}, strings.Title(strings.TrimSuffix(configID, "config")), "SetConfig")
	success, filePathTest := objects.CallStructInterfaceFuncByName(Packages{}, strings.Title(strings.TrimSuffix(configID, "config")), "GetFilePath")
	if success {
		//config, filePath := refVal
		fmt.Println(objects.GetInterfaceFromReflectValue(refVal[0], refVal[0])
		fmt.Println(objects.GetStringFromReflectValue(refVal[1]))
	}*/
	var config interface{}
	packageName := configID
	filePath := ""
	if configID == "cli" {
		cli.F.SetConfig(cli.F{})
		config = cli.F.GetConfig(cli.F{})
		filePath = cli.F.GetFilePath(cli.F{})
	} else if configID == "dev" {
		dev.F.SetConfig(dev.F{})
		config = dev.F.GetConfig(dev.F{})
		filePath = dev.F.GetFilePath(dev.F{})
	} else if configID == "user" {
		user.F.SetConfig(user.F{})
		config = user.F.GetConfig(user.F{})
		filePath = user.F.GetFilePath(user.F{})
	} else if configID == "infrastructure" {
		infra.F.SetConfig(infra.F{})
		config = infra.F.GetConfig(infra.F{})
		filePath = infra.F.GetFilePath(infra.F{})
	} else if configID == "server" {
		server.F.SetConfig(server.F{})
		config = server.F.GetConfig(server.F{})
		filePath = server.F.GetFilePath(server.F{})
	} else if configID == "project" {
		project.F.SetConfig(project.F{})
		config = project.F.GetConfig(project.F{})
		filePath = project.F.GetFilePath(project.F{})
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid configID "+configID+"!", true, true, true)
	}
	result := ""
	if action == "set" {
		setValue(config, filePath, key, value, packageName)
	} else if action == "get" {
		result = getValue(config, key, packageName)
	} else if action == "reset" {
		yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, key, "", packageName))
	} else if action == "init" {
	}
	return result
}

func ValidSettings(configID, configType string, logInvalid bool) bool {
	var valid bool = true
	var config interface{}
	packageName := configID
	var fields []string
	fieldsPrefix := ""
	if configID == "cli" {
		cli.F.SetConfig(cli.F{})
		config = cli.F.GetConfig(cli.F{})
		if configType == "users" {
			fieldsPrefix = "Spec.Users."
			fields = []string{"DefaultID"}
		} else if configType == "projects" {
			fieldsPrefix = "Spec.Projects."
			fields = []string{"DefaultID"}
		} else if configType == "updates" {
			fieldsPrefix = "Spec.Updates."
			fields = []string{"Auto.Status"}
		}
	} else if configID == "user" {
		user.F.SetConfig(user.F{})
		config = user.F.GetConfig(user.F{})
		if configType == "login" {
			fieldsPrefix = "Spec."
			fields = []string{"Auth.JWT.SigningKey"}
		}
	} else if configID == "infrastructure" {
		infra.F.SetConfig(infra.F{})
		config = infra.F.GetConfig(infra.F{})
		if configType == "cluster" {
			fieldsPrefix = "Spec.Cluster."
			fields = []string{"Name", "Auth.Password", "SelfDeletion.Active", "Nodes.DiskSizeGb"}
		} else if configType == "kube" {
			fieldsPrefix = "Spec.Cluster.Auth."
			fields = []string{"KubeConfigPath"}
			if !filesystem.Exists(objects.GetFieldValueFromStruct(config, fieldsPrefix+fields[0], runtime.F.GetCallerInfo(runtime.F{}, true))) {
				valid = false
			}
		} else if configType == vars.InfraProviderGcloud {
			fieldsPrefix = "Spec.Gcloud."
			fields = []string{"Auth.ServiceAccountJSONPath", "ProjectID", "Zone", "MachineType"}
			if !filesystem.Exists(objects.GetFieldValueFromStruct(config, fieldsPrefix+fields[0], runtime.F.GetCallerInfo(runtime.F{}, true))) {
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
		} else if configType == vars.NeuraKubeNameRepo {
			fieldsPrefix = "Spec.Neurakube."
			fields = []string{"VolumeSizeGB"}
		}
	} else if configID == "server" {
		server.F.SetConfig(server.F{})
		config = server.F.GetConfig(server.F{})
		if configType == "useradmin" {
			fieldsPrefix = "Spec.Users."
			fields = []string{"Admin"}
		}
	} else if configID == "project" {
		project.F.SetConfig(project.F{})
		config = project.F.GetConfig(project.F{})
		if configType == "containers" {
			fieldsPrefix = "Spec.Containers.Registry."
			fields = []string{"Address", "Auth.Username", "Auth.Password"}
		} else {
			fields = []string{"Metadata.ID", "Metadata.Name", "Spec.WorkingDir", ""}
		}
	} else if configID == "dev" {
		dev.F.SetConfig(dev.F{})
		config = dev.F.GetConfig(dev.F{})
		if configType == "containers" {
			fieldsPrefix = "Spec.Containers."
			fields = []string{"RepositoryAddrs", "Auth.Username", "Auth.Password"}
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
				logging.Log([]string{"\n", vars.EmojiSettings, vars.EmojiWarning}, "Incomplete or invalid "+configType+" settings ("+configID+"config)!", 0)
				logging.Log([]string{"", vars.EmojiSettings, vars.EmojiAssistant}, "Starting assistant..", 0)
			}
		}
		return valid
	}
}

func setValue(config interface{}, filePath, key, value, packageName string) {
	yaml.StructToFile(filePath, objects.SetFieldValueFromStruct(config, key, value, packageName))
}

func getValue(config interface{}, key, packageName string) string {
	return objects.GetFieldValueFromStruct(config, key, packageName)
}