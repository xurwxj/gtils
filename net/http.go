package net

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/xurwxj/gtils/validators"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

// RemoteReq remote req object
type RemoteReq struct {
	// 要请求的网址
	URL string `json:"url" validate:"required"`
	// 请求方法 大写的POST/GET...
	Method  string `json:"method" validate:"required"`
	Body    io.Reader
	Headers map[string]string
}

// RemoteHTTPClient remote http client use gentleman
func RemoteHTTPClient(reqData RemoteReq, validor *validator.Validate) ([]byte, string, error) {
	validRS := validators.ValidStruct(validor, reqData)
	if validRS != "" {
		return nil, validRS, fmt.Errorf(validRS)
	}
	cli := gentleman.New()
	cli.URL(reqData.URL)
	req := cli.Request()
	req.Method(reqData.Method)
	req.Use(body.Reader(reqData.Body))
	for k, v := range reqData.Headers {
		req.SetHeader(k, v)
	}
	res, err := req.Send()
	if err != nil {
		return nil, "", err
	}
	return res.Bytes(), "", nil
}

// Remote issues a POST to url with header and body resolved from formBody and env, returns the response body.
func Remote(url, method string, body io.Reader, header map[string]string) ([]byte, string, error) {
	var req *http.Request
	var res *http.Response
	var err error
	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	req, err = http.NewRequest(method, url, body)

	if err != nil {
		return nil, "netErr", err
	}
	req = setHeader(req, header)
	res, err = client.Do(req)
	if err != nil {
		return nil, "netErr", err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "parseErr", err
	}
	return data, "", nil
}

func setHeader(req *http.Request, header map[string]string) *http.Request {
	// fmt.Println("setHeader: ", header)
	// req.Header.Set("X-Auth-AppID", ReleaseAppID)
	if len(header) < 1 {
		return req
	}
	for k, h := range header {
		// fmt.Println("setHeader k: ", k)
		// fmt.Println("setHeader h: ", h)
		req.Header.Set(k, h)
	}
	return req
}

// PrintLocalDial connects to addr on the network.
// It returns the net.Conn and an error.
func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   15 * time.Minute,
		KeepAlive: 15 * time.Minute,
	}

	return dial.Dial(network, addr)
}

// ParseHeader parse header []byte to http.Header
func ParseHeader(headerByte []byte) (http.Header, error) {
	var h http.Header
	reader := bufio.NewReader(bytes.NewReader(headerByte))

	reqHeader, err := http.ReadRequest(reader)
	if err == nil {
		h = reqHeader.Header
	}

	return h, err

}

// GetCookieFromHeaderByte get cookie value from headerByte by key name
func GetCookieFromHeaderByte(headerByte []byte, key string) (string, error) {
	h, err := ParseHeader(headerByte)
	if err == nil {
		cookieStrs := h.Get("Cookie")
		cookieStrs = strings.TrimSpace(cookieStrs)
		if cookieStrs != "" {
			for _, cookieObjStr := range strings.Split(cookieStrs, ";") {
				cookieObjStr = strings.TrimSpace(cookieObjStr)
				if cookieObjStr != "" {
					zs := strings.Split(cookieObjStr, "=")
					if len(zs) == 2 {
						k := strings.TrimSpace(zs[0])
						if k == key {
							return strings.TrimSpace(zs[1]), nil
						}
					}
				}
			}
		}
	}
	return "", err
}

// QueryParseToStruct parse query or form body to struct
func QueryParseToStruct(query []byte, out interface{}) error {

	queryMap := QueryBytesToMap(query)
	queryBytes, err := json.Marshal(queryMap)
	// fmt.Println(string(queryBytes))
	if err != nil {
		return err
	}
	return json.Unmarshal(queryBytes, out)
}
