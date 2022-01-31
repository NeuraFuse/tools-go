package emoji

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/neurafuse/tools-go/logging/color"
	"github.com/neurafuse/tools-go/vars"
)

func Println(prefix, emojiKey1, emojiKey2 string, text string) {
	if prefix != "" {
		fmt.Print(prefix)
	}
	var emojiString string
	if emojiKey1 != "" {
		emojiString = ":" + emojiKey1 + ":"
	}
	if emojiKey2 != "" {
		var suffix string
		if emojiKey2 == vars.EmojiSuccess {
			suffix = color.Green(":" + emojiKey2 + ":")
		} else {
			suffix = ":" + emojiKey2 + ":"
		}
		emojiString = emojiString + suffix
	}
	emoji.Println(emojiString + " " + text)
}
