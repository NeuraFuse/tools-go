package updater

import (
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/io"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/logging/emoji"
	"github.com/neurafuse/tools-go/readers/json"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

type Info struct {
	Version struct {
		Recent string `json:"recent"`
	} `json:"version"`
}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Check() {
	update, versionInstalled, versionRecent, err := f.checkVersion()
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Version check failed!", false, false, true) {
		if update {
			if !config.DevConfigActive() {
				if config.Setting("get", "cli", "Spec.Updates.Auto.Status", "") != "disabled" {
					f.update(versionRecent)
				} else {
					emoji.Println("", vars.EmojiGlobe, vars.EmojiWarning, "There is an update available but you have turned off auto. updates.")
					emoji.Println("", vars.EmojiWarning, vars.EmojiInfo, "You can enable auto updates via the CLI Settings.\n")
				}
			} else {
				emoji.Println("", vars.EmojiDev, vars.EmojiGlobe, "Skipped available update due to developer mode.")
			}
		} else {
			emoji.Println("", vars.EmojiGlobe, vars.EmojiSuccess, "Version ("+versionInstalled+") is up to date.\n")
		}
	}
}

func (f F) checkVersion() (bool, string, string, error) {
	var update bool
	versionInstalled := env.F.GetVersion(env.F{})
	versionRecent, err := f.getVersionRecent()
	if versionInstalled != versionRecent {
		update = true
		emoji.Println("", vars.EmojiAssistant, vars.EmojiGlobe, "There is an update available for "+env.F.GetActive(env.F{}, true)+".")
		emoji.Println("", vars.EmojiAssistant, vars.EmojiGlobe, "Version "+versionInstalled+" --> "+versionRecent+"\n")
	}
	return update, versionInstalled, versionRecent, err
}

func (f F) getVersionRecent() (string, error) {
	var versionRecent string
	var err error
	var info *Info = new(Info)
	var url string = f.getUpdateInfoURL()
	if io.F.Reachable(io.F{}, url) {
		json.URLToInterface(url, info)
	} else {
		err = errors.New("Unable to get recent " + env.F.GetActive(env.F{}, true) + " version (" + url + " not reachable)!")
	}
	if err == nil {
		versionRecent = info.Version.Recent
	}
	return versionRecent, err
}

func (f F) getProvider() string {
	var provider string
	provider = "github"
	return provider
}

func (f F) getURLBase() string {
	var protocol string = "https://"
	var tld string = "com"
	return protocol + f.getProvider() + "." + tld + "/"
}

func (f F) getUpdateURLPath() string {
	var url string = f.getURLBase()
	switch f.getProvider() {
	case "github":
		var branch string = "master"
		url = url + vars.OrganizationNameRepo + "/" + env.F.GetActive(env.F{}, false) + "/blob/" + branch + "/"
	case "neurafuse":
		url = url + env.F.GetActive(env.F{}, false) + "/"
	}
	return url + f.GetRepoUpdateDir()
}

func (f F) getReleaseDownloadURL(version string) string {
	var url string = f.getURLBase()
	switch f.getProvider() {
	case "github":
		url = url + vars.OrganizationNameRepo + "/" + env.F.GetActive(env.F{}, false) + "/releases/download/" + version + "/"
		url = url + f.getReleaseDownloadFile()
	case "neurafuse":
		//url = url + env.F.GetActive(env.F{}, false) + "/" TODO:
	}
	return url
}

func (f F) getReleaseDownloadFile() string {
	return env.F.GetActive(env.F{}, false) + "-" + runtime.F.GetOS(runtime.F{}) + "-" + runtime.F.GetOSArchitecture(runtime.F{})
}

func (f F) CreateRepoInfoFile() {
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiSettings}, "Creating repo "+context+" info file..", 0)
	var filePath string = f.GetRepoUpdateDir() + f.getUpdateInfoFile()
	if filesystem.Exists(filePath) {
		filesystem.Delete(filePath, false)
	}
	var info *Info = new(Info)
	info.Version.Recent = env.F.GetVersion(env.F{})
	json.StructToFile(filePath, info)
}

func (f F) getUpdateInfoURL() string {
	return f.getUpdateURLPath() + f.getUpdateInfoFile()
}

func (f F) GetRepoUpdateDir() string {
	var repoUpdateDir string = "releases/"
	return repoUpdateDir
}

func (f F) GetRepoUpdateBuildsDir() string {
	var repoUpdateBuildsDir string = "builds/"
	return f.GetRepoUpdateDir() + repoUpdateBuildsDir
}

func (f F) getUpdateInfoFile() string {
	return "info.json"
}

func (f F) getTmpDir() string {
	return "tmp/" + context + "/"
}

func (f F) update(versionRecent string) {
	emoji.Println("", vars.EmojiProcess, vars.EmojiInfo, "Starting update..")
	url := f.getReleaseDownloadURL(versionRecent)
	tmpPath := f.getTmpDir() + "update/" + f.getReleaseDownloadFile()
	err := io.F.DownloadFile(io.F{}, tmpPath, url)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to download recent "+env.F.GetActive(env.F{}, true)+" version!", false, false, true) {
		filesystem.Delete(f.getReleaseDownloadFile(), false)
		filesystem.Move(tmpPath, f.getReleaseDownloadFile(), false)
		emoji.Println("", vars.EmojiProcess, vars.EmojiSuccess, "Updated "+env.F.GetActive(env.F{}, false)+".")
		emoji.Println("", vars.EmojiProcess, vars.EmojiInfo, "Please restart to apply the updates.")
		terminal.Exit(0, "")
	}
}
