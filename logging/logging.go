package logging

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cheggaaa/pb/v3"
	devConfig "github.com/neurafuse/tools-go/config/dev"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/logging/color"
	"github.com/neurafuse/tools-go/logging/emoji"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/timing"
)

var logMsgLast string
var LogTimeLast time.Time

func Log(style []string, msg string, level int) {
	var print bool
	if env.F.CLI(env.F{}) {
		devConfig.F.SetConfig(devConfig.F{})
	} else {
		level = 0
	}
	if msg != logMsgLast {
		if level == 0 {
			print = true
		} else if level == 1 {
			if devConfig.F.IsLogLevelActive(devConfig.F{}, "info") || devConfig.F.IsLogLevelActive(devConfig.F{}, "debug") {
				print = true
			}
		} else if level == 2 {
			if devConfig.F.IsLogLevelActive(devConfig.F{}, "debug") {
				print = true
			}
		}
		if print {
			emoji.Println(style[0], style[1], style[2], msg)
		}
	}
	logMsgLast = msg
	LogActive()
}

func LogActive() {
	LogTimeLast = timing.GetCurrentTime()
}

var progressSpinner *spinner.Spinner = spinner.New(spinner.CharSets[11], timing.GetTimeDuration(50, "ms"))

func ProgressSpinner(action string) {
	progressSpinner.Color("green", "bold")
	if action == "start" {
		go psController(action)
	} else if action == "stop" {
		psController(action)
	}
}

var psControllerInterrupt bool

func psController(action string) {
	secs := 1
	switch action {
	case "start":
		var started bool
		for {
			if !psControllerInterrupt {
				if secs == 1 {
					if !started {
						progressSpinner.Start()
						started = true
					}
				}
				if secs > 1 {
					progressSpinner.Suffix = color.Green(" Fusing it together.. (" + strings.ToString(secs) + "s)")
				}
				timing.Sleep(1, "s")
				secs++
			} else {
				break
			}
		}
	case "stop":
		psControllerInterrupt = true
		progressSpinner.Stop()
	}
}

func GetProgressBar() *pb.ProgressBar {
	fmt.Println()
	tmpl := `{{ yellow "Starting assistant:" }} {{ bar . "|" "-" (cycle . "↖" "↗" "↘" "↙" ) "." "|"}} {{speed . | yellow }} {{percent .}} {{string . "info_1" | green}} {{string . "info_2" | blue}}` // rndcolor
	var limit int64 = 512 * 512 * 200
	cli_status := pb.ProgressBarTemplate(tmpl).Start64(limit)
	cli_status.SetWidth(40)
	cli_status.SetRefreshRate(timing.GetTimeFormat("s"))
	cli_status.Set("info_1", "").Set("info_2", "")
	fmt.Println()
	return cli_status
}

func PartingLine() {
	var line string
	if env.F.CLI(env.F{}) {
		line = "_____________________________________________________"
	} else if env.F.API(env.F{}) {
		line = "____________________________________________________________"
	}
	fmt.Println(line + "\n")
}
