package releases

import (
	"github.com/neurafuse/tools-go/build"
	"github.com/neurafuse/tools-go/container"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/git"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/updater"
	"github.com/neurafuse/tools-go/vars"
	//"github.com/neurafuse/tools-go/releases/license"
)

type F struct{}

func (f F) Router(cliArgs []string) {
	var mode string
	var modes []string = []string{"For this local environment", "Remote release"}
	if len(cliArgs) < 3 {
		mode = terminal.GetUserSelection("What kind of "+env.F.GetActive(env.F{}, true)+" release do you want to build?", modes, false, false)
	} else {
		mode = cliArgs[2]
	}
	var context string = env.F.GetActive(env.F{}, false)
	f.processBuild(context, mode, modes)
}

func (f F) processBuild(context, mode string, modes []string) {
	//license.F.CreateFile(license.F{})
	git.F.CreateIgnoreFile(git.F{})
	updater.F.CreateRepoInfoFile(updater.F{})
	switch mode {
	case modes[0]:
		f.localBuild(context)
	case modes[1]:
		f.remoteBuild(context)
		f.publishBuild(context)
	}
}

func (f F) publishBuild(context string) {
	var publish string = terminal.GetUserSelection("Do you want to publish this remote build?", []string{}, false, true)
	if publish == "Yes" {
		f.remotePublish(context)
	}
}

func (f F) localBuild(context string) {
	f.createBuildFiles(false)
}

func (f F) remoteBuild(context string) {
	f.createBuildFiles(true)
	container.F.CheckUpdates(container.F{}, context, false, true)
}

func (f F) remotePublish(context string) {
	container.F.CheckUpdates(container.F{}, context, true, true)
}

func (f F) createBuildFiles(crossCompile bool) {
	if crossCompile {
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiInspect}, "Cross-compiling for linux/macOS architecture pairs..", 0)
		osPairs := runtime.F.GetOSArchitecturePairs(runtime.F{}, "linux")
		osPairs = append(osPairs, runtime.F.GetOSArchitecturePairs(runtime.F{}, "macos")...)
		for _, pair := range osPairs {
			build.F.Make(build.F{}, env.F.GetActive(env.F{}, false), pair[0], pair[1], updater.F.GetRepoUpdateBuildsDir(updater.F{}), false, false)
		}
	} else {
		build.F.Make(build.F{}, env.F.GetActive(env.F{}, false), "", "", "", true, true)
	}
}
