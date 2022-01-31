package nlp

import (
	"strings"
)

func VerbToAction(input string) string {
	if strings.HasSuffix(input, "t") {
		input = input+"t"
	} else if strings.HasSuffix(input, "s") {
		input = strings.Trim(input, "s")
		input = input + "g"
	}
	return input
}

func ConvertToPlural(resourceType string) (string, bool) {
	var actionTypeMultiple bool
	if resourceType[len(resourceType)-1:] == "s" {
		actionTypeMultiple = true
	} else {
		resourceType = resourceType + "s"
	}
	return resourceType, actionTypeMultiple
}