package build

import (
	"os"

	"../config"
	buildconfig "../config/build"
	cliconfig "../config/cli"
	devconfig "../config/dev"
	"../crypto/hash"
	dep "../dependencies"
	"../env"
	"../errors"
	"../exec"
	"../filesystem"
	"../logging"
	"../objects/strings"
	"../runtime"
	"../terminal"
	"../updater/golang"
	"../users"
	"../vars"
)

type F struct{}

func (f F) CheckUpdates(module string, handover bool) {
	f.setHandover()
	if !strings.ArrayContains(os.Args, "--"+f.GetFlags()["build"][0]) { // TODO: && checkDo
		checkDo := false // TODO: cobra routing is currently taking place after this check .. better imp. would be: checkDo := buildconfig.F.Setting(buildconfig.F{}, "get", "check", false)
		if env.F.API(env.F{}) {
			if !env.F.Container(env.F{}) {
				checkDo = true
			}
		} else {
			if config.Setting("get", "dev", "Spec.Status", "") == "active" {
				checkDo = true
			}
		}
		if checkDo {
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiInspect}, "Checking "+module+" build..", 0)
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiWarning}, "Don't update files in the "+vars.OrganizationName+" directories while the build process is active.", 0)
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiWarning}, "Don't start concurrent build checks manually.\n", 0)
			hashNow, changed := f.codeAnalysis(module)
			depUpdated := dep.F.CheckBuild(dep.F{})
			if changed || depUpdated {
				logging.Log([]string{"", vars.EmojiDev, vars.EmojiInfo}, "Detected code updates since last build.\n", 0)
				var success bool = false
				for ok := true; ok; ok = !success {
					logging.Log([]string{"", vars.EmojiDev, vars.EmojiProcess}, "Starting rebuild..", 0)
					_, err := f.Make(module, "", "", "", true, false)
					if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Building of new version failed!", false, false, true) {
						success = true
						hashNow, _ = f.codeAnalysis(module)
						f.saveHash(hashNow)
						if handover {
							//args := strings.ArrayRemoveString(os.Args, "./"+moduleExecutable)
							args := []string{"--" + f.GetFlags()["build"][0]}
							args = append(args, "--"+f.GetFlags()["build"][1])
							logging.Log([]string{"\n", vars.EmojiDev, vars.EmojiLink}, "Handover to new build..", 0)
							logging.Log([]string{"", vars.EmojiDev, vars.EmojiInfo}, "Arguments: "+strings.Join(args, " "), 0)
							var err error
							if len(args) == 0 && !env.F.API(env.F{}) {
								logging.Log([]string{"\n", vars.EmojiDev, vars.EmojiWarning}, "To interact with the assistant you have to start the new build manually.", 0)
								logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "Just type "+module+" in your terminal.", 0)
							} else {
								err = exec.WithLiveLogs(env.F.GetActive(env.F{}, false), args, true)
							}
							if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to run program!", false, true, true) {
								terminal.Exit(0, "")
							}
						}
					} else {
						sel := terminal.GetUserSelection("What do you want to do?", []string{"Retry", "Start old build version", "Exit"}, false, false)
						if sel == "Retry" {
							logging.Log([]string{"", vars.EmojiDev, vars.EmojiProcess}, "Retrying to build new version..\n", 0)
						} else if sel == "Start old build version" {
							logging.Log([]string{"", vars.EmojiDev, vars.EmojiWarning}, "Continuing with old build version..\n", 0)
							success = true
						} else if sel == "Exit" {
							terminal.Exit(0, "")
						}
					}
				}
			} else {
				version := f.getVersion(env.F.GetActive(env.F{}, false), false)
				logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "Local build is up to date ("+version+").\n", 0)
			}
		}
	} else {
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "Build check disabled via flag.", 0)
	}
	os.Args = strings.ArrayRemoveString(os.Args, "--"+f.GetFlags()["build"][0])
	os.Args = strings.ArrayRemoveString(os.Args, "--"+f.GetFlags()["build"][1])
}

func (f F) codeAnalysis(module string) (string, bool) {
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiInspect}, "Starting code analysis of module "+module+"..", 0)
	devconfigPath := devconfig.F.GetFilePath(devconfig.F{})
	if env.F.API(env.F{}) {
		filesystem.Move(env.F.GetAPIHTTPCertPath(env.F{}), "../tmp/certs", false)
	} else {
		filesystem.Move(devconfigPath, "../tmp/"+devconfig.FileName, false)
	}
	filesystem.Move("../tools-go", "tools-go", false)
	hashNow := hash.SHA256Folder("../" + module)
	if env.F.API(env.F{}) {
		filesystem.Move("../tmp/certs", env.F.GetAPIHTTPCertPath(env.F{}), false)
	} else {
		filesystem.Move("../tmp/"+devconfig.FileName, devconfigPath, false)
	}
	filesystem.Move("tools-go", "../tools-go", false)
	if env.F.API(env.F{}) {
		cliconfig.F.SetFilePath(cliconfig.F{}, "../"+env.F.GetID(env.F{}, "neuracli")+"/"+users.BasePath)
		userDefault := config.Setting("get", "cli", "Spec.Users.DefaultID", "")
		devconfig.F.SetFilePath(devconfig.F{}, "../"+env.F.GetID(env.F{}, "neuracli")+"/"+users.BasePath, userDefault)
	}
	hashRecent := config.Setting("get", "dev", "Spec.Build."+strings.Title(module)+".HashRecent", "")
	changed := false
	if hashNow != hashRecent {
		changed = true
	}
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "Code analysis finished.", 0)
	return hashNow, changed
}

func (f F) Make(module, targetOS, targetArchitecture, remotePath string, localPath, checkDependencies bool) (string, error) {
	if targetOS == "" {
		targetOS = runtime.F.GetOS(runtime.F{})
	}
	if targetArchitecture == "" {
		targetArchitecture = "amd64"
	}
	emoji := ""
	if module == vars.NeuraCLINameRepo {
		emoji = vars.EmojiClient
	} else if module == vars.NeuraKubeNameRepo {
		emoji = vars.EmojiAPI
	}
	version := f.getVersion(module, true)
	context := " local "
	if remotePath != "" {
		context = " " + remotePath + " "
	}
	logging.Log([]string{"", vars.EmojiProcess, emoji}, "Building "+module+context+"("+version+")..", 0)
	logging.Log([]string{"", vars.EmojiProcess, vars.EmojiInfo}, "Target: "+targetOS+"-"+targetArchitecture, 0)
	goVersion, _, _ := golang.F.GetVersion(golang.F{}, true)
	logging.Log([]string{"", vars.EmojiProcess, vars.EmojiInfo}, "Golang version: "+goVersion, 0)
	os.Setenv("GOOS", targetOS)
	os.Setenv("GOARCH", targetArchitecture)
	buildfile := module
	module = module + "-" + targetOS + "-" + targetArchitecture
	if checkDependencies {
		dep.F.CheckBuild(dep.F{})
	}
	logging.ProgressSpinner("start")
	/*if depUpdate { TODO: Reactivate if (depUpdate bool) no longer needed (same global signatures for dynamic func selection)
		logging.Log([]string{"", emoji, vars.EmojiInfo}, "Updating all dependencies (this may take a while)..", 0)
		exec.WithLiveLogs("go", "get -u all")
	}*/
	err := exec.WithLiveLogs("go", []string{"build", "-o", module, "../" + buildfile + "/" + buildfile + ".go"}, true)
	logging.ProgressSpinner("stop")
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to build program!", false, false, true) {
		buildPath := ""
		if localPath {
			buildPath = module
		} else if remotePath != "" {
			buildPath = remotePath + module
			buildDir := filesystem.GetDirPathFromFilePath(buildPath)
			if filesystem.Exists(buildDir) {
				if filesystem.Exists(buildPath) {
					filesystem.Delete(buildPath, false)
				}
			} else {
				filesystem.CreateDir(buildDir, false)
			}
			filesystem.Copy(module, buildPath, false)
			hash := hash.SHA256File(buildPath)
			filesystem.AppendStringToFile(buildPath+".sha256", hash)
			filesystem.Delete(module, false)
		} else {
			buildPath = strings.Split(buildfile, ".go")[0]
			buildPath = "../" + buildPath + "/" + buildPath
			if filesystem.Exists(buildPath) {
				filesystem.Delete(buildPath, false)
			}
			filesystem.Copy(module, buildPath, false)
			filesystem.Delete(module, false)
		}
		logging.Log([]string{"", vars.EmojiDir, vars.EmojiSuccess}, "Saved executable at: {workdir}/"+buildPath, 0)
		logging.Log([]string{"", vars.EmojiProcess, vars.EmojiSuccess}, "Build successfull.\n", 0)
	}
	return module, err
}

func (f F) saveHash(hash string) {
	config.Setting("set", "dev", "Spec.Build."+env.F.GetActive(env.F{}, true)+".HashRecent", hash)
}

func (f F) getVersion(module string, next bool) string {
	idRecent := config.Setting("get", "dev", "Spec.Build."+strings.Title(module)+".IDRecent", "")
	version := config.Setting("get", "dev", "Spec.Build."+strings.Title(module)+".Version", "")
	if idRecent != "" {
		if next {
			idRecent = strings.ToString(strings.ToInt(idRecent) + 1)
		}
	} else {
		idRecent = "0"
	}
	if next {
		config.Setting("set", "dev", "Spec.Build."+strings.Title(module)+".IDRecent", idRecent)
	}
	return version + idRecent
}

func (f F) GetFlags() map[string][]string {
	flags := make(map[string][]string)
	flags["build"] = []string{"build-check-disable", "build-handover"}
	return flags
}

func (f F) setHandover() {
	buildconfig.F.Setting(buildconfig.F{}, "set", "handover", strings.ArrayContains(os.Args, "--"+f.GetFlags()["build"][1])) // TODO: [shell perspective] replace array contains check with bool var buildHandover (set setting works but bool does not switch at flag occurence)
}
