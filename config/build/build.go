package build

type F struct{}

var checkDo bool
var handover bool

func (f F) Setting(method string, setting string, value bool) bool {
	switch setting {
	case "check":
		if method == "set" {
			checkDo = value
		} else {
			return checkDo
		}
	case "handover":
		if method == "set" {
			handover = value
		} else {
			return handover
		}
	}
	return false
}
