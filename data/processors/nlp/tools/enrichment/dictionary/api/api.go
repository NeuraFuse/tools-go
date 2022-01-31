package api

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/json"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

type Response [1]struct {
	Word      string `json:"word"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string   `json:"definition"`
			Example    string   `json:"example"`
			Synonyms   []string `json:"synonyms,omitempty"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func (f F) GetDefinition(word string) []string {
	url := "https://api.dictionaryapi.dev/api/v2/entries/en_US/" + word
	response := new(Response)
	json.URLToInterface(url, response)
	var defAr []string
	if len(response[0].Meanings) != 0 {
		definition := response[0].Meanings[0].Definitions[0].Definition + "."
		defAr = strings.FieldsFunc(definition, f.Split)
	} else {
		logging.Log([]string{"", vars.EmojiInspect, vars.EmojiWarning}, "Unable to get definition for word: "+word, 0)
	}
	return defAr
}

func (f F) Split(r rune) bool {
	return r == ' ' || r == '.' || r == ',' || r == '\n'
}
