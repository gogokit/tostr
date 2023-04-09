package tostr

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/runtime/protoiface"
)

func NotExportFieldFilter(obj reflect.Value, fieldIdx int) (hitFilter bool) {
	field := obj.Type().Field(fieldIdx)
	fieldName := field.Name
	return fieldName[0] < 'A' || fieldName[0] > 'Z'
}

func ProtobufFieldFilter(obj reflect.Value, fieldIdx int) (hitFilter bool) {
	field := obj.Type().Field(fieldIdx)
	if strings.HasPrefix(field.Type.PkgPath(), "google.golang.org/protobuf") {
		return true
	}
	if !reflect.PtrTo(obj.Type()).Implements(reflect.TypeOf((*protoiface.MessageV1)(nil)).Elem()) {
		return false
	}
	_, inMap := map[string]struct{}{
		"state":                {},
		"sizeCache":            {},
		"unknownFields":        {},
		"XXX_NoUnkeyedLiteral": {},
		"XXX_unrecognized":     {},
		"XXX_sizecache":        {},
	}[field.Name]
	return inMap
}

func FilterByFieldName(names ...string) func(obj reflect.Value, fieldIdx int) (hitFilter bool) {
	needFilters := make(map[string]struct{}, len(names))
	for _, v := range names {
		needFilters[v] = struct{}{}
	}
	return func(obj reflect.Value, fieldIdx int) (hitFilter bool) {
		_, inMap := needFilters[obj.Type().Field(fieldIdx).Name]
		return inMap
	}
}
