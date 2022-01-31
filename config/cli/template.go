package cli

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
	} `json:"metadata"`
	Spec       struct {
		Users struct {
			DefaultID string `json:"defaultID"`
			ActiveID string `json:"activeID"`
		} `json:"users"`
		Updates struct {
			Auto struct {
				Status string `json:"status"`
			} `json:"auto"`
		} `json:"updates"`
	} `json:"spec"`
}