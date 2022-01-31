package python

import (
	"github.com/neurafuse/tools-go/container/runtime/tools"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/exec"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/kubernetes/resources"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(project, module, pathExec, dataPath, serverSyncWaitMsg string) {
	f.checkClusterResources()
	var firstExecRun bool = true
	var embeddingFirstRun bool
	var argsExec []string = []string{pathExec}
	argsExec = append(argsExec, module)
	for {
		if !embeddingFirstRun {
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiInfo}, "pathExec: "+pathExec, 0)
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiInfo}, "argsExec: "+strings.Join(argsExec, " "), 0)
		}
		if filesystem.Exists(pathExec) {
			if firstExecRun {
				timing.Sleep(5, "s")
				firstExecRun = false
			}
			logging.PartingLine()
			logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "Starting project "+project+"..\n", 0)
			exec.WithLiveLogs("python", argsExec, true)
			logging.Log([]string{"", vars.EmojiLink, vars.EmojiInfo}, "Auto restart after 1s..", 0)
			logging.PartingLine()
		} else {
			if !embeddingFirstRun {
				logging.Log([]string{"", vars.EmojiClient, vars.EmojiAPI}, serverSyncWaitMsg, 0)
				tools.F.DataAggregation(tools.F{}, module, project, dataPath)
				embeddingFirstRun = true
			}
		}
		timing.Sleep(1, "s")
	}
}

func (f F) checkClusterResources() {
	if env.F.Container(env.F{}) {
		resources.Check("container", "tpu")
	}
}
