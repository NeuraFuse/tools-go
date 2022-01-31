package golang

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	runtimeGo "runtime"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

const urlVersionsJSON = "https://golang.org/dl/?mode=json"

type GoVersion struct {
	Version string       `json:"version"`
	Stable  bool         `json:"stable"`
	Files   []GoDownload `json:"files"`
}

type GoDownload struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}

func (f F) Check() bool {
	var updated bool
	var versionInstalled string
	var err error
	var versionNewest string
	versionInstalled, err = f.GetVersion(true)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get available versions!", false, false, false) {
		versionNewest, err = f.GetVersion(false)
		var logMsg string
		var logEmoji string
		if versionInstalled == versionNewest {
			logMsg = "Golang is up to date (" + versionInstalled + ")."
			logEmoji = vars.EmojiSuccess
		} else {
			logMsg = "Golang is not up to date (" + versionInstalled + " --> " + versionNewest + ")."
			logEmoji = vars.EmojiWarning
		}
		logging.Log([]string{"", vars.EmojiDev, logEmoji}, logMsg, 0)
	} else {
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiWarning}, "Unable to check for new golang versions.", 0)
	}
	return updated
}

func (f F) GetVersion(local bool) (string, error) {
	var err error
	var availableVersions []GoVersion
	if local {
		return "go" + strings.Trim(runtimeGo.Version(), "go"), err
	} else {
		availableVersions, err = f.getAvailableVersions()
		var newestVersion string
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get available versions!", false, false, false) {
			newestVersion = availableVersions[0].Version
		}
		return newestVersion, err
	}
}

func (f F) getAvailableVersions() ([]GoVersion, error) {
	client := http.Client{
		Timeout: timing.GetTimeDuration(10, "s"),
	}
	req, err := http.NewRequest(http.MethodGet, urlVersionsJSON, nil)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create HTTP request!", false, false, true)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	var availableVersions []GoVersion
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Request failed!", false, false, false) {
		body, err := ioutil.ReadAll(resp.Body)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to read body!", false, false, true)
		err = json.Unmarshal(body, &availableVersions)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to parse JSON!", false, false, true)
	}
	return availableVersions, err
}
