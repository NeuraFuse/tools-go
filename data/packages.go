package data

import (
	"github.com/neurafuse/tools-go/data/processors/nlp/knowledge"
	"github.com/neurafuse/tools-go/data/providers/commoncrawl"
)

type Packages struct {
	Commoncrawl commoncrawl.F
	Knowledge   knowledge.F
}
