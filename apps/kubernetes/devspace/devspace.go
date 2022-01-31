package devspace

import (
	"github.com/neurafuse/tools-go/apps/kubernetes/devspace/templates"
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/container"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/exec"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/io"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Sync(contextID, namespace, imageAddrs string, thread bool) {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Preparing file synchronization..", 0)
	//devspace.Sync(createDevSpaceConfig(templateBasePath, templatePath))
	configPath := f.createDevSpaceConfig(contextID, namespace, imageAddrs)
	f.syncExec(configPath, thread)
}

func (f F) createDevSpaceConfig(contextID, namespace, imageAddrs string) string {
	configInt := f.getConfigInterface(contextID, f.getTemplateBasePath(), f.getPath("template"), namespace, imageAddrs)
	projectConfigPath := f.getPath("projectconfig")
	if filesystem.Exists(projectConfigPath) {
		filesystem.Delete(projectConfigPath, false)
	}
	yaml.StructToFile(projectConfigPath, configInt)
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Created devspace configs.", 0)
	return projectConfigPath
}

func (f F) getConfigInterface(contextID, templateBasePath, templatePath, namespace, imageAddrs string) interface{} {
	dConfig := yaml.FileToStruct(templateBasePath+templatePath, &templates.Default{})
	var repoPreamble string = config.Setting("get", "dev", "Spec.Containers.Registry.Address", "")
	dConfig.(*templates.Default).Images.Remote.Image = repoPreamble + imageAddrs
	var sync []struct {
		ImageName     string   "json:\"imageName\""
		Namespace     string   "json:\"namespace\""
		LocalSubPath  string   "json:\"localSubPath\""
		ContainerPath string   "json:\"containerPath\""
		ExcludePaths  []string "json:\"excludePaths\""
	}
	var syncEntry struct {
		ImageName     string   "json:\"imageName\""
		Namespace     string   "json:\"namespace\""
		LocalSubPath  string   "json:\"localSubPath\""
		ContainerPath string   "json:\"containerPath\""
		ExcludePaths  []string "json:\"excludePaths\""
	}
	syncEntry.ImageName = "remote"
	syncEntry.Namespace = namespace
	var localAppRoot string = config.Setting("get", "project", "Spec.Containers.Sync.PathMappings.LocalAppRoot", "")
	if localAppRoot == filesystem.GetWorkspaceFolderVar() {
		localAppRoot = ""
	}
	syncEntry.LocalSubPath = localAppRoot
	var containerAppRoot string = config.Setting("get", "project", "Spec.Containers.Sync.PathMappings.ContainerAppRoot", "")
	syncEntry.ContainerPath = containerAppRoot
	syncEntry.ExcludePaths = []string{container.F.GetVolumeMountPath(container.F{}) + "training/input"}
	sync = append(sync, syncEntry)
	//dConfig.(*templates.Default).Dev.Sync[0].ImageName = context
	//dConfig.(*templates.Default).Dev.Sync[0].Namespace = namespace
	//dConfig.(*templates.Default).Dev.Sync[0].ExcludePaths = []string{"lost+found"}
	dConfig.(*templates.Default).Dev.Sync = sync
	dConfig.(*templates.Default).Version = "v1beta9"
	return dConfig
}

func (f F) getBasePath() string {
	var basePath string = "../tools-go@" + vars.ToolsGoVersion + "/apps/kubernetes/devspace/"
	return basePath
}

func (f F) getTemplateBasePath() string {
	return f.getBasePath() + "/templates/"
}

func (f F) getPath(id string) string {
	var path string
	var ext string = ".yaml"
	switch id {
	case "template":
		{
			path = f.getTemplateBasePath() + "default" + ext
		}
	case "projectconfig":
		{
			path = filesystem.GetWorkingDir(false) + "." + context + "/" + context + ext
		}
	}
	return path
}

func (f F) syncExec(configPath string, thread bool) {
	f.CheckSetup()
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiLink}, "Starting file synchronization..", 0)
	logging.Log([]string{"\n", vars.EmojiLink, vars.EmojiInfo}, "Please wait a few seconds for your files to get synced into the container..", 0)
	logging.Log([]string{"", vars.EmojiLink, vars.EmojiInfo}, "You can stop the file synchronization with CTRL+C.\n", 0)
	args := []string{"sync", "--config", configPath}
	if thread {
		go exec.WithLiveLogs(context, args, false)
	} else {
		exec.WithLiveLogs(context, args, false)
	}
}

func (f F) CheckSetup() {
	if !f.setupExists() {
		f.install()
	} else {
		f.upgrade()
	}
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "DevSpace setup is now ready to be used.", 0)
}

func (f F) getSetupFilePath() string {
	var setupFilePath string = f.getBasePath() + context
	return setupFilePath
}

func (f F) getInstallPath() string {
	var installPath string = runtime.F.GetOSInstallDir(runtime.F{}) + context
	return installPath
}

func (f F) install() {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Installing DevSpace..", 0)
	var err error = f.downloadSetupFile()
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to download setup file!", false, false, true) {
		filesystem.Move(f.getSetupFilePath(), f.getInstallPath(), true)
		filesystem.GiveProgramPermissions(f.getInstallPath(), runtime.F.GetOSUsername(runtime.F{}))
	}
}

func (f F) upgrade() {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Checking for devspace updates..", 0)
	exec.WithLiveLogs("devspace", []string{"upgrade"}, false)
}

func (f F) Reinstall() {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Reinstalling DevSpace..", 0)
	if f.setupExists() {
		f.uninstall()
	}
	f.install()
}

func (f F) getReleaseDownloadURL() string {
	var orgaName string = "loft-sh"
	var repoName string = "devspace"
	return "https://github.com/" + orgaName + "/" + repoName + "/releases/download/"
}

func (f F) getReleaseDownloadFileURL() string {
	var fileName string = context + "-" + runtime.F.GetOS(runtime.F{}) + "-" + runtime.F.GetOSArchitecture(runtime.F{})
	return f.getReleaseDownloadURL() + f.getVersion() + "/" + fileName
}

func (f F) getVersion() string {
	return "v5.11.0"
}

func (f F) downloadSetupFile() error {
	var err error
	if !filesystem.Exists(f.getSetupFilePath()) {
		err = io.F.DownloadFile(io.F{}, f.getSetupFilePath(), f.getReleaseDownloadFileURL())
	}
	return err
}

func (f F) uninstall() {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Removing existing devspace installation..", 0)
	filesystem.RemoveFile(f.getInstallPath())
}

func (f F) setupExists() bool {
	return filesystem.Exists(f.getInstallPath())
}
