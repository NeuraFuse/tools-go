package request

import (
	"bytes"

	"github.com/neurafuse/tools-go/api/client"
	"github.com/neurafuse/tools-go/objects/strings"
)

type F struct{}

func (f F) Router(context string) string {
	var response string
	var request string = "{\"context\": \"" + context + "\"}"
	var body *bytes.Reader = bytes.NewReader(strings.ToBytes(request))
	response = client.F.Router(client.F{}, "inference/gpt", "POST", "user/infrastructure", "", "", "", body)
	return response
}
