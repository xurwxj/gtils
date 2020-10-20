package net

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Remote issues a POST to url with header and body resolved from formBody and env, returns the response body.
func Remote(url, method string, body io.Reader, header map[string]string) ([]byte, error) {
	var req *http.Request
	var res *http.Response
	var err error
	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	req, err = http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}
	req = setHeader(req, header)
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
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

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	// fmt.Println("PrintLocalDial connect done, use", conn.LocalAddr().String())

	return conn, err
}
