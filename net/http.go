package net

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
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

// CheckIPPublic check ip is not a private ip or multicast ip
func CheckIPPublic(ip string) (bool, error) {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
		return false, err
	}
	targetIP := net.ParseIP(ip)
	if targetIP != nil && !targetIP.IsLinkLocalMulticast() && !targetIP.IsGlobalUnicast() {
		for _, block := range privateIPBlocks {
			if block.Contains(targetIP) {
				return false, fmt.Errorf("%s", block.Network())
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown")
}

// PaginationFromMap generate pagination from map
func PaginationFromMap(m map[string]interface{}, pSize, defaultSize int64) (map[string]interface{}, int, int) {
	var page = 1
	var pageSize = defaultSize
	if pSize > 0 {
		pageSize = pSize
	}
	qPSize, has := m["pageSize"]
	if has {
		qPSizeV, qPSerr := strconv.ParseInt(qPSize.(string), 10, 64)
		if qPSerr == nil {
			pageSize = qPSizeV
		}
	}
	if pageSize < 1 {
		pageSize = 10
	}
	qPage, pHas := m["page"]
	if pHas {
		qPageV, qPErr := strconv.Atoi(qPage.(string))
		if qPErr == nil {
			page = qPageV
		}
	}
	if page < 1 {
		page = 1
	}
	m["page"] = page
	m["pageSize"] = pageSize
	m["offset"] = (page - 1) * int(pageSize)
	// return map[string]interface{}{"page": page, "pageSize": pageSize, "offset": (page - 1) * int(pageSize)}, int(page), int(pageSize)
	return m, int(page), int(pageSize)
}

// QueryBytesToMap query bytes convert to map
// s := "A=B&C=D&E=F"
func QueryBytesToMap(query []byte) map[string]string {
	m := make(map[string]string)
	queries := strings.Split(string(query), "&")
	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q != "" {
			zs := strings.Split(q, "=")
			if len(zs) == 2 {
				k := strings.TrimSpace(zs[0])
				v := strings.TrimSpace(zs[1])
				if k != "" && v != "" {
					m[k] = v
				}
			}
		}
	}
	return m
}

// PaginationInfo generate pagination result info object
func PaginationInfo(page, pageSize, defaultSize, total int) map[string]int {
	if pageSize == 0 {
		pageSize = defaultSize
	}
	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}
	totalPage := total / pageSize
	if total%pageSize > 0 {
		totalPage++
	}
	return map[string]int{"page": page, "pageSize": pageSize, "pages": totalPage, "total": total}
}
