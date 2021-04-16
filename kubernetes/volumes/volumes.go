package volumes

import (
	contextPack "context"
	"fmt"
	"strconv"

	"../../errors"
	"../../logging"
	"../../objects/strings"
	"../../runtime"
	"../../vars"
	"../client"
	"../namespaces"
	"./templates"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type F struct{}

func (f F) Get(namespace, volumeType string, logResult bool) []string {
	client := f.getClient()
	var volumes []string
	if volumeType == "pv" {
		list, err := client.PersistentVolumes().List(contextPack.TODO(), metav1.ListOptions{}) // TODO: Test namespace selection
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+volumeType+" in namespace "+namespaces.Default+"!", false, true, true) {
			if logResult {
				logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Volume "+volumeType+"s:", 0)
				if len(list.Items) == 0 {
					logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
				}
			}
			for i, pv := range list.Items {
				if logResult {
					fmt.Printf(" [" + strconv.Itoa(i) + "] " + pv.Name)
					fmt.Printf(" | " + pv.Spec.PersistentVolumeSource.HostPath.Path)
					fmt.Printf(" | " + pv.Spec.StorageClassName)
					//fmt.Printf(" | "+pv.Spec.Capacity["storage"].String())
					fmt.Printf("\n")
				}
				volumes = append(volumes, pv.Name)
			}
		}
	} else if volumeType == "pvc" {
		list, err := client.PersistentVolumeClaims(namespace).List(contextPack.TODO(), metav1.ListOptions{})
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+volumeType+" in namespace "+namespaces.Default+"!", false, true, true) {
			if logResult {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "Volume "+volumeType+"s:", 0)
				if len(list.Items) == 0 {
					logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
				}
			}
			for i, pvc := range list.Items {
				if logResult {
					fmt.Printf(" [" + strconv.Itoa(i) + "] " + pvc.Name)
					fmt.Printf(" | " + string(*pvc.Spec.VolumeMode))
					fmt.Printf("\n")
				}
				volumes = append(volumes, pvc.Name)
			}
		}
	}
	return volumes
}

func (f F) Create(namespace, contextID, serviceCluster string, volumes [][]string) {
	client := f.getClient()
	for i, volume := range volumes {
		size := volume[1]
		diskType := "ssd"
		volumeType := "pvc"
		name := volumeType + "-" + contextID + "-" + strings.ToString(i+1)
		if !strings.ArrayContains(f.Get(namespace, volumeType, false), name) {
			_, err := client.PersistentVolumeClaims(namespace).Create(contextPack.TODO(), templates.GetConfigPVC(namespace, name, size), metav1.CreateOptions{})
			if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create "+volumeType+" "+name+"!", false, true, true) {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created "+volumeType+" "+name+".", 1)
			}
		}
		volumeType = "pv"
		name = volumeType + "-" + contextID + "-" + strings.ToString(i+1)
		if !strings.ArrayContains(f.Get(namespace, volumeType, false), name) {
			_, err := client.PersistentVolumes().Create(contextPack.TODO(), templates.GetConfigPV(namespace, name, size, diskType, serviceCluster), metav1.CreateOptions{})
			if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create "+volumeType+" "+name+"!", false, true, true) {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created "+volumeType+" "+name+".", 1)
			}
		}
	}
}

func (f F) Delete(namespace, contextID string, volumes [][]string) {
	client := f.getClient()
	for i, _ := range volumes {
		volumeType := "pvc"
		name := volumeType + "-" + contextID + "-" + strings.ToString(i+1)
		if strings.ArrayContains(f.Get(namespace, volumeType, false), name) {
			deletePolicy := metav1.DeletePropagationForeground
			err := client.PersistentVolumeClaims(namespace).Delete(contextPack.TODO(), name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			})
			if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete "+volumeType+" "+name+"!", false, true, true) {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Deleted "+volumeType+" "+name+".", 1)
			}
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), ""+volumeType+" "+name+" already deleted.", true, false, true)
		}
		volumeType = "pv"
		name = volumeType + "-" + contextID + "-" + strings.ToString(i+1)
		if strings.ArrayContains(f.Get(namespace, volumeType, false), name) {
			deletePolicy := metav1.DeletePropagationForeground
			err := client.PersistentVolumes().Delete(contextPack.TODO(), name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			})
			if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete "+volumeType+" "+name+"!", false, true, true) {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Deleted "+volumeType+" "+name+".", 1)
			}
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), ""+volumeType+" "+name+" already deleted.", true, false, true)
		}
	}
}

func (f F) Exists(namespace, contextID string) bool {
	contextID = contextID + "-1"
	if !strings.ArrayContains(f.Get(namespace, "pvc", false), "pvc-"+contextID) {
		return false
	} else if !strings.ArrayContains(f.Get(namespace, "pv", false), "pv-"+contextID) {
		return false
	}
	return true
}

func (f F) getClient() corev1typed.CoreV1Interface {
	return client.F.GetAuth(client.F{}).CoreV1()
}
