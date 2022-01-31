package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"
	buildconfig "github.com/neurafuse/tools-go/config/build"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging/color"
	"github.com/neurafuse/tools-go/logging/emoji"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

func Init(skipIntro bool) {
	listenOSInterrupt()
	intro(skipIntro)
}

func UserSelectionFiles(label, objectType, filePath string, search, exclude []string, withAdd bool, returnAbsolutePath bool) string {
	var success bool
	var hintlogged bool
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
				emoji.Println("\n", vars.EmojiDir, vars.EmojiInfo, "No "+strings.Join(search, ", ")+" file found in: "+filePath)
				emoji.Println("", vars.EmojiDir, vars.EmojiInfo, "Please move a suitable file to this location.")
				emoji.Println("", vars.EmojiDir, vars.EmojiInfo, "Auto refreshing every second..")
				hintlogged = true
			}
			timing.Sleep(200, "ms")
		}
	}
	var absolutePath string
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
		emoji.Println("", vars.EmojiWarning, vars.EmojiInfo, "Unable to get user selection via prompt!")
		if buildconfig.F.Setting(buildconfig.F{}, "get", "handover", false) {
			emoji.Println("", vars.EmojiWarning, vars.EmojiInfo, "This occurs if the assistant gets called after a build handover.")
			emoji.Println("", vars.EmojiWarning, vars.EmojiInfo, "Please restart manually to workaround.\n")
		}
		Exit(0, "")
	}
	return result
}

func GetUserInput(instruction string) string {
	reader := bufio.NewReader(os.Stdin)
	emoji.Println("", vars.EmojiInfo, "", instruction)
	input, err := reader.ReadString('\n')
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get user input", false, true, true)
	input = strings.TrimSuffix(input, "\n")
	return input
}

func intro(skip bool) {
	if !skip {
		var module string = env.F.GetActive(env.F{}, false)
		var asciiArtColor string
		var licenseType string
		var footer string
		if module == vars.NeuraCLINameID {
			module = vars.NeuraCLIName
			asciiArtColor = "cyan"
			licenseType = vars.NeuraCLILicenseType
			footer = "     " + vars.OrganizationName + " | " + vars.OrganizationWebsite + " | © " + licenseType + " (" + env.F.GetVersion(env.F{}) + ")"
		} else if module == vars.NeuraKubeNameID {
			module = vars.NeuraKubeName
			asciiArtColor = "cyan"
			licenseType = vars.NeuraKubeLicenseType
			footer = "        " + vars.OrganizationName + " | " + vars.OrganizationWebsite + " | © " + licenseType + " (" + env.F.GetVersion(env.F{}) + ")"
		}
		fmt.Println()
		printASCIIArt(module, "small", asciiArtColor)
		//year := strings.ToString(timing.GetCurrent("year"))
		fmt.Println(color.Red(footer))
		partingLine()
	} else {
		fmt.Println("")
	}
}

func partingLine() {
	var line string
	if env.F.CLI(env.F{}) {
		line = "_____________________________________________________"
	} else if env.F.API(env.F{}) {
		line = "____________________________________________________________"
	}
	fmt.Println(line + "\n")
}

func printASCIIArt(text string, style string, color string) {
	art := figure.NewColorFigure(text, style, color, true)
	art.Print()
}

func Exit(errorCode int, msg string) {
	if !buildconfig.F.Setting(buildconfig.F{}, "get", "handover", false) {
		var emojiKey string
		var newline string = "\n"
		if errorCode == 0 {
			emojiKey = vars.EmojiWavingHand
			msg = "Everything done."
		} else if errorCode == 1 {
			emojiKey = "upside_down_face"
			msg = "Unfortunately an unexpected error occured.\n"
		} else if errorCode == 2 {
			emojiKey = vars.EmojiWarning
			errorCode = 0
		}
		emoji.Println(newline, vars.EmojiAssistant, emojiKey, msg+"\n")
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
	cmds = []string{"export PATH=$PATH:" + path}
	var created bool
	for _, cmd := range cmds {
		if !filesystem.FileContainsString(shellConfigPath, cmd) {
			filesystem.AppendStringToFile(shellConfigPath, cmd)
			created = true
		}
	}
	if created {
		emoji.Println("\n", vars.EmojiProcess, vars.EmojiInfo, "Created terminal alias "+envID+".")
		emoji.Println("", vars.EmojiProcess, vars.EmojiInfo, "You can now start "+envID+" by just typing "+envID+" in your terminal.")
		emoji.Println("", vars.EmojiProcess, vars.EmojiInfo, "Please restart your terminal for this change to take effect.\n")
	}
}
