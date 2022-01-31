package dev

type Default struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
	} `json:"metadata"`
	Spec       struct {
		Status   string `json:"status"`
		LogLevel string `json:"logLevel"`
		API      struct {
			Address string `json:"address"`
		} `json:"api"`
		CI      struct {
			Mode string `json:"mode"`
		} `json:"ci"`
		Containers struct {
			Registry struct {
				Address string `json:"address"`
				Auth        struct {
					Username     string `json:"username"`
					Password     string `json:"password"`
				} `json:"auth"`
			} `json:"registry"`
		} `json:"containers"`
		Build struct {
			Neuracli struct {
				Version    string `json:"version"`
				IDRecent   string `json:"idRecent"`
				HashRecent string `json:"hashRecent"`
			} `json:"neuracli"`
			Neurakube struct {
				Version    string `json:"version"`
				IDRecent   string `json:"idRecent"`
				HashRecent string `json:"hashRecent"`
				Container  struct {
					HashRecent string `json:"hashRecent"`
				} `json:"container"`
			} `json:"neurakube"`
			Accelerator struct {
				Base struct {
					Pytorch struct {
						GPU struct {
							HashRecent string `json:"hashRecent"`
						} `json:"gpu"`
						TPU struct {
							HashRecent string `json:"hashRecent"`
						} `json:"tpu"`
					} `json:"pytorch"`
				} `json:"base"`
				Remote struct {
					Pytorch struct {
						GPU struct {
							HashRecent string `json:"hashRecent"`
							Lightning  struct {
								HashRecent string `json:"hashRecent"`
							} `json:"lightning"`
						} `json:"gpu"`
						TPU struct {
							HashRecent string `json:"hashRecent"`
							Lightning  struct {
								HashRecent string `json:"hashRecent"`
							} `json:"lightning"`
						} `json:"tpu"`
					} `json:"pytorch"`
				} `json:"remote"`
				App struct {
					Pytorch struct {
						GPU struct {
							HashRecent string `json:"hashRecent"`
							Lightning  struct {
								HashRecent string `json:"hashRecent"`
							} `json:"lightning"`
						} `json:"gpu"`
						TPU struct {
							HashRecent string `json:"hashRecent"`
							Lightning  struct {
								HashRecent string `json:"hashRecent"`
							} `json:"lightning"`
						} `json:"tpu"`
					} `json:"pytorch"`
				} `json:"app"`
				Inference struct {
					Pytorch struct {
						GPU struct {
							HashRecent string `json:"hashRecent"`
							Lightning  struct {
								HashRecent string `json:"hashRecent"`
							} `json:"lightning"`
						} `json:"gpu"`
						TPU struct {
							HashRecent string `json:"hashRecent"`
							Lightning  struct {
								HashRecent string `json:"hashRecent"`
							} `json:"lightning"`
						} `json:"tpu"`
					} `json:"pytorch"`
				} `json:"inference"`
			} `json:"accelerator"`
		} `json:"build"`
	} `json:"spec"`
}