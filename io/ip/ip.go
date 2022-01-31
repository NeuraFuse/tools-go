package ip

import (
	"github.com/neurafuse/tools-go/objects/strings"
)

func SplitBlocks(ipS string) (string, string) {
	var ip []string = strings.Split(ipS, ".")
	var network string = ip[0] + ip[1] + ip[2]
	var host string = ip[3]
	return network, host
}
