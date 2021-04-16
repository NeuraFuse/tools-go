package project

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
	} `json:"metadata"`
	Spec       struct {
		WorkingDir string `json:"workingDir"`
		Containers struct {
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