package devspace

import (
	"../../../config"
	"../../../env"
	"../../../errors"
	"../../../exec"
	"../../../filesystem"
	"../../../io"
	"../../../logging"
	"../../../projects"
	"../../../readers/yaml"
	"../../../runtime"
	"../../../vars"
	"./templates"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Sync(context, namespace, imageAddrs string, thread bool) {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Preparing file synchronization..", 0)
	//devspace.Sync(createDevSpaceConfig(templateBasePath, templatePath))
	configPath := f.createDevSpaceConfig(context, namespace, imageAddrs)
	f.syncExec(configPath, thread)
}

func (f F) createDevSpaceConfig(context, namespace, imageAddrs string) string {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Creating devspace configs..", 0)
	configInt := f.getConfigInterface(context, f.getTemplateBasePath(), f.getPath("template"), namespace, imageAddrs)
	projectConfigPath := f.getPath("projectconfig")
	if filesystem.Exists(projectConfigPath) {
		filesystem.Delete(projectConfigPath, false)
	}
	yaml.StructToFile(projectConfigPath, configInt)
	return projectConfigPath
}

func (f F) getConfigInterface(context, templateBasePath, templatePath, namespace, imageAddrs string) interface{} {
	dConfig := yaml.FileToStruct(templateBasePath+templatePath, &templates.Default{})
	repoPreamble := config.Setting("get", "dev", "Spec.Containers.Registry.Address", "")
	dConfig.(*templates.Default).Images.Remote.Image = repoPreamble + imageAddrs
	dConfig.(*templates.Default).Dev.Sync[0].ImageName = context
	dConfig.(*templates.Default).Dev.Sync[0].Namespace = namespace
	dConfig.(*templates.Default).Dev.Sync[0].ExcludePaths = []string{"lost+found"}
	return dConfig
}

func (f F) getBasePath() string {
	var basePath string = "../tools-go/apps/kubernetes/devspace/"
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
			path = projects.F.GetWorkingDir(projects.F{}) + context + "/" + context + ext
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
		filesystem.Move(f.getSetupFilePath(), f.getInstallPath(), false)
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
