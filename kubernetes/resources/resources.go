package resources

import (
	"os"

	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/kubernetes/daemonsets"
	"github.com/neurafuse/tools-go/kubernetes/daemonsets/templates"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/vars"
)

func Check(context, resourceType string) {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInspect}, "Checking resourceType "+resourceType+" drivers/environment..", 0)
	if resourceType == "gpu" {
		drivers := []string{"nvidia-driver-installer"}
		var createdNew bool
		for _, driverID := range drivers {
			if !daemonsets.F.Exists(daemonsets.F{}, driverID) {
				createdNew = true
				logging.Log([]string{"", vars.EmojiInspect, vars.EmojiWarning}, "The driver (daemonset) "+driverID+" does not exist yet.", 0)
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Starting creation..", 0)
				ds := templates.GetConfig("nvidia", infraConfig.F.GetProviderIDActive(infraConfig.F{}), "nvidia-driver-installer", "cos")
				daemonsets.F.Create(daemonsets.F{}, ds)
			}
		}
		var status string
		if createdNew {
			status = "now"
		} else {
			status = "already"
		}
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "All necessary "+resourceType+" drivers (daemonsets) are "+status+" created.\n", 0)
	} else if resourceType == "tpu" && context == "container" {
		envKey := "XRT_TPU_CONFIG"
		if os.Getenv(envKey) == "" {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Creating "+resourceType+" environment variable "+envKey+"..", 0)
			tpuEndpoints := os.Getenv("KUBE_GOOGLE_CLOUD_TPU_ENDPOINTS")
			tpuEndpoints = strings.TrimPrefix(tpuEndpoints, "grpc://")
			envValue := "tpu_worker;0;" + tpuEndpoints
			os.Setenv(envKey, envValue)
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created "+resourceType+" environment variable "+envKey+"=\""+envValue+"\" .", 0)
		}
		envKey = "TPU_SPLIT_COMPILE_AND_EXECUTE"
		if os.Getenv(envKey) == "" {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Creating "+resourceType+" environment variable "+envKey+"..", 0)
			envValue := "1"
			os.Setenv(envKey, envValue)
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created "+resourceType+" environment variable "+envKey+"=\""+envValue+"\" .", 0)
		}
	}
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, resourceType+" drivers/environment are now ready.\n", 0)
}
