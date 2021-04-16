package json

import (
	"encoding/json"
	"net/http"

	"../../errors"
	"../../filesystem"
	"../../runtime"
	"../../timing"
)

var myClient = &http.Client{Timeout: timing.GetTimeDuration(5, "s")}

func URLToInterface(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func FileToStruct(filePath string, i interface{}) interface{} {
	if filesystem.Exists(filePath) {
		err := json.Unmarshal(filesystem.FileToBytes(filePath), &i)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	}
	return i
}

func StructToFile(filePath string, i interface{}) {
	out, err := json.MarshalIndent(i, "", "   ")
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	filesystem.SaveByteArrayToFile(out, filePath)
}