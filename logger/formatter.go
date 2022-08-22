package logger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"seb7887/logsafe/masker"
)

func formatArgs(m masker.Masker, args ...interface{}) string {
	str := ""
	for _, v := range args {
		v, _ := m.MaskSensitiveData(v)
		value := reflect.ValueOf(v)
		if value.Kind() == reflect.Ptr {
			str += jsonToString(&v)
		} else {
			str += fmt.Sprintf("%s", v)
		}
	}
	return str
}

func jsonToString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
