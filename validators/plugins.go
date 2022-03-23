package validators

import (
	"reflect"
	"strings"

	json "github.com/json-iterator/go"

	"github.com/fatih/structs"
	validator "github.com/go-playground/validator/v10"
	"github.com/xurwxj/gtils/base"
)

// ValidStruct return json string error
func ValidStruct(Validate *validator.Validate, tu interface{}) string {
	targets := []interface{}{}
	s := reflect.ValueOf(tu)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		for i := 0; i < s.Len(); i++ {
			targets = append(targets, s.Index(i).Interface())
		}
	} else {
		targets = append(targets, tu)
	}
	// fmt.Println(targets)

	rs := make(map[string]string)
	for _, u := range targets {
		err := Validate.Struct(u)
		if err != nil {
			su := structs.New(u)
			for _, verr := range err.(validator.ValidationErrors) {
				// tParam := verr.Param()
				// fmt.Println("verr.Field(): ", verr.Field())
				ff, ok := su.FieldOk(verr.Field())
				if ok {
					// fmt.Println("su.Field(verr.Field()): ", ff)
					// fmt.Println("su.Field(verr.Field()).Tag(\"json\"): ", ff.Tag("json"))
					rs[ff.Tag("json")] = verr.Tag()
				} else {
					// fmt.Println("verr StructField: ", verr.StructField())
					// fmt.Println("verr.Tag(): ", verr.Tag())
					rs[verr.Field()] = verr.Tag()
				}
			}
		}
	}
	if len(rs) > 0 {
		rsByte, err := json.Marshal(rs)
		if err == nil {
			return string(rsByte)
		} else {
			return "unknownErr"
		}
	}
	return ""
}

//IsCronOn check field is a formal cronable string
func IsCronOn(fl validator.FieldLevel) bool {
	return base.CronAbleFormat(fl.Field().String())
}

// CantBothEmpty depecrated
func CantBothEmpty(fl validator.FieldLevel) bool {
	tf, tk, hasTargetField := fl.GetStructFieldOK()
	if hasTargetField && tk == reflect.String && strings.TrimSpace(tf.String()) != "" {
		return true
	}
	if hasTargetField && tk == reflect.Slice && tf.Len() > 0 {
		return true
	}
	if hasTargetField && tk == reflect.Int && tf.Int() > 0 {
		return true
	}
	if hasTargetField && tk == reflect.Int64 && tf.Int() > 0 {
		return true
	}
	cf := fl.Field()
	ck := cf.Kind()
	if ck == reflect.String && strings.TrimSpace(cf.String()) != "" {
		return true
	}
	if ck == reflect.Slice && cf.Len() > 0 {
		return true
	}
	if ck == reflect.Int && cf.Int() > 0 {
		return true
	}
	if ck == reflect.Int64 && cf.Int() > 0 {
		return true
	}
	return false
}
