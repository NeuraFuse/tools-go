package namespaces

import (
	"context"
	"fmt"
	"strconv"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/kubernetes/client"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type F struct{}

var Default string = vars.OrganizationNameRepo

func (f F) Init() {
	f.Create(Default)
}

func (f F) Get(logResult bool) (error, []string) {
	var err error
	var list *apiv1.NamespaceList
	var namespaces []string
	list, err = f.getClient().List(context.TODO(), metav1.ListOptions{})
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+runtime.F.GetCallerInfo(runtime.F{}, true)+"!", false, false, true) {
		if logResult {
			logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, strings.Title(runtime.F.GetCallerInfo(runtime.F{}, true))+":", 0)
			if len(list.Items) == 0 {
				logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
			}
		}
		for i, n := range list.Items {
			if logResult {
				fmt.Printf(" [" + strconv.Itoa(i) + "] " + n.Name)
				fmt.Printf("\n")
			}
			namespaces = append(namespaces, n.Name)
		}
	}
	return err, namespaces
}

func (f F) Create(id string) {
	nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: id}}
	_, err := f.getClient().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create namespace "+id+" (exists or unauthorized)!", false, false, false) {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created namespace "+id+".", 0)
	}
}

func (f F) Delete(id string) {
	if f.Exists(id) {
		err := f.getClient().Delete(context.TODO(), id, metav1.DeleteOptions{})
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete namespace "+id+"!", false, true, true) {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Deleted namespace "+id+".", 0)
		}
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), ""+id+" already deleted.", true, false, true)
	}
}

func (f F) filter(ns []apiv1.Namespace, filter []string) []string {
	var filtered []string
	var doFilter bool
	for _, eA := range ns {
		for _, eF := range filter {
			if eF == eA.Name {
				doFilter = true
			}
		}
		if !doFilter {
			filtered = append(filtered, eA.Name)
		}
		doFilter = false
	}
	return filtered
}

func (f F) Exists(id string) bool {
	var list []string
	_, list = f.Get(false)
	return strings.ArrayContains(list, id)
}

func (f F) getClient() corev1typed.NamespaceInterface {
	return client.F.GetAuth(client.F{}).CoreV1().Namespaces()
}
