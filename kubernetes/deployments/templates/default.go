package templates

import (
	"../../../config"
	"../../../objects/strings"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfig(id, imageAddrs string, serviceCluster, resourcesSpec string, volumes, containerPort [][]string) *appsv1.Deployment {
	ports := getContainerPorts(containerPort)
	resources, annotations := getResources(resourcesSpec)
	volumeMounts, volumesAPIv1 := getVolumes(id, volumes)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": id,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": id,
					},
					Annotations: annotations,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:         id,
							Image:        imageAddrs,
							VolumeMounts: volumeMounts,
							Ports:        ports,
							Resources:    resources,
						},
					},
					Volumes: volumesAPIv1,
					NodeSelector: map[string]string{
						//NodeSelectorTerms: nodeSelectorTerm(serviceCluster),
						"service-cluster": serviceCluster,
					},
				},
			},
		},
	}
	return deployment
}

func getContainerPorts(ports [][]string) []apiv1.ContainerPort {
	var containerPorts []apiv1.ContainerPort
	for i, port := range ports {
		var containerPort apiv1.ContainerPort
		containerPort.Name = "port-" + strings.ToString(i+1)
		containerPort.ContainerPort = strings.ToInt32(port[0])
		if port[1] == "TCP" {
			containerPort.Protocol = apiv1.ProtocolTCP
		} else {
			containerPort.Protocol = apiv1.ProtocolUDP
		}
		containerPorts = append(containerPorts, containerPort)
	}
	return containerPorts
}

func getVolumes(id string, volumes [][]string) ([]apiv1.VolumeMount, []apiv1.Volume) {
	pvName := "pv-" + id
	pvcName := "pvc-" + id
	if strings.Contains(id, "inference") { // TODO: Refactor
		pvName = "pv-app-1"
		pvcName = "pvc-app-1"
		volumes = [][]string{{"/app/lightning/pytorch/data", "100Gi"}}
	}
	var volumeMounts []apiv1.VolumeMount
	var volumesAPIv1 []apiv1.Volume
	for i, volume := range volumes {
		var volumeMount apiv1.VolumeMount
		pvID := pvName + "-" + strings.ToString(i+1)
		volumeMount.Name = pvID
		volumeMount.MountPath = volume[0]
		volumeMounts = append(volumeMounts, volumeMount)

		var volumeAPIv1 apiv1.Volume
		volumeAPIv1.Name = pvID
		var volumeSource apiv1.VolumeSource
		pvcID := pvcName + "-" + strings.ToString(i+1)
		volumeSource.PersistentVolumeClaim = &apiv1.PersistentVolumeClaimVolumeSource{ClaimName: pvcID} //pvcVolumeSource
		volumeAPIv1.VolumeSource = volumeSource
		volumesAPIv1 = append(volumesAPIv1, volumeAPIv1)
	}
	return volumeMounts, volumesAPIv1
}

func getResources(res string) (apiv1.ResourceRequirements, map[string]string) {
	var resQuantity string = "0"
	var resName string = ""
	const GPUKey string = "nvidia.com/gpu"
	const TPUKey string = "cloud-tpus.google.com/"
	resRequested := false
	annotations := make(map[string]string)
	if res == "gpu" {
		resName = GPUKey
		resQuantity = "1"
		resRequested = true
	} else if res == "tpu" {
		tpuVersion := config.Setting("get", "infrastructure", "Spec.Gcloud.Accelerator.TPU.Version", "")
		resName = TPUKey + tpuVersion
		resQuantity = config.Setting("get", "infrastructure", "Spec.Gcloud.Accelerator.TPU.Cores", "")
		resRequested = true
		annotations["tf-version.cloud-tpus.google.com"] = config.Setting("get", "infrastructure", "Spec.Gcloud.Accelerator.TPU.TF.Version", "")
	}
	var resList apiv1.ResourceList
	if resRequested {
		resList = apiv1.ResourceList{apiv1.ResourceName(resName): resource.MustParse(resQuantity)}
	} else {
		resList = apiv1.ResourceList{}
	}
	var resReq apiv1.ResourceRequirements
	resReq.Limits = resList
	return resReq, annotations
}

func nodeSelectorTerm(serviceCluster string) []apiv1.NodeSelectorTerm {
	nodeSelectorTerm := []apiv1.NodeSelectorTerm{apiv1.NodeSelectorTerm{}}
	nodeSelectorTerm[0].MatchExpressions = nodeSelectorRequirement(serviceCluster)
	return nodeSelectorTerm
}

func nodeSelectorRequirement(serviceCluster string) []apiv1.NodeSelectorRequirement {
	nodeSelectorRequirement := []apiv1.NodeSelectorRequirement{apiv1.NodeSelectorRequirement{}}
	nodeSelectorRequirement[0].Key = "service-cluster"
	nodeSelectorRequirement[0].Values = []string{serviceCluster}
	return nodeSelectorRequirement
}

func int32Ptr(i int32) *int32 { return &i }
