package templates

type DaemonSetTemplate struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Labels    struct {
			K8SApp string `json:"k8s-app"`
		} `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Selector struct {
			MatchLabels struct {
				K8SApp string `json:"k8s-app"`
			} `json:"matchLabels"`
		} `json:"selector"`
		UpdateStrategy struct {
			Type string `json:"type"`
		} `json:"updateStrategy"`
		Template struct {
			Metadata struct {
				Labels struct {
					Name   string `json:"name"`
					K8SApp string `json:"k8s-app"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Affinity struct {
					NodeAffinity struct {
						RequiredDuringSchedulingIgnoredDuringExecution struct {
							NodeSelectorTerms []struct {
								MatchExpressions []struct {
									Key      string `yaml:vars.EmojiCrypto`
									Operator string `json:"operator"`
								} `json:"matchExpressions"`
							} `json:"nodeSelectorTerms"`
						} `json:"requiredDuringSchedulingIgnoredDuringExecution"`
					} `json:"nodeAffinity"`
				} `json:"affinity"`
				Tolerations []struct {
					Operator string `json:"operator"`
				} `json:"tolerations"`
				HostNetwork bool `json:"hostNetwork"`
				HostPID     bool `json:"hostPID"`
				Volumes     []struct {
					Name     string `json:"name"`
					HostPath struct {
						Path string `json:"path"`
					} `json:"hostPath"`
				} `json:"volumes"`
				InitContainers []struct {
					Image           string `json:"image"`
					ImagePullPolicy string `json:"imagePullPolicy"`
					Name            string `json:"name"`
					Resources       struct {
						Requests struct {
							CPU string `json:"cpu"`
						} `json:"requests"`
					} `json:"resources"`
					SecurityContext struct {
						Privileged bool `json:"privileged"`
					} `json:"securityContext"`
					Env []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"env"`
					VolumeMounts []struct {
						Name      string `json:"name"`
						MountPath string `json:"mountPath"`
					} `json:"volumeMounts"`
				} `json:"initContainers"`
				Containers []struct {
					Image string `json:"image"`
					Name  string `json:"name"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}
