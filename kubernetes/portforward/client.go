package portforward

import (
	// "flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	// "k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	// "path/filepath"
	"github.com/neurafuse/tools-go/kubernetes/client"
	"github.com/neurafuse/tools-go/vars"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/objects/strings"
	goStrings "strings"
	"sync"
	"syscall"
)

type PortForwardAPodRequest struct {
	// RestConfig is the kubernetes config
	RestConfig *rest.Config
	// Pod is the selected pod for this port forwarding
	Pod v1.Pod
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int
	// PodPort is the target port for the pod
	PodPort int
	// Steams configures where to write or read input from
	Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

func Connect(namespace, podID string, localPort int, podPort int) {
	var forwardDetails string = "localhost:"+strings.ToString(localPort)+" <-> "+podID+":"+strings.ToString(podPort)
	var wg sync.WaitGroup
	wg.Add(1)

	// stopCh control the port forwarding lifecycle. When it gets closed the
	// port forward will terminate
	stopCh := make(chan struct{}, 1)
	// readyCh communicate when the port forward is ready to get traffic
	readyCh := make(chan struct{})
	// stream is used to tell the port forwarder where to place its output or
	// where to expect input if needed. For the port forwarding we just need
	// the output eventually
	stream := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	// managing termination signal from the terminal. As you can see the stopCh
	// gets closed to gracefully handle its termination.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "Port-forwarding tunnel is going to be closed.", 0)
		close(stopCh)
		wg.Done()
	}()
	
	go func() {
		// PortForward the pod specified from its port 9090 to the local port
		// 8080
		err := PortForwardAPod(PortForwardAPodRequest{
			RestConfig: client.F.GetRestConfig(client.F{}),
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podID,
					Namespace: namespace,
				},
			},
			LocalPort: localPort,
			PodPort:   podPort,
			Streams:   stream,
			StopCh:    stopCh,
			ReadyCh:   readyCh,
		})
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to port-forward: "+forwardDetails, false, true, true)
	}()

	select {
	case <-readyCh:
		break
	}
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Port-forwarding tunnel is ready.", 0)
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiClient}, forwardDetails, 0)
	wg.Wait()
}

func PortForwardAPod(req PortForwardAPodRequest) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Pod.Namespace, req.Pod.Name)
	hostIP := goStrings.TrimLeft(req.RestConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
