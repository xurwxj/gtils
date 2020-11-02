package net

import (
	"fmt"
	"net"
)

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
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
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
