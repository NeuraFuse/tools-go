package filter

import (
	"regexp"
	"unicode"

	"../../errors"
	"../../objects/strings"
	"../../runtime"
	"github.com/microcosm-cc/bluemonday"
)

func SentenceMaxWords(text string, max int) string {
	if getSentenceMaxWords(text) > max {
		return ""
	} else {
		return text
	}
}

func getSentenceMaxWords(text string) int {
	maxWords, nWords := 0, 0
	inWord := false
	for _, r := range text {
		switch r {
		case '.', '?', '!':
			inWord = false
			if maxWords < nWords {
				maxWords = nWords
			}
			nWords = 0
		default:
			if unicode.IsSpace(r) {
				inWord = false
			} else if inWord == false {
				inWord = true
				nWords++
			}
		}
		if maxWords < nWords {
			maxWords = nWords
		}
	}
	return maxWords
}

func WordsLength(text, mode string, limit int) string { // TODO: only mode min not working!
	var output string
	sentences := strings.ToSentences(text)
	for _, sentence := range sentences {
		words := strings.SentenceToWords(sentence)
		for _, word := range words {
			if mode == "min" {
				if strings.GetWordLength(word) >= limit {
					output = output + word
				}
			} else if mode == "max" {
				if strings.Count(word) <= limit {
					output = output + word
				}
			} else {
				errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Imcompatible mode: "+mode, true, false, true)
			}
		}
	}
	return output
}

func UnrealWords(text string) string { // TODO: Refactor loops @WordsLength
	var output string
	sentences := strings.ToSentences(text)
	for _, sentence := range sentences {
		words := strings.SentenceToWords(sentence)
		for _, word := range words {
			save := true
			var arLetter []rune
			for _, letter := range word {
				arLetter = append(arLetter, letter)
			}
			for _, letter := range arLetter {
				iOccurence := 0
				iPuncts := 0
				for _, letterI := range arLetter {
					if letter == letterI {
						iOccurence++
					}
					if unicode.IsPunct(letterI) {
						iPuncts++
					}
				}
				if iOccurence > 2 { // Has letters that are occuring more than if times
					save = false
				}
				if iPuncts > 1 { // Word contains puncts more than if times
					save = false
				}
			}
			if save {
				output = output + word
			}
		}
	}
	return output
}

func Casing(text string) string {
	return strings.ToLower(text)
}

func Digits(text string) string {
	re := regexp.MustCompile("[0-9]+")
	return re.ReplaceAllString(text, "")
}

func SpecialChars(text string, exceptions []string) string {
	//text = RemoveAlphaNumeric(text, exceptions)
	text = RemoveAlphaNumericRegex(text, exceptions)
	//text = RemoveExcept(text, exceptions)
	return text
}

func RemoveExcept(text string, exceptions []string) string {
	var output string
	for _, char := range text {
		var save bool
		for _, exception := range exceptions {
			for _, letter := range exception {
				if char == letter {
					save = true
				}
			}
		}
		if save {
			output = output + string(char)
		}
	}
	return output
}

func RemoveAlphaNumericRegex(text string, exceptions []string) string {
	r := strings.Join(exceptions, "")
	re := regexp.MustCompile("[" + r + "]+")
	output := re.ReplaceAllString(text, "")
	output = strings.Replace(output, "-", "", -1)
	return RemoveChineseChars(output)
}

func RemoveChineseChars(text string) string {
	re := regexp.MustCompile(`\p{Han}*`)
	return re.ReplaceAllString(text, "")
}

func DuplicateLines(text, contentNew, sep string) string {
	//textSent := strings.ToSentences(text)
	//contentNewSent := strings.ToSentences(contentNew)
	textSent := strings.Split(text, sep)
	contentNewSent := strings.Split(contentNew, sep)
	/*
	var contentNewFiltered []string
	for _, sentenceNew := range contentNewSent {
		//if !strings.ArrayContains(contentNewSent, sentence) { // TODO: Abstract via ArrayArrayContains
		var add bool = true
		for _, sentence := range textSent {
			if sentence == sentenceNew {
				add = false
			}
			//contentNewSent = strings.ArrayRemoveString(contentNewSent, sentence)
		}
		if add {
			contentNewFiltered = append(contentNewFiltered, sentenceNew)
		}
		textSent = append(textSent, contentNewFiltered...)
	}*/
	textSent = append(textSent, contentNewSent...)
	return strings.Join(unique(textSent), sep)
}

func unique(stringSlice []string) []string {
    keys := make(map[string]bool)
    list := []string{} 
    for _, entry := range stringSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}

func TrimSpaceMultiline(text, sep string) string {
	lines := strings.ToArray(text, sep)
	output := ""
	for i, line := range lines {
		if i != 0 {
			output = output + "\n"
		}
		output = output + strings.TrimSpace(line)
	}
	output = strings.Replace(output, "  ", "", -1)
	return output
}

func EmptyLines(text string) string {
	re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	text = re.ReplaceAllString(text, "")
	return text
}

func Code(text, codeType string) string {
	switch codeType {
	case "html":
		p := bluemonday.StrictPolicy()
		text = p.Sanitize(text)
	}
	return text
}