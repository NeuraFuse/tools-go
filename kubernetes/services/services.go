package services

import (
	contextPack "context"
	"fmt"
	"strconv"

	"../../errors"
	"../../logging"
	"../../objects/strings"
	"../../runtime"
	"../../timing"
	"../../vars"
	"../client"
	"../namespaces"
	"./templates"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type F struct{}

func (f F) getClient(namespace string) corev1typed.ServiceInterface {
	if namespace == "" {
		namespace = namespaces.Default
	}
	test := client.F.GetAuth(client.F{}).CoreV1().Services(namespace)
	return test
}

func (f F) getServiceList(namespace string) (*apiv1.ServiceList, error) {
	list, err := f.getClient(namespace).List(contextPack.TODO(), metav1.ListOptions{})
	return list, err
}

func (f F) Get(namespace, contextID string) *apiv1.Service {
	service, err := f.getClient(namespace).Get(contextPack.TODO(), contextID, metav1.GetOptions{})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get service "+contextID+" in namespace "+namespace+"!", false, true, true)
	return service
}

func (f F) GetList(namespace string, logResult bool) []string {
	list, err := f.getServiceList(namespace)
	var services []string
	if logResult {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, strings.Title(runtime.F.GetCallerInfo(runtime.F{}, true))+":", 0)
	}
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to list services in namespace "+namespace+":", false, true, true) {
		if len(list.Items) == 0 {
			if logResult {
				logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
			}
		}
		for i, s := range list.Items {
			if logResult {
				fmt.Printf("[" + strconv.Itoa(i) + "] " + s.Name)
				fmt.Printf(" | " + s.Spec.ClusterIP)
				fmt.Printf(" | " + string(s.Spec.Type))
				for iP, port := range s.Spec.Ports {
					fmt.Printf(" | ")
					fmt.Printf(string(port.Protocol) + ":" + strconv.Itoa(int(port.Port)))
					if iP > 0 && !(len(s.Spec.Ports) == iP+1) {
						fmt.Printf(", ")
					}
				}
				if len(s.Spec.ClusterIPs) != 0 {
					fmt.Printf(" | " + strings.Join(s.Spec.ClusterIPs, ","))
				}
				if len(s.Spec.ExternalIPs) != 0 {
					fmt.Printf(" | " + strings.Join(s.Spec.ExternalIPs, ","))
				}
				fmt.Printf("\n")
			}
			services = append(services, s.Name)
		}
	}
	return services
}

func (f F) GetLoadBalancerIP(namespace, contextID string) string {
	var lbIP string
	var success bool = false
	logging.ProgressSpinner("start")
	var loggedWaiting bool = false
	list, err := f.getServiceList(namespace)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to fetch services in namespace "+namespaces.Default+":", false, true, true) {
		if len(list.Items) != 0 {
			for i, s := range list.Items {
				if list.Items[i].ObjectMeta.Name == contextID {
					success = true
					if len(s.Status.LoadBalancer.Ingress) == 0 {
						if !loggedWaiting {
							logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for creation of "+contextID+" service endpoints..", 0)
							loggedWaiting = true
						}
						timing.TimeOut(1, "s")
					} else {
						logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Resolved endpoints for service "+contextID+".", 2)
						lbIP = s.Status.LoadBalancer.Ingress[0].IP
						success = true
						logging.ProgressSpinner("stop")
					}
				}
			}
		}
		if !success {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to find service "+contextID+"!", true, false, true)
		}
	}
	return lbIP
}

func (f F) GetClusterIP(namespace, contextID string) string {
	service := f.Get(namespace, contextID)
	return service.Spec.ClusterIP
}

func (f F) Create(namespace, contextID, clusterIP string, ports [][]string) {
	if !f.Exists(namespace, contextID) {
		service := templates.GetConfig(contextID, clusterIP, ports)
		_, err := f.getClient(namespace).Create(contextPack.TODO(), service, metav1.CreateOptions{})
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create service "+contextID+"!", false, false, true) {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created service "+contextID+".", 0)
		}
	} else {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "Service "+contextID+" already exists.", 0)
	}
}

func (f F) Delete(namespace, contextID string) {
	if f.Exists(namespace, contextID) {
		deletePolicy := metav1.DeletePropagationForeground
		if err := f.getClient(namespace).Delete(contextPack.TODO(), contextID, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete service!", true, true, true)
		} else {
			var success bool = false
			logging.ProgressSpinner("start")
			for ok := true; ok; ok = !success {
				if f.Exists(namespace, contextID) {
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for the service "+contextID+" to be deleted..", 1)
				} else {
					success = true
					logging.ProgressSpinner("stop")
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Deleted service "+contextID+".", 1)
				}
				timing.TimeOut(1, "s")
			}
		}
	} else {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "Service "+contextID+" already deleted.", 1)
	}
}

func (f F) Exists(namespace, contextID string) bool {
	return strings.ArrayContains(f.GetList(namespace, false), contextID)
}
