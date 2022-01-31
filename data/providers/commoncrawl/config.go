package commoncrawl

import (
	//"os"
	"path"
)

// Config is the preset variables for your extractor
type Config struct {
	BaseURI     string
	WetPaths    string
	DataFolder  string
	MatchFolder string
	Start       int
	Stop        int
}

func GetConfiguration() Config {
	//cwd, err := os.Getwd()

	//if err != nil {
		//panic(err)
	//}
	cwd := "data/providers/commoncrawl/"
	return Config{
		Start:       0,
		Stop:        1,
		BaseURI:     "https://commoncrawl.s3.amazonaws.com/",
		WetPaths:    path.Join(cwd, "warc.paths"),
		DataFolder:  path.Join(cwd, "output/crawl-data"),
		MatchFolder: path.Join(cwd, "output/match-data"),
	}
}
