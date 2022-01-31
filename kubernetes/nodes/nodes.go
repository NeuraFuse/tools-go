package nodes

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/kubernetes/client"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type F struct{}

func (f F) Get(logResult bool) []string {
	list, err := f.getClient().List(context.TODO(), metav1.ListOptions{})
	var nodes []string
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+runtime.F.GetCallerInfo(runtime.F{}, true)+"!", false, false, true) {
		if logResult {
			logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, strings.Title(runtime.F.GetCallerInfo(runtime.F{}, true))+":", 0)
			if len(list.Items) == 0 {
				if logResult {
					logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
				}
			}
		}
		for i, d := range list.Items {
			if logResult {
				fmt.Printf(" [" + strconv.Itoa(i) + "] " + d.Name)
				fmt.Printf("\n")
			}
			nodes = append(nodes, d.Name)
		}
	}
	return nodes
}

func (f F) Taint() {

}

func (f F) getClient() corev1typed.NodeInterface {
	return client.F.GetAuth(client.F{}).CoreV1().Nodes()
}
