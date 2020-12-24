package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

// GetByteArrayFromInterface get byte array from interface
func GetByteArrayFromInterface(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

// GetStringFromInterface get string from interface
func GetStringFromInterface(i interface{}) string {
	if i == nil {
		return ""
	}
	// fmt.Println(i)
	rs, err := GetByteArrayFromInterface(i)
	if err != nil {
		return ""
	}
	return string(rs)
}

// StructToJSON convert struct to json by struct's json tag
func StructToJSON(item interface{}) (interface{}, error) {
	var rs interface{}
	dataByte, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(dataByte, &rs)
	}
	return rs, err
}

// StructSliceToJSONTagMap convert struct slice to map by struct's json tag
func StructSliceToJSONTagMap(items interface{}) ([]interface{}, error) {
	rs := make([]interface{}, 0)
	dataByte, err := json.Marshal(items)
	if err == nil {
		err = json.Unmarshal(dataByte, &rs)
	}
	return rs, err
}

// StructToJSONTagMap convert struct to map by struct's json tag
func StructToJSONTagMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		if tag != "" && tag != "-" {
			field := reflectValue.Field(i).Interface()
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToJSONTagMap(field)
			} else {
				res[tag] = field
			}
		} else {
			if v.Field(i).Type.Kind() == reflect.Struct {
				for t, tv := range StructToJSONTagMap(reflectValue.Field(i).Interface()) {
					res[t] = tv
				}
			}
		}
	}
	return res
}

// StructToFormMap decodes an object into a map,
// which key is the string stored under the "form" key in the struct field's tag
// and value from the struct field.
func StructToFormMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Type.Kind() == reflect.String {
			data[obj1.Field(i).Tag.Get("form")] = obj2.Field(i).Interface()
		}
		if obj1.Field(i).Type.Kind() == reflect.Int || obj1.Field(i).Type.Kind() == reflect.Int64 {
			data[obj1.Field(i).Tag.Get("form")] = fmt.Sprintf("%d", obj2.Field(i).Interface())
		}
	}
	return data
}

// String2Bytes convert string to []byte
func String2Bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2String convert []byte to string
func Bytes2String(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}
