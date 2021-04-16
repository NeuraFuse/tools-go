package container

import (
	"../../neurakube/infrastructure/ci/base"
	"../build"
	"../config"
	"../crypto/hash"
	"../env"
	"../errors"
	"../filesystem"
	"../logging"
	"../objects/strings"
	"../runtime"
	"../terminal"
	"../vars"
	"./api/docker"
)

type F struct{}

func (f F) CheckUpdates(context string, push, release bool) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Checking for container image(s) updates for context "+context+"..", 0)
	var packagePaths [][]string
	var envFramework string
	var contextType string
	contextConfig := strings.TrimPrefix(context, "develop/")   // TODO: Refactor
	contextConfig = strings.TrimSuffix(contextConfig, "-base") // TODO: Refactor
	if context != env.F.GetContext(env.F{}, vars.NeuraKubeNameRepo, false) {
		envFramework = config.Setting("get", "infrastructure", "Spec."+strings.Title(contextConfig)+".Environment.Framework", "")
		contextType = config.Setting("get", "infrastructure", "Spec."+strings.Title(contextConfig)+".Type", "")
	}
	templatesBasePath := "../neurakube/infrastructure/ci/"
	dockerfilePath := "/Dockerfile"
	dockerfilePathBase := templatesBasePath + "base/templates/" + envFramework + "/" + contextType + dockerfilePath + "_base"
	neurakubeDockerfilePath := "../neurakube" + dockerfilePath
	dockerfilePathModules := templatesBasePath + strings.TrimSuffix(context, "-base") + "/templates/" + envFramework + "/" + contextType + dockerfilePath
	modulePath := "../lightning"
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
		if context == env.F.GetContext(env.F{}, vars.NeuraKubeNameRepo, false) {
			imageBase = false
			packagePaths = append(packagePaths, []string{templatesBasePath + "api/templates" + dockerfilePath, neurakubeDockerfilePath})
		} else if context == "develop/remote" {
			imageBase = false
			packagePaths = append(packagePaths, []string{dockerfilePathModules, neurakubeDockerfilePath})
		} else if context == "app" || context == "inference" {
			imageBase = false
			packagePaths = append(packagePaths, []string{dockerfilePathModules, neurakubeDockerfilePath})
			packagePaths = append(packagePaths, []string{modulePath, "../neurakube/lightning"})
		} else {
			if !strings.Contains(context, "base") {
				errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create image with unknown context: "+context, true, true, true)
			}
		}
		if !imageBase && !f.checkUpdateStatus(context, "../lightning") {
			success = true
		}
		if !success {
			updateExecuted = true
			imageAddrs := f.GetImgAddrs(contextConfig, imageBase, release)
			f.importPackage(packagePaths)
			if !strings.HasSuffix(context, "-base") {
				_, err := build.F.Make(build.F{}, vars.NeuraKubeNameRepo, "linux", "amd64", "", false, true)
				errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to build "+vars.NeuraKubeNameRepo+"!", false, true, true)
			}
			docker.Initialize(release)
			docker.BuildImage("../"+vars.NeuraKubeNameRepo, []string{imageAddrs})
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
	} else if context == env.F.GetContext(env.F{}, vars.NeuraKubeNameRepo, false) {
		module = context
	} else {
		if context == "remote" {
			module = "develop/" + context + "/"
		} else {
			module = context + "/"
		}
		frameworkAcc = true
	}
	var orgaName string
	var repoAddrs string = f.GetImgRepoAddrs(release)
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

func (f F) GetImgRepoAddrs(release bool) string {
	var repoAddrs string
	var configID string
	var configKey string = "Containers.RepoAddrs"
	if release {
		configID = "dev"
	} else {
		configID = "project"
	}
	if config.ValidSettings(configID, "containers", false) {
		repoAddrs = config.Setting("get", configID, configKey, "")
	} else {
		repoKind := terminal.GetUserSelection("Whith what kind of container repository do you want to interact with?", []string{"Private", "Public (e.g. Docker Hub)"}, false, false)
		repoAddrs = terminal.GetUserInput("What is the base path of the " + strings.ToLower(repoKind) + " container image repository?")
		config.Setting("set", configID, configKey, repoAddrs)
	}
	repoSlash := "/"
	if repoAddrs == vars.OrganizationNameRepo {
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiWarning}, "Interaction with official "+vars.OrganizationNameRepo+" container repository..", 0)
		sel := terminal.GetUserSelection("Do you want to continue?", []string{}, false, true)
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
	contextID := strings.TrimPrefix(context, "develop/") // TODO: Refactor
	contextID = strings.TrimSuffix(contextID, "-base")
	contextID = strings.Title(contextID)
	if strings.HasSuffix(context, "-base") {
		envFramework := base.F.GetEnvFramework(base.F{}, contextID, true)
		resType := base.F.GetType(base.F{}, contextID, true)
		configAddrs := "Build.Accelerator.Base." + envFramework + "." + resType + ".HashRecent"
		hashes = append(hashes, hash.SHA256Folder(filesystem.GetDirPathFromFilePath(path)))
		hashes = append(hashes, config.Setting("get", "dev", configAddrs, ""))
		config.Setting("set", "dev", configAddrs, hashes[0])
	} else if context == env.F.GetContext(env.F{}, vars.NeuraKubeNameRepo, false) {
		build.F.CheckUpdates(build.F{}, "neurakube", false)
		configAddrs := "Build.Neurakube."
		hashes = append(hashes, config.Setting("get", "dev", configAddrs+"HashRecent", ""))
		hashes = append(hashes, config.Setting("get", "dev", configAddrs+"Container.HashRecent", ""))
		config.Setting("set", "dev", configAddrs+"Container.HashRecent", hashes[0])
	} else {
		envFramework := base.F.GetEnvFramework(base.F{}, contextID, true)
		resType := base.F.GetType(base.F{}, contextID, true)
		build.F.CheckUpdates(build.F{}, "neurakube", false)
		configAddrs := "Build.Accelerator." + contextID + "." + envFramework + "."
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