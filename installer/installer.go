package installer

import (
	"../config"
	"../env"
	"../filesystem"
	"../logging"
	"../runtime"
	"../terminal"
	"../vars"
)

type F struct{}

func (f F) CheckLocalSetup() {
	devConfigStatus := config.Setting("get", "dev", "Spec.Status", "")
	if devConfigStatus != "active" {
		envActive := env.F.GetActive(env.F{}, false)
		envActiveTitle := env.F.GetActive(env.F{}, true)
		workingDir := filesystem.GetWorkingDir()
		if workingDir != f.getOSInstallDir() {
			sel := terminal.GetUserSelection("Do you want to install "+envActiveTitle+"?", []string{}, false, true)
			if sel == "Yes" {
				f.install(workingDir, envActive, envActiveTitle)
				terminal.Exit(0, "")
			}
		}
	}
}

func (f F) getOSInstallDir() string {
	return runtime.F.GetOSInstallDir(runtime.F{}) + env.F.GetActive(env.F{}, false) + "/"
}

func (f F) install(workingDir, envActive, envActiveTitle string) {
	logging.Log([]string{"", vars.EmojiProcess, vars.EmojiInfo}, "Installing "+envActiveTitle+"..", 0)
	logging.Log([]string{"", vars.EmojiProcess, vars.EmojiCrypto}, "You will be asked for temporary admin permissions\nto install "+envActiveTitle+" to "+f.getOSInstallDir(), 0)
	exec := runtime.F.GetRunningExecutable(runtime.F{})
	setupFilePath := workingDir + exec
	installFilePath := f.getOSInstallDir() + exec
	var aborted bool
	sel := terminal.GetUserSelection("Do you want to proceed with the installation process?", []string{}, false, true)
	if sel == "Yes" {
		filesystem.Delete(installFilePath, true)
		if !aborted {
			filesystem.Copy(setupFilePath, installFilePath, true)
			filesystem.GiveProgramPermissions(f.getOSInstallDir(), runtime.F.GetOSUsername(runtime.F{}))
			sel := terminal.GetUserSelection("Do you want to delete the setup file?", []string{}, false, true)
			if sel == "Yes" {
				filesystem.Delete(setupFilePath, false)
			}
			terminal.CreateAlias(envActive, f.getOSInstallDir())
			logging.Log([]string{"", vars.EmojiProcess, vars.EmojiSuccess}, "Installed "+envActiveTitle+".", 0)
		}
	} else {
		aborted = true
	}
	if aborted {
		logging.Log([]string{"", vars.EmojiProcess, vars.EmojiInfo}, "Aborted installation of "+envActiveTitle+".", 0)
	}
}
