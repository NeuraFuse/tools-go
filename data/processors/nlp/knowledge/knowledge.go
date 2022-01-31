package knowledge

import (
	"github.com/neurafuse/tools-go/data/processors/nlp/tools/filters/web"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/nlp/filter"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/vars"
	//dict "../tools/enrichment/dictionary"
)

type F struct{}

func (f F) Live(index, contentNew, format string) string {
	index = web.F.Router(web.F{}, index, contentNew, format)
	index = f.basicFilters(index)
	f.metricsLive(index)
	return index
}

func (f F) basicFilters(index string) string {
	index = filter.Casing(index)
	index = filter.Digits(index)
	charFilter := f.getCharFilter()
	if len(charFilter) > 0 {
		index = filter.SpecialChars(index, charFilter)
	}
	//index = filter.WordsLength(index, "min", 3)
	index = filter.WordsLength(index, "max", 16)
	//contentNew = filter.SentenceMaxWords(contentNew, 20) // Filter too complex sentences
	//contentNew = filter.DuplicateLines(index, contentNew, lineSep) // TODO: Reactivate after debugging
	return index
}

var metricsC int

func (f F) metricsLive(index string) {
	metricsC++
	if metricsC == 100 {
		logMsg := "Filtered " + strings.ToString(len(index)) + " lines live."
		logging.Log([]string{"", vars.EmojiCompression, vars.EmojiInfo}, logMsg, 0)
		metricsC = 0
	}
}

func (f F) PostProcess(text string) (string, int) {
	text = web.F.Post(web.F{}, text)
	//input = dict.F.AddDefinitions(dict.F{}, input, []string{"N"})
	return text, strings.LinesCount(text)
}

func (f F) getCharFilter() []string {
	return []string{"]", "^", "\\\\", "[", "{", "}", "!", "?", "и", "к", "п", "↳", "ệ", "…", "›", "é", "„", "|", "$", "“", "”", "α", "β", "¯", "⋆", ",", "►", "∥", "↓", "→", "↑", "↖", "↗", "*", "+", "▹", "▾", "←", ":", ";", "’", ".", "_", "=", "↩", "~", "/", "(", ")"} // 92 backslash // []string{`\`, ".", "?", "!", "(", ")", "{", "}"})
}
