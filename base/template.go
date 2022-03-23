package base

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	json "github.com/json-iterator/go"
)

// RenderText parses t as a template body, then applies it to data and returns the result.
func RenderText(t string, data interface{}) string {
	funcMap := template.FuncMap{
		"trim": func(s string) string { return strings.TrimSpace(s) },
		"json": func(s interface{}) string {
			sb, err := json.Marshal(s)
			if err != nil {
				return ""
			}
			// return fmt.Sprintf("%s", string(sb))
			return strings.ReplaceAll(string(sb), "\"", "\\\"")
		},
	}
	tp, err := template.New("main").Funcs(funcMap).Parse(t)
	if err != nil {
		fmt.Println("RenderText parse err: ", err)
		return ""
	}
	buffer := new(bytes.Buffer)
	if err := tp.Execute(buffer, data); err != nil {
		fmt.Println("RenderText execute err: ", err)
		return ""
	}
	return string(buffer.Bytes())
}
