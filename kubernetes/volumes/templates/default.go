package templates

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func GetConfigPVC(namespace, id, size string) *apiv1.PersistentVolumeClaim {
	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
			Namespace: namespace,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: apiv1.ResourceRequirements {
				Limits: apiv1.ResourceList{
					"storage": resource.MustParse(size),
				},
				Requests: apiv1.ResourceList{
					"storage": resource.MustParse(size),
				},
			},
			//VolumeName: id,
			//StorageClassName: "fast",
			//VolumeMode: apiv1.PersistentVolumeMode{apiv1.PersistentVolumeFilesystem},
		},
	}
	return pvc
}

func GetConfigPV(namespace, id, size, diskType, serviceCluster string) *apiv1.PersistentVolume {
	pv := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
			Namespace: namespace,
		},
		Spec: apiv1.PersistentVolumeSpec{
			Capacity: apiv1.ResourceList{
				"storage": resource.MustParse(size),
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource {
				/*GCEPersistentDisk: &apiv1.GCEPersistentDiskVolumeSource {
					PDName: "pd-"+id,
					FSType: "ext4",
				},
				Local: &apiv1.LocalVolumeSource {
					Path: "/data",
					//FSType: apiv1.String("ext4"),
				},*/
				HostPath: &apiv1.HostPathVolumeSource {
					Path: "/data",
					//Type: ,
				},
			},
			AccessModes: []apiv1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			StorageClassName: "fast",
			NodeAffinity: &apiv1.VolumeNodeAffinity {
				Required: &apiv1.NodeSelector {
					NodeSelectorTerms: []apiv1.NodeSelectorTerm {
						{
							MatchExpressions: []apiv1.NodeSelectorRequirement {
								{
									Key: "service-cluster",
									Operator: "In",
									Values: []string{serviceCluster},
								},
							},
						},
					},
				},
			},
		},
	}
	return pv
}