package tensorboard

import (
	"../../../exec"
	"../../../logging"
	"../../../vars"
	"../../../env"
)

type F struct{}

var ContainerPorts []string = []string{"6006", "TCP"}
var Protocol string = "http"

func (f F) Start() {
	var logDir string
	baseDir := "lightning/pytorch/data/training/logs"
	var tbArgs []string
	if env.F.Container(env.F{}) {
		logDir = baseDir
		tbArgs = append(tbArgs, "--bind_all")
	} else {
		logDir = "../" + baseDir
		f.Info("localhost")
	}
	tbArgs = append(tbArgs, []string{"--logdir", logDir}...)
	go exec.WithLiveLogs("tensorboard", tbArgs, false)
}

func (f F) Info(ip string) {
	url := Protocol + "://" + ip + ":" + ContainerPorts[0]
	logging.Log([]string{"", vars.EmojiRobot, vars.EmojiInfo}, "The TensorBoard web interface is now available at: "+url+"\n", 0)
}
