package metrics

import (
	"strconv"

	"github.com/hhatto/gocloc"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
)

func DevStats(scope string) {
	var paths []string
	var frame string
	if scope == "all" {
		frame = "Git"
		paths = []string{"github.com/neurafuse/tools-go/git"}
	} else if scope == vars.OrganizationNameRepo {
		frame = vars.OrganizationName
		paths = []string{"../neuracli@" + vars.NeuraCLIVersion, "../neurakube@" + vars.NeuraKubeVersion, "../tools-go@" + vars.ToolsGoVersion, "../lightning-py"}
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
	var codeLinesInt int = int(result.Total.Code)
	var codeLines string = strconv.Itoa(codeLinesInt)
	var codeLinesA4 string = strconv.Itoa(codeLinesInt / 72)
	var filesTotal string = strconv.Itoa(int(result.Total.Total))
	logging.Log([]string{"", vars.EmojiStatistics, vars.EmojiInfo}, "Lines: "+codeLines+" | A4: "+codeLinesA4+" pages | Files: "+filesTotal+"\n", 1)
	//Log([]string{"", vars.EmojiInfo, ""}, "Comments:      " + strconv.Itoa(int(result.Total.Comments)), 1)
	//Log([]string{"", vars.EmojiInfo, ""}, "Blank spaces:  " + strconv.Itoa(int(result.Total.Blanks)), 1)
}

func DevStatsExceptions(hide bool) {
	exceptions := []string{"../tools-go@" + vars.ToolsGoVersion + "/nlp/sentences/data/english.json", "../tmp/english.json"}
	if hide {
		filesystem.Move(exceptions[0], exceptions[1], false)
	} else {
		filesystem.Move(exceptions[1], exceptions[0], false)
	}
}
