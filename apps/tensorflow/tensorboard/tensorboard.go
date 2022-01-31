package tensorboard

import (
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/exec"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var Protocol string = "http"

func (f F) GetContainerPorts() []string {
	var containerPorts []string = []string{"6006", "TCP"}
	return containerPorts
}

func (f F) Start() {
	var logDir string
	var baseDir string = "lightning-py/data/training/logs"
	var tbArgs []string
	if env.F.Container(env.F{}) {
		logDir = baseDir
		tbArgs = append(tbArgs, "--bind_all")
	} else {
		logDir = "../" + baseDir
	}
	tbArgs = append(tbArgs, []string{"--logdir", logDir}...)
	go exec.WithLiveLogs("tensorboard", tbArgs, false)
	f.LogInfo("localhost")
}

func (f F) LogInfo(ip string) {
	var url string = Protocol + "://" + ip + ":" + f.GetContainerPorts()[0]
	logging.Log([]string{"", vars.EmojiRobot, vars.EmojiInfo}, "The TensorBoard web interface is now available at: "+url+"\n", 0)
}
