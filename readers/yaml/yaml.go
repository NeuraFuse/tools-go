package yaml

import (
	"../../../tools-go/errors"
	"../../../tools-go/filesystem"
	"../../../tools-go/runtime"
	"github.com/ghodss/yaml"
	yamlv2 "gopkg.in/yaml.v2"
)

func FileToStruct(filePath string, i interface{}) interface{} {
	if filesystem.Exists(filePath) {
		//err := yamlv2.Unmarshal(filesystem.FileToBytes(filePath), &i)
		err := yaml.Unmarshal(filesystem.FileToBytes(filePath), &i)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	}
	return i
}

func StructToFile(filePath string, i interface{}) {
	//out, err := yamlv2.Marshal(i)
	out, err := yaml.Marshal(i)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	filesystem.SaveByteArrayToFile(out, filePath)
}

func ToBytes(i interface{}) []byte {
	out, err := yamlv2.Marshal(i)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	return out
}

func ToJSON(i interface{}) []byte {
	json, err := yaml.YAMLToJSON(ToBytes(i))
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	return json
}

func ConvertJSON(json []byte) []byte {
	y, err := yaml.JSONToYAML(json)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	return y
}
