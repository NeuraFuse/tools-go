package templates

type Default struct {
	Version string `json:"version"`
	Images  struct {
		Remote struct {
			Image string `json:"image"`
		} `json:"remote"`
	} `json:"images"`
	Dev struct {
		Sync []struct {
			ImageName     string   `json:"imageName"`
			Namespace     string   `json:"namespace"`
			LocalSubPath  string   `json:"localSubPath"`
			ContainerPath string   `json:"containerPath"`
			ExcludePaths  []string `json:"excludePaths"`
		} `json:"sync"`
	} `json:"dev"`
}
