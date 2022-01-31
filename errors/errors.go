package errors

import (
	"errors"
	"fmt"
	"os"

	"github.com/neurafuse/tools-go/logging/color"
	"github.com/neurafuse/tools-go/logging/emoji"
	"github.com/neurafuse/tools-go/vars"
)

func Check(err error, callerInfo, msg string, hasFailed, exit, log bool) bool {
	var error bool
	if err != nil || hasFailed {
		error = true
		if log {
			emoji.Println("\n", vars.EmojiError, vars.EmojiInfo, color.Red(callerInfo+": "+msg))
			if err != nil {
				emoji.Println("", vars.EmojiWarning, vars.EmojiInspect, color.Red("Details:"))
				emoji.Println("", vars.EmojiInspect, vars.EmojiInfo, color.Red(err.Error()+"\n"))
			}
		}
		if exit {
			emoji.Println("", vars.EmojiError, vars.EmojiInfo, "Exiting due to unrecoverable error..")
			emoji.Println("\n", vars.EmojiWavingHand, "", "Good bye.")
			fmt.Println()
			os.Exit(1)
		}
	}
	return error
}

func New(msg string) error {
	return errors.New(msg)
}
