package golang

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	runtimeGo "runtime"

	"../../errors"
	"../../filesystem"
	"../../filesystem/compression"
	"../../io"
	"../../logging"
	"../../objects/strings"
	"../../runtime"
	"../../terminal"
	"../../timing"
	"../../vars"
)

type F struct{}

const urlBase = "https://dl.google.com/go/"
const urlVersionsJSON = "https://golang.org/dl/?mode=json"
const goDirLocal = "/usr/local/go"

func (f F) Check() bool {
	var updated bool
	var versionInstalled string
	var err error
	var versionNewest string
	var urlNewest string
	versionInstalled, _, err = f.GetVersion(true)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get available versions!", false, false, false) {
		versionNewest, urlNewest, err = f.GetVersion(false)
		if versionInstalled != versionNewest {
			updated = f.update(versionInstalled, versionNewest, urlNewest)
		} else {
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "Golang is up to date ("+versionInstalled+").", 0)
		}
	} else {
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiWarning}, "Unable to check for new golang version.", 0)
	}
	return updated
}

func (f F) update(versionInstalled, versionNewest, urlNewest string) bool {
	updated := false
	extractedPath := f.cleanup()
	sel := terminal.GetUserSelection("Do you want to update golang (local "+versionInstalled+") to newest version ("+versionNewest+")?", []string{}, false, true)
	if sel == "Yes" {
		logging.Log([]string{"\n", "", vars.EmojiProcess}, "Updating golang..", 0)
		hostOS := runtime.F.GetOS(runtime.F{})
		if hostOS == "darwin" {
			downloadPath := filesystem.GetWorkingDir() + versionNewest
			io.F.DownloadFile(io.F{}, downloadPath, urlNewest)
			compression.ExtractTarGz(downloadPath)
			logging.Log([]string{"", vars.EmojiDir, vars.EmojiProcess}, "Deleting old installation..", 0)
			logging.Log([]string{"", vars.EmojiDir, vars.EmojiCrypto}, "Asking for sudo permission..", 0)
			filesystem.Delete(goDirLocal, true)
			logging.Log([]string{"", vars.EmojiDir, vars.EmojiProcess}, "Moving new setup files..", 0)
			filesystem.Move(extractedPath, goDirLocal, true)
			filesystem.GiveProgramPermissions(goDirLocal, runtime.F.GetOSUsername(runtime.F{}))
			terminal.CreateAlias("go", goDirLocal+"/bin")
			f.cleanup()
			updated = true
		} else {
			logging.Log([]string{"", vars.EmojiProcess, vars.EmojiWarning}, "Your host OS ("+hostOS+") is not supported yet.", 0)
		}
	}
	return updated
}

func (f F) cleanup() string {
	extractedPath := filesystem.GetWorkingDir() + "go/"
	if filesystem.Exists(extractedPath) {
		logging.Log([]string{"", vars.EmojiProcess, ""}, "Cleaning up remaining artefacts from previous golang update..", 0)
		filesystem.Delete(extractedPath, false)
	}
	return extractedPath
}

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

func (f F) GetVersion(local bool) (string, string, error) {
	var err error
	var availableVersions []GoVersion
	if local {
		return "go" + strings.Trim(runtimeGo.Version(), "go"), "", err
	} else {
		availableVersions, err = f.getAvailableVersions()
		var newestVersion string
		var url string
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get available versions!", false, false, false) {
			newestVersion = availableVersions[0].Version
			url = urlBase + availableVersions[0].Files[1].Filename
		}
		return newestVersion, url, err
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
