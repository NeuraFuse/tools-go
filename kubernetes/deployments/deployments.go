package deployments

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1typed "k8s.io/client-go/kubernetes/typed/apps/v1"
	apiApps "k8s.io/api/apps/v1"
)

type F struct{}

func (f F) GetList(namespace string, logResult bool) []string {
	client := f.getClient(namespace)
	list, err := client.List(contextPack.TODO(), metav1.ListOptions{})
	var deployments []string
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get deployment in namespace "+namespace+":", false, false, true) {
		if logResult {
			logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, strings.Title(runtime.F.GetCallerInfo(runtime.F{}, true))+":", 0)
			if len(list.Items) == 0 {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "There are no deployment.", 0)
			}
		}
		for i, d := range list.Items {
			if logResult {
				fmt.Printf(" [" + strconv.Itoa(i) + "] " + d.Name)
				fmt.Printf(" | " + strconv.Itoa(int(*d.Spec.Replicas)) + " Pod")
				fmt.Printf("\n")
			}
			deployments = append(deployments, d.Name)
		}
	}
	return deployments
}

func (f F) GetID(namespace, contextID string) *apiApps.Deployment {
	client := f.getClient(namespace)
	deployment, err := client.Get(contextPack.TODO(), contextID, metav1.GetOptions{})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get deployment with ID "+contextID+" in namespace "+namespace+"!", false, false, true)
	return deployment
}

func (f F) Create(namespace, contextID, imageAddrs, serviceCluster, resources string, volumes, containerPort [][]string) {
	client := f.getClient(namespace)
	deployment := templates.GetConfig(contextID, imageAddrs, serviceCluster, resources, volumes, containerPort)
	_, err := client.Create(contextPack.TODO(), deployment, metav1.CreateOptions{})
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), ""+contextID+" already exists.", false, false, true) {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created deployment "+contextID+".", 0)
	} else {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "Fail recovery: Deployment recreation..", 0)
		f.Delete(namespace, contextID)
		f.Create(namespace, contextID, imageAddrs, serviceCluster, resources, volumes, containerPort)
	}
}

func (f F) Delete(namespace, contextID string) {
	deletePolicy := metav1.DeletePropagationForeground
	client := f.getClient(namespace)
	if err := client.Delete(contextPack.TODO(), contextID, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), contextID+" already deleted.", true, false, true)
	} else {
		var success bool = false
		logging.ProgressSpinner("start")
		for ok := true; ok; ok = !success {
			if strings.ArrayContains(f.GetList(namespace, false), contextID) {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for the deployment "+contextID+" to be deleted..", 0)
				timing.TimeOut(1, "s")
			} else {
				success = true
				logging.ProgressSpinner("stop")
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Deleted deployment "+contextID+".", 1)
			}
		}
	}
}

func (f F) Exists(namespace, contextID string) bool {
	return strings.ArrayContains(f.GetList(namespace, false), contextID)
}

func (f F) getClient(namespace string) appsv1typed.DeploymentInterface {
	if namespace == "" {
		namespace = namespaces.Default
	}
	return client.F.GetAuth(client.F{}).AppsV1().Deployments(namespace)
}