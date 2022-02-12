package tostr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var defaultConfig = Config{
	ToString: func(o reflect.Value) (objStr string) {
		if f, inMap := toStringMap[o.Type()]; inMap {
			return f(o)
		}

		return fmt.Sprintf("<Fatal error：format func for %v not register>", o.Type())
	},
	FastSpecifyToStringProbe: func(o reflect.Value) (hasSpecifyToString bool) {
		_, inMap := toStringMap[o.Type()]
		return inMap
	},
	FilterStructField: []func(reflect.Value, int) bool{
		func(obj reflect.Value, fieldIdx int) (hitFilter bool) {
			return !obj.Type().Field(fieldIdx).IsExported()
		},
	},
	WarnSize: func(num int) *int {
		return &num
	}(1e4),
	ResultSizeWarnCallback: func(str string) (shouldContinue bool) {
		return false
	},
}

var toStringMap = map[reflect.Type]func(obj reflect.Value) string{
	reflect.TypeOf(time.Time{}): func(obj reflect.Value) string {
		return "{time.Time:\"" + obj.Interface().(time.Time).Format("2006-01-02 15:04:05.000") + "\"}"
	},
	reflect.TypeOf([]byte{}): func(obj reflect.Value) string {
		return strconv.Quote(string(obj.Interface().([]byte)))
	},
	reflect.TypeOf(json.RawMessage{}): func(obj reflect.Value) string {
		var s []byte = obj.Interface().(json.RawMessage)
		return strconv.Quote(string(s))
	},
}
