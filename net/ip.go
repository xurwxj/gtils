package net

import (
	"fmt"
	"net"

	json "github.com/json-iterator/go"

	"github.com/xurwxj/viper"
)

var privateIPs = []string{
	"127.0.0.0/8", // IPv4 loopback
	"100.64.0.0/10",
	"198.18.0.0/15",
	"255.255.255.255/32",
	"192.0.0.0/24",    //RFC5736
	"198.51.100.0/24", //RFC5737
	"192.0.2.0/24",    //RFC5737
	"203.0.113.0/24",  //RFC5737
	"10.0.0.0/8",      // RFC1918
	"172.16.0.0/12",   // RFC1918
	"192.168.0.0/16",  // RFC1918
	"169.254.0.0/16",  // RFC3927 link-local
	"::1/128",         // IPv6 loopback
	"fe80::/10",       // IPv6 link-local
	"fc00::/7",        // IPv6 unique local addr
}

// CheckIPPublic check ip is not a private ip or multicast ip
func CheckIPPublic(ip string) (bool, error) {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range privateIPs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		} else {
			return false, err
		}
	}
	targetIP := net.ParseIP(ip)
	if targetIP != nil && !targetIP.IsLinkLocalMulticast() {
		for _, block := range privateIPBlocks {
			if block.Contains(targetIP) {
				return false, fmt.Errorf("%s", block.Network())
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown")
}

func checkIPIn(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range privateIPs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		} else {
			return false
		}
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// GetAvailableIP get available ip for visiting out of this server, ipv can be v4 or v6
func GetAvailableIP(ipv string) string {
	ip := net.IP{}
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err == nil {
				for _, addr := range addrs {
					switch v := addr.(type) {
					case *net.IPNet:
						if v.IP.IsGlobalUnicast() {
							if ipv == "v4" && v.IP.To4() == nil {
								break
							}
							if ipv == "v6" && v.IP.To16() == nil {
								break
							}
							if !checkIPIn(v.IP) {
								ip = v.IP
							}
						}
					case *net.IPAddr:
						if v.IP.IsGlobalUnicast() {
							if ipv == "v4" && v.IP.To4() == nil {
								break
							}
							if ipv == "v6" && v.IP.To16() == nil {
								break
							}
							if !checkIPIn(v.IP) {
								ip = v.IP
							}
						}
					}
				}
			}
		}
	}
	return ip.String()
}

// GetVenuscnIPRS get ip info from venus api
func GetVenuscnIPRS(ip string) (rs AliIP) {
	venuscnURL := viper.GetString("ip.venuscn.url")
	venuscnCode := viper.GetString("ip.venuscn.code")
	if venuscnURL == "" || venuscnCode == "" {
		return
	}
	url := fmt.Sprintf("%s%s", venuscnURL, ip)
	resByte, _, err := Remote(url, "GET", nil, map[string]string{
		"accept":        "application/json",
		"content-type":  "application/json",
		"Authorization": fmt.Sprintf("APPCODE %s", venuscnCode),
	})
	if err != nil {
		return
	}
	var res ipDetail
	err = json.Unmarshal(resByte, &res)
	if err == nil && res.Code == 0 && res.Data.CityID != "lcoal" && res.Data.CountryID != "IANA" {
		rs = res.Data
	}
	if err == nil && res.Code == 0 && (res.Data.CityID == "lcoal" || res.Data.CountryID == "IANA") {
		res.Data.CountryID = "local"
		res.Data.Country = "local"
		rs = res.Data
	}
	return
}

type ipDetail struct {
	Code int   `json:"code"`
	Data AliIP `json:"data"`
}

// AliIP ali ip parse
type AliIP struct {
	IP        string `json:"ip"`
	Country   string `json:"country"`
	Area      string `json:"area"`
	Region    string `json:"region"`
	City      string `json:"city"`
	County    string `json:"county"`
	ISP       string `json:"isp"`
	CountryID string `json:"country_id"`
	AreaID    string `json:"area_id"`
	RegionID  string `json:"region_id"`
	CityID    string `json:"city_id"`
	CountyID  string `json:"county_id"`
	ISPID     string `json:"isp_id"`
}
