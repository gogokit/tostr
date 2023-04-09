package tostr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

// 返回全局默认配置的深拷贝
func GetDefaultConfig() Config {
	return defaultConfig.Clone()
}

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
	FilterStructField: []func(reflect.Value, int) bool{NotExportFieldFilter, ProtobufFieldFilter},
	WarnSize: func(num int) *int {
		return &num
	}(1e4),
	ResultSizeWarnCallback: func(str string) (shouldContinue bool) {
		return false
	},
}

var toStringMap = map[reflect.Type]func(obj reflect.Value) string{
	reflect.TypeOf(time.Time{}): func(obj reflect.Value) string {
		if obj.CanInterface() {
			return "{time.Time:\"" + obj.Interface().(time.Time).Format("2006-01-02 15:04:05.000") + "\"}"
		}
		ptr := (*time.Time)((unsafe.Pointer)(obj.UnsafeAddr()))
		return "{time.Time:\"" + ptr.Format("2006-01-02 15:04:05.000") + "\"}"
	},
	reflect.TypeOf([]byte{}): func(obj reflect.Value) string {
		if obj.CanInterface() {
			return strconv.Quote(string(obj.Interface().([]byte)))
		}
		ptr := (*[]byte)((unsafe.Pointer)(obj.UnsafeAddr()))
		return strconv.Quote(string(*ptr))
	},
	reflect.TypeOf(json.RawMessage{}): func(obj reflect.Value) string {
		if obj.CanInterface() {
			return strconv.Quote(string(obj.Interface().(json.RawMessage)))
		}
		ptr := (*json.RawMessage)((unsafe.Pointer)(obj.UnsafeAddr()))
		return strconv.Quote(string(*ptr))
	},
}
