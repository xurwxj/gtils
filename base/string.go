package base

import (
	"math/rand"
	"reflect"
	"regexp"
	"time"
	"unicode"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// SpecStringRandSeq get specific seq
func SpecStringRandSeq(n int, runeStr string) string {
	if runeStr == "" {
		runeStr = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	var letters = []rune(runeStr)
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Contain 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// Shuffle get Shuffle ints
func Shuffle(a []int) []int {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}

// Escape escape spec character
func Escape(source string) string {
	var j int = 0
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

// DecodingFromString get encoding from filepath
func DecodingFromString(f string) string {

	e, _, _ := charset.DetermineEncoding([]byte(f), "")
	if e != nil {
		rs, err := e.NewDecoder().String(f)
		if err == nil {
			return rs
		}
	}
	return ""
}

// DecodingKoreanString get encoding from filepath
func DecodingKoreanString(f string) string {

	es := korean.All
	for _, e := range es {
		rs, err := e.NewDecoder().String(f)
		if err == nil {
			return rs
		}
	}
	return ""
}

// DecodingGBKString get encoding from filepath
func DecodingGBKString(f string) string {

	es := simplifiedchinese.All
	for _, e := range es {
		rs, err := e.NewDecoder().String(f)
		if err == nil {
			return rs
		}
	}
	return ""
}

// DecodingBIG5String get encoding from filepath
func DecodingBIG5String(f string) string {

	e := traditionalchinese.Big5
	rs, err := e.NewDecoder().String(f)
	if err == nil {
		return rs
	}
	return ""
}

// DecodingJPString get encoding from filepath
func DecodingJPString(f string) string {

	es := japanese.All
	for _, e := range es {
		rs, err := e.NewDecoder().String(f)
		if err == nil {
			return rs
		}
	}
	return ""
}

// HasJP check string contain japanese
func HasJP(data string) bool {
	for _, v := range data {
		if unicode.Is(unicode.Hiragana, v) || unicode.Is(unicode.Katakana, v) {
			return true
		}
	}
	return false
}

// HasGBK check string contain gbk
func HasGBK(data string) bool {
	for _, v := range data {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// HasGBKReg check string contain gbk
func HasGBKReg(data string) bool {
	var reg = regexp.MustCompile("[\u4e00-\u9fa5]$")
	for _, v := range data {
		if reg.MatchString(string(v)) {
			return true
		}
	}
	return false
}

// HasJPReg check string contain japanese
func HasJPReg(data string) bool {
	var reg = regexp.MustCompile("[\u3000-\u303F\u3040-\u309F\u30A0-\u30FF\uFF00-\uFFEF\u2605-\u2606\u2190-\u2195\u203B]|[一-龠]|[ぁ-ん]|[ァ-ヴー]")
	for _, v := range data {
		if reg.MatchString(string(v)) {
			return true
		}
	}
	return false
}

// KeepRegexCharacter keep all character matched in regex
func KeepRegexCharacter(str, regexStr string) (rs string, err error) {
	if str == "" {
		return
	}
	reg, err := regexp.Compile(regexStr)
	if err != nil {
		return
	}
	rs = reg.ReplaceAllString(str, "")
	return
}

// KeepRegexCharacter keep all character matched in regex
func SubStringByLen(str string, t int) (rs string) {
	if str == "" {
		return
	}
	rs = str[:t]
	return
}
