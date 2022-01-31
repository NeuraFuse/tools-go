package objects

import (
	"reflect"
	"strings"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
)

/*func GetInterfaceFromReflectValue(i type, any interface{}) interface{} {
	return reflect.ValueOf(any).Interface().(i)
}*/

func GetStringFromReflectValue(ref reflect.Value) string {
	return reflect.ValueOf(ref).String()
}

func CallStructInterfaceFuncByName(any interface{}, name string, funcName string, args ...interface{}) (bool, []reflect.Value) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	ref := reflect.ValueOf(any).FieldByName(name)
	method := ref.MethodByName(funcName)
	var refVal []reflect.Value
	if method.String() == "<invalid Value>" {
		return false, refVal
	} else {
		refVal = method.Call(inputs)
		return true, refVal
	}
}

func StructFieldValuesExisting(str interface{}, fieldsPrefix string, fields []string, packageName string) bool {
	var notEmpty bool = true
	for _, field := range fields {
		if GetFieldValueFromStruct(str, fieldsPrefix+field, packageName) == "" {
			notEmpty = false
		}
	}
	return notEmpty
}

func GetFieldValueFromStruct(str interface{}, key, packageName string) string {
	return GetReflectValueFromStruct(str, key, packageName).String()
}

func SetFieldValueFromStruct(str interface{}, key, value, packageName string) interface{} {
	str, ref := reflectFields(str, key, packageName)
	ref.SetString(value)
	return str
}

func GetReflectValueFromStruct(str interface{}, key, packageName string) reflect.Value {
	_, ref := reflectFields(str, key, packageName)
	return ref
}

func reflectFields(str interface{}, key, packageName string) (interface{}, reflect.Value) {
	ref := reflect.ValueOf(str).Elem()
	fields := strings.Split(key, ".")
	for _, field := range fields {
		if ref.FieldByName(field).IsValid() {
			ref = ref.FieldByName(field)
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get setting for key: "+key, true, false, true)
			if ref.Kind() != reflect.String {
				errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "The setting has an invalid type (not string) for key: "+key, true, true, true)
			}
		}
	}
	return str, ref
}

func callFunc(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "The number of params is not adapted.", true, true, true)
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
