package infrastructure

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
	} `json:"metadata"`
	Spec       struct {
		Gcloud struct {
			Auth struct {
				ServiceAccountJSONPath string `json:"serviceAccountJSONPath"`
			} `json:"auth"`
			ProjectID              string `json:"projectID"`
			Zone                   string `json:"zone"`
			MachineType            string `json:"machineType"`
			Accelerator            struct {
				GPU struct {
					MachineType string `json:"machineType"`
					Type        string `json:"type"`
					Node        struct {
						DiskSizeGb string `json:"diskSizeGb"`
					} `json:"node"`
				} `json:"gpu"`
				TPU struct {
					MachineType string `json:"machineType"`
					Version     string `json:"version"`
					Cores       string `json:"cores"`
					TF          struct {
						Version string `json:"version"`
					} `json:"tf"`
					Node struct {
						DiskSizeGb string `json:"diskSizeGb"`
					} `json:"node"`
				} `json:"tpu"`
			} `json:"accelerator"`
		} `json:"gcloud"`
		Cluster struct {
			Name           string `json:"name"`
			Auth struct {
				Password       string `json:"password"`
				KubeConfigPath string `json:"kubeConfigPath"`
			} `json:"auth"`
			SelfDeletion   struct {
				Active            string `json:"active"`
				TimeDurationHours string `json:"timeDurationHours"`
			} `json:"selfDeletion"`
			Nodes struct {
				DiskSizeGb string `json:"diskSizeGb"`
			} `yaml:nodes`
		} `json:"cluster"`
		Neurakube struct {
			VolumeSizeGB string `json:"volumeSizeGB"`
			AutoCreation string `json:"autoCreation"`
			Cache        struct {
				Endpoint      string `json:"endpoint"`
				ConnectStatus string `json:"connectStatus"`
			} `json:"cache"`
		} `json:"neurakube"`
		Remote struct {
			Environment struct {
				IDE       string `json:"ide"`
				Framework string `json:"framework"`
			} `json:"environment"`
			Type         string `json:"type"`
			VolumeSizeGB string `json:"volumeSizeGB"`
			SelfDeletion struct {
				Active            string `json:"active"`
				TimeDurationHours string `json:"timeDurationHours"`
			} `json:"selfDeletion"`
			NodePools struct {
				Dedicated string `json:"dedicated"`
			} `json:"nodePools"`
		} `json:"remote"`
		App struct {
			Environment struct {
				IDE       string `json:"ide"`
				Framework string `json:"framework"`
			} `json:"environment"`
			Type         string `json:"type"`
			VolumeSizeGB string `json:"volumeSizeGB"`
			SelfDeletion struct {
				Active            string `json:"active"`
				TimeDurationHours string `json:"timeDurationHours"`
			} `json:"selfDeletion"`
			NodePools struct {
				Dedicated string `json:"dedicated"`
			} `json:"nodePools"`
		} `json:"app"`
		Inference struct {
			Environment struct {
				IDE       string `json:"ide"`
				Framework string `json:"framework"`
			} `json:"environment"`
			Type         string `json:"type"`
			VolumeSizeGB string `json:"volumeSizeGB"`
			SelfDeletion struct {
				Active            string `json:"active"`
				TimeDurationHours string `json:"timeDurationHours"`
			} `json:"selfDeletion"`
			NodePools struct {
				Dedicated string `json:"dedicated"`
			} `json:"nodePools"`
		} `json:"inference"`
	} `json:"spec"`
}