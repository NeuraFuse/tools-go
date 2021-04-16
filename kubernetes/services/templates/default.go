package templates

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"../../../../tools-go/objects/strings"
)

func GetConfig(id, clusterIP string, containerPorts [][]string) *apiv1.Service {
	ports := getServicePorts(containerPorts)
	service := &apiv1.Service {
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
			Type: "LoadBalancer",
		},
	}
	return service
}

func getServicePorts(containerPorts [][]string) []apiv1.ServicePort {
	var servicePorts []apiv1.ServicePort
	for i, port := range containerPorts {
		var servicePort apiv1.ServicePort
		servicePort.Name = "port-"+strings.ToString(i+1)
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