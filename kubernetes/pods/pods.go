package pods

import (
	"bytes"
	contextPack "context"
	"fmt"
	"io"
	"strconv"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/kubernetes/client"
	"github.com/neurafuse/tools-go/kubernetes/deployments"
	"github.com/neurafuse/tools-go/kubernetes/namespaces"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	ostrings "github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type F struct{}

func (f F) GetIDFromName(namespace, podName string) string {
	var podID string
	var pods []string = f.Get(namespace, false)
	for _, pod := range pods {
		if strings.HasPrefix(pod, strings.ToLower(podName)) {
			podID = pod
		}
	}
	var err error
	if podID == "" {
		err = errors.New("Unable to find podID for podName: "+podName+" in namespace "+namespace+"!")
	}
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get information!", false, true, true)
	return podID
}

func (f F) GetList(namespace, contextID string) (*apiv1.PodList, error) {
	options := metav1.ListOptions{
		LabelSelector: "app=" + contextID,
	}
	list, err := f.getClient(namespace).List(contextPack.TODO(), options)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get pod list for deployment "+contextID+".", false, true, true)
	return list, err
}

func (f F) Get(namespace string, logResult bool) []string {
	var pods []string
	deployments := deployments.F.Get(deployments.F{}, namespace, false)
	if logResult {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInfo}, "Pods:", 0)
	}
	if len(deployments) != 0 {
		for i, contextID := range deployments {
			list, _ := f.GetList(namespace, contextID)
			if len(list.Items) == 0 {
				logging.Log([]string{"", "", ""}, "There are no "+runtime.F.GetCallerInfo(runtime.F{}, true)+".", 0)
			}
			for _, p := range (*list).Items {
				pods = append(pods, p.Name)
				if logResult {
					fmt.Printf("[" + strconv.Itoa(i) + "] " + p.Name)
					fmt.Printf(" | " + string(p.Status.Phase))
					for iV, volume := range p.Spec.Volumes {
						fmt.Printf(" | ")
						fmt.Printf(volume.Name)
						if iV > 0 && !(len(p.Spec.Volumes) == iV+1) {
							fmt.Printf(", ")
						}
					}
					if string(p.Spec.NodeName) != "" {
						fmt.Printf(" | " + string(p.Spec.NodeName))
					}
					fmt.Printf("\n")
				}
			}
		}
	} else {
		if logResult {
			logging.Log([]string{"", "", ""}, "There are no deployments and therefore also no pods.", 0)
		}
	}
	return pods
}

func (f F) Delete(namespace, contextID string) {
	podID, _ := f.getPod(namespace, contextID)
	err := f.getClient(namespace).Delete(contextPack.TODO(), podID.Name, metav1.DeleteOptions{})
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), ""+contextID+" already deleted.", false, false, true) {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Pods from deployment "+contextID+" deleted.", 1)
	}
}

var logscontextIDLast string

//var logsErrStreamCh chan error
var logsErrStream error

func (f F) Logs(namespace, contextID, waitForStatusInLog string, parallel bool, initWaitDuration int) {
	err := f.WaitForPhase(namespace, contextID, "running", initWaitDuration)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get logs for "+contextID+"!", false, false, true) {
		if contextID != logscontextIDLast {
			logscontextIDLast = contextID
			var finished bool
			//logsErrStreamCh = make(chan error)
			for ok := true; ok; ok = !finished {
				if parallel {
					go f.LogsRoutine(namespace, contextID, waitForStatusInLog)
				} else {
					f.LogsRoutine(namespace, contextID, waitForStatusInLog)
				}
				//errStreamCh := <-logsErrStreamCh
				if logsErrStream != nil {
					err = f.WaitForPhase(namespace, contextID, "running", 1)
					if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to stream pod logs because the deployment "+contextID+" doesn't exist!", false, false, true) {
						finished = true
					}
				} else {
					finished = true
				}
			}
		}
	}
}

var InterruptPodLogsLive bool

func (f F) LogsRoutine(namespace, contextID, waitForStatusInLog string) {
	var logsLast string
	if waitForStatusInLog != "" {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiInspect}, "Waiting for pod log status \""+waitForStatusInLog+"\" to occur in "+contextID+"..", 0)
	} else {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInspect}, "Streaming pod logs from "+contextID+" to terminal..", 0)
		logging.PartingLine()
	}
	var logsCurrent string
	var errStream error
	for {
		if !InterruptPodLogsLive {
			logsCurrent, errStream = f.GetLogs(namespace, contextID)
			if errStream == nil {
				if logsCurrent != logsLast {
					if logsLast != "" && waitForStatusInLog == "" {
						logging.Log([]string{"\n", "", ""}, strings.TrimPrefix(logsCurrent, logsLast), 0) // strings.Trim(logsCurrent, logsLast)
					} else if waitForStatusInLog == "" {
						logging.Log([]string{"\n", "", ""}, logsCurrent, 0)
					}
					if waitForStatusInLog != "" {
						if strings.Contains(logsCurrent, waitForStatusInLog) {
							InterruptPodLogsLive = true
						}
					}
				}
				logsLast = logsCurrent
				timing.Sleep(1, "s")
			} else {
				InterruptPodLogsLive = true
			}
		} else {
			if waitForStatusInLog != "" {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInspect}, "Specific log \""+waitForStatusInLog+"\" occured in "+contextID+".", 0)
			} else {
				logging.PartingLine()
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "Ended container logs streaming to terminal due to interrupt.\n", 0)
			}
			InterruptPodLogsLive = false
			break
		}
	}
	logsErrStream = errStream // TODO: chan logsErrStreamCh <- errStream
}

func (f F) GetLogs(namespace, contextID string) (string, error) {
	pod, _ := f.getPod(namespace, contextID)
	podLogOpts := apiv1.PodLogOptions{}
	req := f.getClient(namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, errStream := req.Stream(contextPack.TODO())
	var logsStr string
	if !errors.Check(errStream, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+contextID+" pod logs because it was deleted!", false, false, true) {
		defer podLogs.Close()
		buf := new(bytes.Buffer)
		_, errCopyBuf := io.Copy(buf, podLogs)
		errors.Check(errCopyBuf, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to copy information from podLogs to buffer!", false, false, true)
		logsStr = buf.String()
	}
	return logsStr, errStream
}

func (f F) GetContainers(namespace, contextID string) []apiv1.Container {
	pod, _ := f.getPod(namespace, contextID)
	return pod.Spec.Containers
}

func (f F) GetContainerNamesList(namespace, contextID string) []string {
	var containers []apiv1.Container = f.GetContainers(namespace, contextID)
	var containerNames []string
	for _, container := range containers {
		containerNames = append(containerNames, container.Name)
	}
	return containerNames
}

func (f F) GetContainerIDByName(namespace, contextID, containerName string) int {
	var containers []apiv1.Container = f.GetContainers(namespace, contextID)
	for i, container := range containers {
		if containerName == container.Name {
			return i
		}
	}
	errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get container ID by name!", true, false, true)
	return -1
}

func (f F) getContainer(namespace, contextID string, containerID int) apiv1.Container {
	return f.GetContainers(namespace, contextID)[containerID]
}

func (f F) GetContainerImgAddrs(namespace, contextID string, containerID int) string {
	return f.getContainer(namespace, contextID, containerID).Image
}

func (f F) getPod(namespace, contextID string) (apiv1.Pod, error) {
	var podReturn apiv1.Pod
	list, err := f.GetList(namespace, contextID)
	for _, pod := range (*list).Items {
		podReturn = pod
		break
	}
	return podReturn, err
}

func (f F) WaitForPhase(namespace, contextID, phase string, initWaitDuration int) error {
	if initWaitDuration == 0 {
		initWaitDuration = 1
	}
	phase = strings.Title(phase)
	phaseCurrent, err := f.GetPhase(namespace, contextID)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to wait for "+contextID+" pod phase!", false, false, true) {
		if !(phaseCurrent == phase) {
			i := 0
			var success bool
			for ok := true; ok; ok = !success {
				phaseCurrent, _ = f.GetPhase(namespace, contextID)
				if phaseCurrent == phase {
					success = true
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Pod "+contextID+" is now in phase "+strings.ToLower(phase)+".", 0)
				} else {
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting (max. "+ostrings.ToString(initWaitDuration)+" minutes) for "+contextID+" pod phase to be "+strings.ToLower(phase)+"..", 0)
					timing.Sleep(1, "s")
				}
				i++
				if i == initWaitDuration*60 {
					err = errors.New("Waited more than " + ostrings.ToString(initWaitDuration) + " minutes for " + contextID + " pod phase " + phase + "!")
					break
				}
			}
		}
	}
	return err
}

func (f F) GetPhase(namespace, contextID string) (string, error) {
	var pod apiv1.Pod
	var phase string
	var err error
	var success bool
	for ok := true; ok; ok = !success {
		pod, err = f.getPod(namespace, contextID)
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Waiting for "+contextID+" pods to be created to fetch status..", false, false, true) {
			phase = string(pod.Status.Phase)
			success = true
		} else {
			timing.Sleep(1, "s")
		}
	}
	return phase, err
}

func (f F) WaitForCondition(namespace, contextID, conditionType string, waitMinsMax int) error {
	var conditions []apiv1.PodCondition
	var err error
	conditions, err = f.GetConditions(namespace, contextID)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to get "+contextID+" pod conditions!", false, false, true) {
		var success bool
		i := 0

		for ok := true; ok; ok = !success {
			conditions, err = f.GetConditions(namespace, contextID)
			if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to get "+contextID+" pod conditions!", false, false, true) {
				break
			}
			var conditionTypeI string
			for _, condition := range conditions {
				conditionTypeI = string(condition.Type)
				if conditionTypeI == conditionType {
					success = true
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Pod "+contextID+" is now in condition "+conditionType+".", 0)
					break
				} else {
					logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting (max. "+ostrings.ToString(waitMinsMax)+" minutes) for pod condition to be "+conditionType+"..", 0)
					timing.Sleep(1, "s")
				}
			}
			i++
			if i == waitMinsMax*60 {
				err = errors.New("Waited more than " + ostrings.ToString(waitMinsMax) + " minutes for " + contextID + " pod condition " + conditionType + "!")
				break
			}
		}
	}
	return err
}

func (f F) GetConditions(namespace, contextID string) ([]apiv1.PodCondition, error) {
	pod, err := f.getPod(namespace, contextID)
	return pod.Status.Conditions, err
}

func (f F) Exists(namespace, contextID string) bool {
	return strings.ArrayContains(f.Get(namespace, false), contextID)
}

func (f F) getClient(namespace string) corev1typed.PodInterface {
	if namespace == "" {
		namespace = namespaces.Default
	}
	return client.F.GetAuth(client.F{}).CoreV1().Pods(namespace)
}
