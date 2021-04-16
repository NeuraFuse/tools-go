package user

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
	} `json:"metadata"`
	Spec       struct {
		Auth       struct {
			JWT    struct {
				SigningKey string `json:"signingKey"`
			} `json:"jwt"`
		} `json:"auth"`
	} `json:"spec"`
}