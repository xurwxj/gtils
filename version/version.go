package version

import (
	"fmt"
	"strconv"
	"strings"
)

// version compare, judge curVer is smaller than targetVer
func CompareVersion(targetVer, curVer, upgradeMode string) bool {
	if curVer == targetVer {
		return false
	}
	curVer, targetVer = MatchLenVersion(curVer, targetVer)
	// fmt.Println("targetVer: ", targetVer)
	// fmt.Println("curVer: ", curVer)
	targetVers := strings.Split(targetVer, ".")
	curVers := strings.Split(curVer, ".")
	for i, v := range targetVers {
		if tvI, err := strconv.ParseInt(v, 10, 64); err == nil {
			if cvI, err := strconv.ParseInt(curVers[i], 10, 64); err == nil {
				// if the larget version is not permit to upgrade cross generation and target ver is a next generation version, ignore this version
				if (i == 0 && upgradeMode != "ignore" && cvI < tvI) || cvI > tvI {
					return false
				}
				if cvI < tvI {
					return true
				}
			}
		}
	}
	return false
}

func BetweenVersion(curVer, targetVer1, targetVer2 string) bool {
	return CompareVersionGE(curVer, targetVer1) && !CompareVersionGE(curVer, targetVer2)
}

// judge curVer is greater than or equal to targetVer
func CompareVersionGE(curVer, targetVer string) bool {
	if curVer == targetVer {
		return true
	}
	curVer, targetVer = MatchLenVersion(curVer, targetVer)
	// fmt.Println("curVer: ", curVer)
	// fmt.Println("targetVer: ", targetVer)
	targetVers := strings.Split(targetVer, ".")
	curVers := strings.Split(curVer, ".")
	for i, v := range curVers {
		cv, err := strconv.ParseInt(v, 10, 64)
		tv, terr := strconv.ParseInt(targetVers[i], 10, 64)
		if cv > tv && err == nil && terr == nil {
			return true
		} else if cv < tv && err == nil && terr == nil {
			return false
		}
	}
	return true
}

func MatchLenVersion(curVer, targetVer string) (string, string) {
	// fmt.Println("curVer: ", curVer)
	// fmt.Println("targetVer: ", targetVer)
	curSepLen := strings.Count(curVer, ".")
	tarSepLen := strings.Count(targetVer, ".")
	if curSepLen > tarSepLen {
		for i := 0; i < curSepLen-tarSepLen; i++ {
			if strings.LastIndex(targetVer, ".") == len(targetVer)-1 {
				targetVer = fmt.Sprintf("%s0", targetVer)
			} else {
				targetVer = fmt.Sprintf("%s.0", targetVer)
			}
		}
	} else if curSepLen < tarSepLen {
		for i := 0; i < tarSepLen-curSepLen; i++ {
			if strings.LastIndex(curVer, ".") == len(curVer)-1 {
				curVer = fmt.Sprintf("%s0", curVer)
			} else {
				curVer = fmt.Sprintf("%s.0", curVer)
			}
		}
	}
	return curVer, targetVer
}

func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}
