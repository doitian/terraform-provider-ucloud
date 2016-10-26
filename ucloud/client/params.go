package client

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func safeIsNil(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return val.IsNil()
	default:
		return false
	}
}

func stringify(v reflect.Value) (str string, ok bool) {
	ok = false

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ok = true
		if number := v.Int(); number != 0 {
			str = strconv.FormatInt(number, 10)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ok = true
		if number := v.Uint(); number != 0 {
			str = strconv.FormatUint(number, 10)
		}

	case reflect.String:
		ok = true
		str = v.String()
	}

	return
}

func BuildParams(req interface{}) (url.Values, error) {
	ret := url.Values{}
	err := AddParams(ret, req)
	return ret, err
}

func AddParams(params url.Values, req interface{}) error {
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	valType := val.Type()
	typeName := valType.Name()
	action := strings.TrimSuffix(typeName, "Request")
	if len(action) != len(typeName) {
		params.Set("Action", action)
	}

	for i := 0; i < val.NumField(); i++ {
		typeField := valType.Field(i)
		field := val.Field(i)
		if safeIsNil(field) {
			continue
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		fieldKind := field.Kind()
		fieldName := typeField.Tag.Get("ArgName")
		if fieldName == "" {
			fieldName = typeField.Name
		}

		if fieldValue, ok := stringify(field); ok {
			if fieldValue != "" {
				params.Set(fieldName, fieldValue)
			}
		} else if fieldKind == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elemValue, ok := stringify(elem); ok {
					params.Set(fieldName+"."+strconv.Itoa(j), elemValue)
				} else {
					return fmt.Errorf("Cannot convert %s to params in slice", elem.Kind())
				}
			}
		} else {
			return fmt.Errorf("Cannot convert %s to params", fieldKind)
		}
	}

	return nil
}
