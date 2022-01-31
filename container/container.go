package container

import (
	"github.com/neurafuse/tools-go/build"
	"github.com/neurafuse/tools-go/ci/base"
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/container/api/docker"
	"github.com/neurafuse/tools-go/crypto/hash"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) GetVolumeMountPath() string {
	var projectName string = config.Setting("get", "project", "Metadata.Name", "")
	var containerAppDataRoot string = config.Setting("get", "project", "Spec.Containers.Sync.PathMappings.ContainerAppDataRoot", "")
	return env.F.GetContainerWorkingDir(env.F{}) + "/" + projectName + "/" + containerAppDataRoot
}

func (f F) CheckUpdates(context string, push, release bool) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Checking for container image(s) updates for context "+context+"..", 0)
	var packagePaths [][]string
	var envFramework string
	var contextType string
	var contextConfig string = strings.TrimPrefix(context, "develop/") // TODO: Refactor
	contextConfig = strings.TrimSuffix(contextConfig, "-base")         // TODO: Refactor
	if context != env.F.GetContext(env.F{}, vars.NeuraKubeNameID, false) {
		envFramework = config.Setting("get", "infrastructure", "Spec."+strings.Title(contextConfig)+".Environment.Framework", "")
		contextType = config.Setting("get", "infrastructure", "Spec."+strings.Title(contextConfig)+".Type", "")
	}
	var templatesBasePath string = "../tools-go@" + vars.ToolsGoVersion + "/ci/"
	var dockerfilePath string = "/Dockerfile"
	var dockerfilePathBase string = templatesBasePath + "base/templates/" + envFramework + "/" + contextType + dockerfilePath + "_base"
	var neurakubeDockerfilePath string = "../neurakube@" + vars.NeuraKubeVersion + "/" + dockerfilePath
	var dockerfilePathModules string = templatesBasePath + strings.TrimSuffix(context, "-base") + "/templates/" + envFramework + "/" + contextType + dockerfilePath
	var modulePath string = "../lightning"
	var imageBase bool
	var success bool
	var updateExecuted bool
	for ok := true; ok; ok = !success {
		if strings.HasSuffix(context, "-base") {
			if f.checkUpdateStatus(context, dockerfilePathBase) {
				packagePaths = append(packagePaths, []string{dockerfilePathBase, neurakubeDockerfilePath})
				imageBase = true
			} else {
				context = strings.TrimSuffix(context, "-base")
			}
		}
		if context == env.F.GetContext(env.F{}, vars.NeuraKubeNameID, false) {
			imageBase = false
			packagePaths = append(packagePaths, []string{templatesBasePath + "api/templates" + dockerfilePath, neurakubeDockerfilePath})
		} else if context == "develop" {
			imageBase = false
			packagePaths = append(packagePaths, []string{dockerfilePathModules, neurakubeDockerfilePath})
		} else if context == "app" || context == "inference" {
			imageBase = false
			packagePaths = append(packagePaths, []string{dockerfilePathModules, neurakubeDockerfilePath})
			packagePaths = append(packagePaths, []string{modulePath, "../neurakube@" + vars.NeuraKubeVersion + "/lightning"})
		} else {
			if !strings.Contains(context, "base") {
				errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create image with unknown context: "+context, true, true, true)
			}
		}
		if !success {
			updateExecuted = true
			var imageAddrs string = f.GetImgAddrs(contextConfig, imageBase, release)
			f.importPackage(packagePaths)
			if !strings.HasSuffix(context, "-base") {
				_, err := build.F.Make(build.F{}, vars.NeuraKubeNameID, "linux", "amd64", "", false, true)
				errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to build "+vars.NeuraKubeNameID+"!", false, true, true)
			}
			docker.Initialize(release)
			docker.BuildImage("../"+vars.NeuraKubeNameID, []string{imageAddrs})
			f.removePackage(packagePaths)
			if push {
				docker.PushImage(imageAddrs)
			}
			if strings.HasSuffix(context, "-base") {
				context = strings.TrimSuffix(context, "-base")
			} else {
				success = true
			}
		}
	}
	var updateDone string = "already"
	if updateExecuted {
		updateDone = "now"
	}
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiSuccess}, "All container image(s) for context "+context+" are "+updateDone+" prepared.\n", 0)
}

func (f F) GetImgAddrs(context string, baseIMG, release bool) string {
	var baseSuffix string
	var module string
	var frameworkAcc bool
	var imageAddrs string
	if baseIMG {
		baseSuffix = "/base"
		module = ""
		frameworkAcc = true
	} else if context == env.F.GetContext(env.F{}, vars.NeuraKubeNameID, false) {
		module = context
	} else {
		module = context + "/"
		frameworkAcc = true
	}
	var orgaName string
	var repoAddrs string = f.GetRegistryAddress(release)
	if !f.isPublicRepoAddrs(repoAddrs) {
		orgaName = vars.OrganizationNameRepo + "/"
	}
	if frameworkAcc {
		frameworkAccAddrs := base.F.GetEnvFramework(base.F{}, context, false) + "/" + base.F.GetType(base.F{}, context, false) + baseSuffix
		imageAddrs = orgaName + module + frameworkAccAddrs + ":latest"
	} else {
		imageAddrs = orgaName + module + ":latest"
	}
	return repoAddrs + imageAddrs
}

func (f F) GetRegistryAddress(release bool) string {
	var repoAddrs string
	var configID string
	var configKey string = "Spec.Containers.Registry.Address"
	if release {
		configID = "dev"
	} else {
		configID = "project"
	}
	if config.ValidSettings(configID, "containers/registry", false) {
		repoAddrs = config.Setting("get", configID, configKey, "")
	} else {
		var repoKind string = terminal.GetUserSelection("Whith what kind of container repository do you want to interact with?", []string{"Private", "Public (e.g. Docker Hub)"}, false, false)
		repoAddrs = terminal.GetUserInput("What is the base path of the " + strings.ToLower(repoKind) + " container image repository?")
		config.Setting("set", configID, configKey, repoAddrs)
	}
	var repoSlash string = "/"
	if repoAddrs == vars.OrganizationNameRepo {
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiWarning}, "Interaction with official "+vars.OrganizationNameRepo+" container repository..", 0)
		var sel string = terminal.GetUserSelection("Do you want to continue?", []string{}, false, true)
		if sel == "No" {
			config.Setting("set", configID, configKey, "")
			terminal.Exit(0, "")
		}
		repoSlash = ""
	}
	return repoAddrs + repoSlash
}

func (f F) isPublicRepoAddrs(addrs string) bool {
	var public bool
	if !strings.Contains(addrs, ".") {
		public = true
	}
	return public
}

func (f F) importPackage(packagePaths [][]string) {
	for i, _ := range packagePaths {
		if filesystem.Exists(packagePaths[i][1]) {
			filesystem.Delete(packagePaths[i][1], false)
		}
		filesystem.Copy(packagePaths[i][0], packagePaths[i][1], false)
	}
}

func (f F) removePackage(packagePaths [][]string) {
	for i, _ := range packagePaths {
		if filesystem.Exists(packagePaths[i][1]) {
			filesystem.Delete(packagePaths[i][1], false)
		}
	}
}

func (f F) checkUpdateStatus(context, path string) bool {
	var update bool
	var hashes []string
	var contextID string = strings.TrimPrefix(context, "develop/") // TODO: Refactor
	contextID = strings.TrimSuffix(contextID, "-base")
	contextID = strings.Title(contextID)
	if strings.HasSuffix(context, "-base") {
		var envFramework string = base.F.GetEnvFramework(base.F{}, contextID, true)
		var resType string = base.F.GetType(base.F{}, contextID, true)
		var configAddrs string = "Spec.Build.Accelerator.Base." + envFramework + "." + resType + ".HashRecent"
		hashes = append(hashes, hash.SHA256Folder(filesystem.GetDirPathFromFilePath(path)))
		hashes = append(hashes, config.Setting("get", "dev", configAddrs, ""))
		config.Setting("set", "dev", configAddrs, hashes[0])
	} else if context == env.F.GetContext(env.F{}, vars.NeuraKubeNameID, false) {
		build.F.CheckUpdates(build.F{}, "neurakube", false)
		var configAddrs string = "Spec.Build.Neurakube."
		hashes = append(hashes, config.Setting("get", "dev", configAddrs+"HashRecent", ""))
		hashes = append(hashes, config.Setting("get", "dev", configAddrs+"Container.HashRecent", ""))
		config.Setting("set", "dev", configAddrs+"Container.HashRecent", hashes[0])
	} else {
		var envFramework string = base.F.GetEnvFramework(base.F{}, contextID, true)
		var resType string = base.F.GetType(base.F{}, contextID, true)
		build.F.CheckUpdates(build.F{}, "neurakube", false)
		var configAddrs string = "Spec.Build.Accelerator." + contextID + "." + envFramework + "."
		hashes = append(hashes, config.Setting("get", "dev", "Spec.Build.Neurakube.HashRecent", ""))
		hashes = append(hashes, config.Setting("get", "dev", configAddrs+resType+".HashRecent", ""))
		config.Setting("set", "dev", configAddrs+resType+".HashRecent", hashes[0])
		hashes = append(hashes, hash.SHA256Folder(path))
		hashes = append(hashes, config.Setting("get", "dev", configAddrs+resType+".Lightning.HashRecent", ""))
		config.Setting("set", "dev", configAddrs+resType+".Lightning.HashRecent", hashes[2])
	}
	for i, hash := range hashes {
		if i+1 != len(hashes) {
			if hash != hashes[i+1] {
				update = true
				logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInspect}, "Detected an update for container image "+context+".", 0)
			}
		}
	}
	return update
}
