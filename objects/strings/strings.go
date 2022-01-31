package strings

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
)

type Builder struct {
	// contains filtered or unexported fields
}

func ToBytes(input string) []byte {
	return []byte(input)
}

func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func BoolToString(input bool) string {
	return strconv.FormatBool(input)
}

func LinesCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

func SentenceToWords(s string) []string {
	var words []string
	var word string
	iWords := 0
	for iL, letter := range s {
		word = word + RuneToString(letter)
		if unicode.IsSpace(letter) {
			words = append(words, word)
			iWords++
			word = ""
		} else if unicode.IsPunct(letter) {
			for iP, letterP := range s {
				if iP == iL+1 {
					if unicode.IsSpace(letterP) || letterP == '\n' {
						words = append(words, word)
						iWords++
						word = ""
					}
				}
			}
		}
		/*if unicode.IsSpace(letter) || unicode.IsPunct(letter) {
			words = append(words, word)
			iWords++
			word = ""
		}*/
	}
	return words
}

func FieldsFunc(input string, f func(rune) bool) []string {
	return strings.FieldsFunc(input, f)
}

func ToUpper(input string) string {
	return strings.ToUpper(input)
}

func Count(input string) int {
	return strings.Count(input, "") - 1
}

func GetWordLength(word string) int {
	var length int
	for _, letter := range word {
		if RuneToString(letter) != " " {
			length++
		}
	}
	return length
}

func Fields(input string) []string {
	return strings.Fields(input)
}

func ToSentences(input string) []string {
	var output []string
	var sentence string
	for _, l := range input {
		sentence = sentence + string(l)
		if l == '\n' {
			output = append(output, sentence)
			sentence = ""
		}
	}
	return output
}

func ArrayRemoveEmptyStrings(input []string) []string {
	var output []string
	for _, entry := range input {
		if entry != "" {
			output = append(output, entry)
		}
	}
	return output
}

func ArrayRemoveString(input []string, searchStr string) []string {
	var output []string
	for _, entry := range input {
		if entry != searchStr {
			output = append(output, entry)
		}
	}
	return output
}

func ArrayContains(input []string, searchStr string) bool {
	for _, entry := range input {
		if entry == searchStr {
			return true
		}
	}
	return false
}

func ArrayRemoveIndex(ar []string, i int) []string {
	return append(ar[:i], ar[i+1:]...)
}

func Split(input, sep string) []string {
	return strings.Split(input, sep)
}

func Join(input []string, sep string) string {
	return strings.Join(input, sep)
}

func RuneToString(input rune) string {
	output := TrimPrefix(strconv.QuoteRune(input), "'")
	output = TrimSuffix(output, "'")
	output = Replace(output, `\n`, "\n", -1)
	return output
}

func ToLower(input string) string {
	return strings.ToLower(input)
}

func TrimSpace(input string) string {
	return strings.TrimSpace(input)
}

func ToArray(input, sep string) []string {
	return strings.Split(input, sep)
}

func ArrayToString(ar []string, sep string) string {
	return strings.Join(ar, sep)
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func Contains(input, searchstr string) bool {
	return strings.Contains(input, searchstr)
}

func Title(input string) string {
	return strings.Title(input)
}

func ToString(input int) string {
	return strconv.Itoa(input)
}

func Int64ToString(input int64) string {
	return strconv.Itoa(int(input))
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func ToInt(input string) int {
	i, err := strconv.Atoi(input)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to cast string to int32!", false, true, true)
	return i
}

func ToInt32(input string) int32 {
	return int32(ToInt(input))
}

func Trim(input, cutset string) string {
	return strings.Trim(input, cutset)
}

func TrimPrefix(input, prefix string) string {
	return strings.TrimPrefix(input, prefix)
}

func HasPrefix(input, prefix string) bool {
	return strings.HasPrefix(input, prefix)
}

func HasSuffix(input, suffix string) bool {
	return strings.HasSuffix(input, suffix)
}

func TrimSuffix(input, suffix string) string {
	return strings.TrimSuffix(input, suffix)
}

func Replace(input, old, new string, n int) string {
	return strings.Replace(input, old, new, n)
}
