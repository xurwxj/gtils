package base

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"unsafe"

	json "github.com/json-iterator/go"

	"github.com/spf13/cast"
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

// StructToJSONString convert struct to json string by struct's json tag
func StructToJSONString(item interface{}) (string, error) {
	dataByte, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return string(dataByte), nil
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
func StructToJSONTagMap(item interface{}) (res map[string]interface{}) {
	if item == nil {
		return
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if res == nil {
		res = make(map[string]interface{})
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		if tag != "" && tag != "-" {
			field := reflectValue.Field(i).Interface()
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToJSONTagMap(field)
			} else if v.Field(i).Type.Kind() == reflect.Slice {
				sliceResult, err := StructSliceToJSONTagMap(field)
				if err == nil {
					res[tag] = sliceResult
				}
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
	return
}

// StructToJSONTagtringMap convert struct to string map by struct's json tag
func StructToJSONTagtringMap(item interface{}) (res map[string]string) {
	if item == nil {
		return
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if res == nil {
		res = make(map[string]string)
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		if tag != "" && tag != "-" {
			field := reflectValue.Field(i).String()
			if v.Field(i).Type.Kind() == reflect.String {
				res[tag] = field
			} else if v.Field(i).Type.Kind() == reflect.Struct {
				for t, tv := range StructToJSONTagtringMap(reflectValue.Field(i).Interface()) {
					res[t] = tv
				}
			} else {
				continue
			}
		} else {
			if v.Field(i).Type.Kind() == reflect.Struct {
				for t, tv := range StructToJSONTagtringMap(reflectValue.Field(i).Interface()) {
					res[t] = tv
				}
			}
		}
	}
	return
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

// StructToDBMap decodes an object into a map of sqlx,
// which key is the string stored under the "db" key in the struct field's tag
// and value from the struct field.
func StructToDBMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Type.Kind() == reflect.String {
			data[obj1.Field(i).Tag.Get("db")] = obj2.Field(i).Interface()
		}
		if obj1.Field(i).Type.Kind() == reflect.Int || obj1.Field(i).Type.Kind() == reflect.Int64 {
			data[obj1.Field(i).Tag.Get("db")] = fmt.Sprintf("%d", obj2.Field(i).Interface())
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

// Round 四舍五入，ROUND_HALF_UP 模式实现
// 返回将 val 根据指定精度 precision（十进制小数点后数字的数目）进行四舍五入的结果。precision 也可以是负数或零。
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// ByteSliceStringToByteSlice 把字节数组的字符串转换成标准的字节数组
func ByteSliceStringToByteSlice(s string) (rs []byte) {
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	g := strings.Split(s, " ")
	for _, gv := range g {
		rs = append(rs, byte(cast.ToInt64(gv)))
	}
	return
}
