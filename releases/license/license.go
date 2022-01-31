package license

import (
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) CreateFile() {
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiSettings}, "Creating "+context+" file..", 0)
	var filePath string = "LICENSE"
	if filesystem.Exists(filePath) {
		filesystem.Delete(filePath, false)
	}
	var content string = f.get()
	filesystem.AppendStringToFile(filePath, content)
}

func (f F) get() string {
	var license string
	license = "123"
	return license
}
