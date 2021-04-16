package server

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
	} `json:"metadata"`
	Spec       struct {
		Users struct {
			Admin string `json:"admin"`
		}
	} `json:"spec"`
}