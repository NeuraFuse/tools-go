package metrics

import (
	"strconv"

	"../errors"
	"../filesystem"
	"../logging"
	"../runtime"
	"../vars"
	"github.com/hhatto/gocloc"
)

func DevStats(scope string) {
	var paths []string
	var frame string = ""
	if scope == "all" {
		frame = "Git"
		paths = []string{"../../git"}
	} else if scope == vars.OrganizationNameRepo {
		frame = vars.OrganizationName
		paths = []string{"../neuracli", "../neurakube", "../tools-go"}
	}
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()
	processor := gocloc.NewProcessor(languages, options)
	DevStatsExceptions(true)
	result, err := processor.Analyze(paths)
	DevStatsExceptions(false)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Code analysation failed!", false, false, true)
	/*for _, file := range result.Files {
		fmt.Println(file)
	}
	for _, lang := range result.Languages {
		fmt.Println(lang)
	}*/
	logging.Log([]string{"", vars.EmojiDev, vars.EmojiStatistics}, "Dev statistics ("+frame+")", 1)
	codeLinesInt := int(result.Total.Code)
	codeLines := strconv.Itoa(codeLinesInt)
	codeLinesA4 := strconv.Itoa(codeLinesInt / 72)
	filesTotal := strconv.Itoa(int(result.Total.Total))
	logging.Log([]string{"", vars.EmojiStatistics, vars.EmojiInfo}, "Lines: "+codeLines+" | A4: "+codeLinesA4+" pages | Files: "+filesTotal+"\n", 1)
	//Log([]string{"", vars.EmojiInfo, ""}, "Comments:      " + strconv.Itoa(int(result.Total.Comments)), 1)
	//Log([]string{"", vars.EmojiInfo, ""}, "Blank spaces:  " + strconv.Itoa(int(result.Total.Blanks)), 1)
}

func DevStatsExceptions(hide bool) {
	exceptions := []string{"../tools-go/nlp/sentences/data/english.json", "../tmp/english.json"}
	if hide {
		filesystem.Move(exceptions[0], exceptions[1], false)
	} else {
		filesystem.Move(exceptions[1], exceptions[0], false)
	}
}
