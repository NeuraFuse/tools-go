package templates

type Default struct {
	Version        string `json:"version"`
	Configurations []struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Request      string `json:"request"`
		Port         int    `json:"port,omitempty"`
		Host         string `json:"host,omitempty"`
		Secret       string `json:"secret,omitempty"`
		PathMappings []struct {
			LocalRoot  string `json:"localRoot"`
			RemoteRoot string `json:"remoteRoot"`
		} `json:"pathMappings,omitempty"`
		Mode    string   `json:"mode,omitempty"`
		Program string   `json:"program,omitempty"`
		Args    []string `json:"args,omitempty"`
		Console string   `json:"console,omitempty"`
	} `json:"configurations"`
}