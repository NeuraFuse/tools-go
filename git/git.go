package git

import (
	"../filesystem"
	"../env"
	"../runtime"
	"../vars"
	"../logging"
)

type F struct{}
var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) CreateIgnoreFile() {
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiSettings}, "Creating "+context+" ignore file..", 0)
	var filePath string = ".gitignore"
	if filesystem.Exists(filePath) {
		filesystem.Delete(filePath, false)
	}
	for _, path := range f.getIgnorePaths() {
		filesystem.AppendStringToFile(filePath, path)
	}
}

func (f F) getIgnorePaths() []string {
	var executableFileName string = env.F.GetActive(env.F{}, false)+"-"+runtime.F.GetOS(runtime.F{})+"-"+runtime.F.GetOSArchitecture(runtime.F{})
	var paths []string = []string{"users", "releases/builds", executableFileName, "**/.DS_Store"}
	if env.F.API(env.F{}) {
		paths = append(paths, "server/http/certs")
	}
	return paths
}