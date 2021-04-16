package daemonsets

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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1typed "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type F struct{}

var namespaceDefault string = "kube-system"

func (f F) Get(logResult bool) []string {
	client := f.getClient()
	list, err := client.List(contextPack.TODO(), metav1.ListOptions{})
	var daemonsets []string
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+runtime.F.GetCallerInfo(runtime.F{}, true)+" in namespace "+namespaceDefault+"!", false, false, true) {
		if logResult {
			logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, strings.Title(runtime.F.GetCallerInfo(runtime.F{}, true))+":", 0)
			if len(list.Items) == 0 {
				logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
			}
		}
		for i, d := range list.Items {
			if logResult {
				fmt.Printf(" [" + strconv.Itoa(i) + "] " + d.Name)
				fmt.Printf("\n")
			}
			daemonsets = append(daemonsets, d.Name)
		}
	}
	return daemonsets
}

func (f F) Create(ds *appsv1.DaemonSet) {
	contextID := ds.ObjectMeta.Name
	client := f.getClient()
	_, err := client.Create(contextPack.TODO(), ds, metav1.CreateOptions{})
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), ""+contextID+" already exists.", false, false, true) {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created "+runtime.F.GetCallerInfo(runtime.F{}, true)+" "+contextID+".", 1)
	}
}

func (f F) Exists(contextID string) bool {
	return strings.ArrayContains(f.Get(false), contextID)
}

func (f F) getClient() appsv1typed.DaemonSetInterface {
	return client.F.GetAuth(client.F{}).AppsV1().DaemonSets(namespaceDefault)
}
