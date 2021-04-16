package emoji

import (
	"fmt"
	"github.com/kyokomi/emoji"
	"../../vars"
	"../color"
)

func Println(prefix, emojiKey1, emojiKey2 string, text string) {
	if prefix != "" {
		fmt.Print(prefix)
	}
	emojiString := ""
	if emojiKey1 != "" {
		emojiString = ":" + emojiKey1 + ":"
	}
	if emojiKey2 != "" {
		suffix := ""
		if emojiKey2 == vars.EmojiSuccess {
			suffix = color.Green(":" + emojiKey2 + ":")
		} else {
			suffix = ":" + emojiKey2 + ":"
		}
		emojiString = emojiString + suffix
	}
	emoji.Println(emojiString + " " + text)
}