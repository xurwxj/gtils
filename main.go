package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/gtils/net"
	"github.com/xurwxj/gtils/validators"
)

func main() {
	vs()
}

func vs() {
	v := validator.New()
	d := []factoryMemberForm{
		{FactoryID: "888", Username: "999", IsAdmin: false},
		{FactoryID: "rrr", Username: "ttt", IsAdmin: true},
	}
	fmt.Println(validators.ValidStruct(v, d))
	fmt.Println("s: ", validators.ValidStruct(v, d[0]))
}

type factoryMemberForm struct {
	FactoryID string `form:"factoryID" json:"factoryID" validate:"required"`
	Username  string `form:"username" json:"username" validate:"required"`
	IsAdmin   bool   `form:"isAdmin" json:"isAdmin"`
}

func td() {
	s := FetchPLEForm{}
	net.QueryParseToStruct([]byte("modelCode=einscan-s&serialEncryption=einscan-plus2LGMHJHGGJJFDCFGKLJCBGHNIKJFBKKGEEEBLMKDE"), &s)
	fmt.Println(s)
}

// FetchPLEForm for plefetch form
type FetchPLEForm struct {
	ModelCode        string `db:"model_code" form:"modelCode" query:"modelCode" json:"modelCode" validate:"required"`
	SerialEncryption string `db:"serial_encryption" form:"serialEncryption" query:"serialEncryption" json:"serialEncryption" validate:"required"`
}

func testUnzip() {
	// InitLog()
	// s := sys.GetOsInfo(Log)
	// fmt.Println("s: ", s)
	// fmt.Println("v0.3.4: ", base.IsValidSemver("v0.3.4"))
	// fmt.Println("0.3.4: ", base.IsValidSemver("0.3.4"))
	// fmt.Println("0.3.4.8: ", base.IsValidSemver("0.3.4.8"))
	// fmt.Println("v0.3.4.8: ", base.IsValidSemver("v0.3.4.8"))
	// fmt.Println("0.3.4.8.9: ", base.IsValidSemver("0.3.4.8.9"))
	// fmt.Println("0.3.4.8.v9: ", base.IsValidSemver("0.3.4.8.v9"))
	// header := `User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16_0) AppleWebKit/537.36 (KHTML, like Gecko) Postman/4.1.1 Chrome/47.0.2526.73 Electron/0.36.2 Safari/537.36\r\nHost: 127.0.0.1:7080\r\nContent-Length: 0\r\nConnection: keep-alive\r\nCache-Control: no-cache\r\nPostman-Token: d1673e11-8aea-bd9f-1a89-1109cf2bcb29\r\nAccept: */*\r\nAccept-Encoding: gzip, deflate\r\nAccept-Language: en-US\r\nCookie: s=ui8\r\n\r\n`
	// s, err := net.GetCookieFromHeaderByte([]byte{71, 69, 84, 32, 47, 100, 115, 32, 72, 84, 84, 80, 47, 49, 46, 49, 13, 10, 85, 115, 101, 114, 45, 65, 103, 101, 110, 116, 58, 32, 77, 111, 122, 105, 108, 108, 97, 47, 53, 46, 48, 32, 40, 77, 97, 99, 105, 110, 116, 111, 115, 104, 59, 32, 73, 110, 116, 101, 108, 32, 77, 97, 99, 32, 79, 83, 32, 88, 32, 49, 48, 95, 49, 54, 95, 48, 41, 32, 65, 112, 112, 108, 101, 87, 101, 98, 75, 105, 116, 47, 53, 51, 55, 46, 51, 54, 32, 40, 75, 72, 84, 77, 76, 44, 32, 108, 105, 107, 101, 32, 71, 101, 99, 107, 111, 41, 32, 80, 111, 115, 116, 109, 97, 110, 47, 52, 46, 49, 46, 49, 32, 67, 104, 114, 111, 109, 101, 47, 52, 55, 46, 48, 46, 50, 53, 50, 54, 46, 55, 51, 32, 69, 108, 101, 99, 116, 114, 111, 110, 47, 48, 46, 51, 54, 46, 50, 32, 83, 97, 102, 97, 114, 105, 47, 53, 51, 55, 46, 51, 54, 13, 10, 72, 111, 115, 116, 58, 32, 49, 50, 55, 46, 48, 46, 48, 46, 49, 58, 55, 48, 56, 48, 13, 10, 67, 111, 110, 116, 101, 110, 116, 45, 76, 101, 110, 103, 116, 104, 58, 32, 48, 13, 10, 67, 111, 110, 110, 101, 99, 116, 105, 111, 110, 58, 32, 107, 101, 101, 112, 45, 97, 108, 105, 118, 101, 13, 10, 67, 97, 99, 104, 101, 45, 67, 111, 110, 116, 114, 111, 108, 58, 32, 110, 111, 45, 99, 97, 99, 104, 101, 13, 10, 80, 111, 115, 116, 109, 97, 110, 45, 84, 111, 107, 101, 110, 58, 32, 52, 51, 54, 48, 55, 54, 97, 52, 45, 101, 53, 48, 49, 45, 99, 100, 49, 49, 45, 48, 50, 52, 101, 45, 56, 57, 101, 102, 57, 102, 53, 52, 57, 52, 48, 97, 13, 10, 65, 99, 99, 101, 112, 116, 58, 32, 42, 47, 42, 13, 10, 65, 99, 99, 101, 112, 116, 45, 69, 110, 99, 111, 100, 105, 110, 103, 58, 32, 103, 122, 105, 112, 44, 32, 100, 101, 102, 108, 97, 116, 101, 13, 10, 65, 99, 99, 101, 112, 116, 45, 76, 97, 110, 103, 117, 97, 103, 101, 58, 32, 101, 110, 45, 85, 83, 13, 10, 67, 111, 111, 107, 105, 101, 58, 32, 115, 61, 117, 105, 56, 59, 32, 103, 61, 57, 57, 48, 13, 10, 13, 10}, "s")
	// fmt.Println("0.3.4.8.09: ", s)
	// fmt.Println("0.3.4.8.09: ", err)

	f := "../../../git.shining3d.com/cloud/algorithm/tmp/ŠpičákováVěra_2021-01-06_022.zip"
	base.Unzip(f, fmt.Sprintf("../../../git.shining3d.com/cloud/algorithm/tmp/mesh/%s", base.Md5String(f)), "")
	f = "../../../git.shining3d.com/cloud/algorithm/tmp/2021-01-08_003_111_谭彩红.zip"
	base.Unzip(f, fmt.Sprintf("../../../git.shining3d.com/cloud/algorithm/tmp/mesh/%s", base.Md5String(f)), "")
}
