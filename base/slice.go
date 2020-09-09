package base

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FindInStringSlice returns whether slice s contains string t.
func FindInStringSlice(s []string, t string) bool {
	for _, v := range s {
		if strings.TrimSpace(v) == strings.TrimSpace(t) {
			return true
		}
	}
	return false
}

//  FindInInt64Slice returns whether slice s contains int64 t
func FindInInt64Slice(s []int64, t int64) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

func FindInIntSlice(s []int, t int) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

func IntSeq(ns []int) string {
	if len(ns) == 0 {
		return ""
	}

	// Appr. 3 chars per num plus the comma.
	estimate := len(ns) * 4
	b := make([]byte, 0, estimate)
	// Or simply
	//   b := []byte{}
	for _, n := range ns {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}

func IntSliceUniqMap(s []int) []int {
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func StringSliceUniqMap(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func GetBytesFromInterface(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SetStructField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func InsertSlice(slice, insertion []interface{}, index int) []interface{} {
	result := make([]interface{}, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}
