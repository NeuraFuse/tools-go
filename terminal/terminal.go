package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	buildconfig "../config/build"
	"../env"
	"../errors"
	"../filesystem"
	"../logging"
	"../logging/color"
	"../objects/strings"
	"../runtime"
	"../timing"
	"../vars"
	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"
)

func Init(skipIntro bool) {
	listenOSInterrupt()
	intro(skipIntro)
}

func UserSelectionFiles(label, objectType, filePath string, search, exclude []string, withAdd bool, returnAbsolutePath bool) string {
	var success bool = false
	var hintlogged bool = false
	var fileList []string
	if !filesystem.Exists(filePath) {
		filesystem.CreateDir(filePath, false)
	}
	for ok := true; ok; ok = !success {
		fileList = filesystem.Explorer(objectType, filePath, search, exclude)
		if len(fileList) != 0 {
			success = true
		} else {
			if !hintlogged {
				logging.Log([]string{"\n", vars.EmojiDir, vars.EmojiInfo}, "No "+strings.Join(search, ", ")+" file found in: "+filePath, 0)
				logging.Log([]string{"", vars.EmojiDir, vars.EmojiInfo}, "Please move a suitable file to this location.", 0)
				logging.Log([]string{"", vars.EmojiDir, vars.EmojiInfo}, "Auto refreshing every second..", 0)
				hintlogged = true
			}
			timing.TimeOut(200, "ms")
		}
	}
	var absolutePath string = ""
	if returnAbsolutePath {
		absolutePath = filePath + "/"
	}
	return absolutePath + GetUserSelection(label, fileList, withAdd, false)
}

func GetUserSelection(label string, options []string, withAdd, binaryOptions bool) string {
	// label = "assistant: " + label
	/*template := promptui.SelectTemplates {
		Label:    "{{ . }}?",
	}*/
	if binaryOptions {
		options = []string{"Yes", "No"}
	}
	var result string
	var err error
	if withAdd {
		prompt := promptui.SelectWithAdd{
			Label:    label,
			Items:    options,
			AddLabel: "Create new value for config",
		}
		_, result, err = prompt.Run()
		fmt.Println()
		if strings.HasPrefix(result, "Default: ") {
			result = strings.ReplaceAll(result, "Default: ", "")
			if result == "blank" {
				result = ""
			}
		}
		if strings.HasPrefix(result, "Example: ") {
			result = strings.ReplaceAll(result, "Example: ", "")
		}
	} else {
		prompt := promptui.Select{
			Label: label,
			Items: options,
		}
		_, result, err = prompt.Run()
		fmt.Println()
		if strings.HasPrefix(result, "Default: ") {
			result = strings.ReplaceAll(result, "Default: ", "")
			if result == "blank" {
				result = ""
			}
		}
	}
	if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, false) {
		logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "Unable to get user selection via prompt!", 0)
		if buildconfig.F.Setting(buildconfig.F{}, "get", "handover", false) {
			logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "This occurs if the assistant gets called after a build handover.", 0)
			logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "Please restart manually to workaround.\n", 0)
		}
		Exit(0, "")
	}
	return result
}

func GetUserInput(instruction string) string {
	reader := bufio.NewReader(os.Stdin)
	logging.Log([]string{"", vars.EmojiInfo, ""}, instruction, 0)
	input, err := reader.ReadString('\n')
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get user input", false, true, true)
	input = strings.TrimSuffix(input, "\n")
	return input
}

func intro(skip bool) {
	if !skip {
		module := env.F.GetActive(env.F{}, false)
		var asciiArtColor string
		var licenseType string
		var footer string
		if module == vars.NeuraCLINameRepo {
			module = vars.NeuraCLIName
			asciiArtColor = "cyan"
			licenseType = vars.NeuraCLILicenseType
			footer = "    " + vars.OrganizationName + " | " + vars.OrganizationWebsite + " | © " + licenseType+" ("+env.F.GetVersion(env.F{})+")"
		} else if module == vars.NeuraKubeNameRepo {
			module = vars.NeuraKubeName
			asciiArtColor = "cyan"
			licenseType = vars.NeuraKubeLicenseType
			footer = "       " + vars.OrganizationName + " | " + vars.OrganizationWebsite + " | © " + licenseType+" ("+env.F.GetVersion(env.F{})+")"
		}
		fmt.Println()
		printASCIIArt(module, "small", asciiArtColor)
		//year := strings.ToString(timing.GetCurrent("year"))
		fmt.Println(color.Red(footer))
		logging.PartingLine()
	} else {
		fmt.Println("")
	}
}

func printASCIIArt(text string, style string, color string) {
	art := figure.NewColorFigure(text, style, color, true)
	art.Print()
}

func Exit(errorCode int, msg string) {
	if !buildconfig.F.Setting(buildconfig.F{}, "get", "handover", false) {
		var emoji string
		var newline string
		if errorCode == 0 {
			emoji = vars.EmojiWavingHand
			msg = "Everything done."
		} else if errorCode == 1 {
			newline = "\n"
			emoji = "upside_down_face"
			msg = "Unfortunately an unexpected error occured on our side. Sorry for the inconvenience.\n"
		} else if errorCode == 2 {
			newline = "\n"
			emoji = vars.EmojiWarning
			errorCode = 0
		}
		logging.Log([]string{newline, emoji, ""}, msg+"\n", 0)
	}
	os.Exit(errorCode)
}

func listenOSInterrupt() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		Exit(0, "")
	}()
}

func getShell() string {
	os := runtime.F.GetOS(runtime.F{})
	var shell string
	switch os {
	case "linux":
		shell = "bash"
	case "darwin":
		shell = "zsh"
	}
	return shell
}

func CreateAlias(envID, path string) {
	var shellConfigPath string = filesystem.GetUserHomeDir() + "/." + getShell() + "rc"
	var cmds []string
	if envID == vars.NeuraCLINameRepo || envID == vars.NeuraKubeNameRepo || envID == "go" {
		os := runtime.F.GetOS(runtime.F{})
		archi := runtime.F.GetOSArchitecture(runtime.F{})
		executable := envID
		if envID != "go" {
			executable = envID+"-"+os+"-"+archi
		}
		cmds = []string{"\nalias " + envID + "='cd " + path + " && ./" + executable +"'", "export PATH=\"" + path + ":$PATH\""}
	} else {
		cmds = []string{"\nalias " + envID + "='" + path + "'", "export PATH=\"" + path + ":$PATH\""}
	}
	created := false
	for _, cmd := range cmds {
		if !filesystem.FileContainsString(shellConfigPath, cmd) {
			filesystem.AppendStringToFile(shellConfigPath, cmd)
			created = true
		}
	}
	if created {
		logging.Log([]string{"\n", vars.EmojiProcess, vars.EmojiInfo}, "Created terminal alias "+envID+".", 0)
		logging.Log([]string{"", vars.EmojiProcess, vars.EmojiInfo}, "You can now start "+envID+" by just typing "+envID+" in your terminal.", 0)
		logging.Log([]string{"", vars.EmojiProcess, vars.EmojiInfo}, "Please restart your terminal for this change to take effect.\n", 0)
	}
}