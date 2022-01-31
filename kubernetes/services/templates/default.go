package templates

import (
	"github.com/neurafuse/tools-go/objects/strings"
	projectConfig "github.com/neurafuse/tools-go/config/project"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfig(id, clusterIP string, containerPorts [][]string) *apiv1.Service {
	ports := getPorts(containerPorts)
	var typ apiv1.ServiceType = getType()
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
			Labels: map[string]string{
				"app": id,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": id,
			},
			ClusterIP: clusterIP,
			Type:      typ,
		},
	}
	return service
}

func getType() apiv1.ServiceType {
	var t apiv1.ServiceType
	if projectConfig.F.NetworkMode(projectConfig.F{}, "port-forward") {
		t = "ClusterIP"
	} else if projectConfig.F.NetworkMode(projectConfig.F{}, "remote-url") {
		t = "LoadBalancer"
	}
	return t
}

func getPorts(containerPorts [][]string) []apiv1.ServicePort {
	var servicePorts []apiv1.ServicePort
	for i, port := range containerPorts {
		var servicePort apiv1.ServicePort
		servicePort.Name = "port-" + strings.ToString(i+1)
		servicePort.Port = strings.ToInt32(port[0])
		if port[1] == "TCP" {
			servicePort.Protocol = apiv1.ProtocolTCP
		} else {
			servicePort.Protocol = apiv1.ProtocolUDP
		}
		servicePorts = append(servicePorts, servicePort)
	}
	return servicePorts
}
