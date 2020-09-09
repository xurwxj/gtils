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

func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	// fmt.Println("v: ", reflect.ValueOf(data))
	elem := reflect.ValueOf(data)
	size := elem.NumField()
	// fmt.Println("s: ", size)

	for i := 0; i < size; i++ {
		// fmt.Println("f: ", elem.Field(i))
		// fmt.Println("t: ", elem.Type().Field(i))
		field := elem.Type().Field(i).Name
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

func StructToJsonTagMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Tag.Get("json")
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

func StructToJsonTagMap2(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	return result
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
