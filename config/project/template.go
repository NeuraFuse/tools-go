package project

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
	} `json:"metadata"`
	Spec       struct {
		Infrastructure struct {
			Cluster struct {
				ID string `json:"id"`
			} `json:"cluster"`
		} `json:"infrastructure"`
		App struct {
			Kind string `json:"kind"`
		} `json:"app"`
		Containers struct {
			Network struct {
				Mode string `json:"mode"`
			} `json:"network"`
			Sync struct {
				RemoteURL string `json:"remoteURL"`
				PathMappings struct {
					LocalAppRoot string `json:"localAppRoot"`
					LocalIDERoot string `json:"localIDERoot"`
					ContainerAppRoot string `json:"containerAppRoot"`
					ContainerAppDataRoot string `json:"containerAppDataRoot"`
				} `json:"pathMappings"`
			} `json:"sync"`
			Registry struct {
				Address string `json:"address"`
				Auth        struct {
					Username     string `json:"username"`
					Password     string `json:"password"`
				} `json:"auth"`
			} `json:"registry"`
		} `json:"containers"`
	} `json:"spec"`
}