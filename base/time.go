package base

import (
	"encoding/binary"
	"strconv"
	"strings"
	"time"
)

// MysqlTimeFormat format time to mysql
func MysqlTimeFormat(t time.Time) string {
	e := time.Time{}
	if t == e {
		t = time.Now().UTC()
	}
	return t.Format("2006-01-02 15:04:05.999999")
}

// GetMysqlNowUTC get utc now for mysql
func GetMysqlNowUTC() string {
	return MysqlTimeFormat(time.Time{})
}

// TimeDecode unmarshals a time from a slice.
func TimeDecode(b []byte) time.Time {
	i := int64(binary.BigEndian.Uint64(b))
	return time.Unix(i, 0)
}

// TimeEncode marshals a time into a slice.
func TimeEncode(t time.Time) []byte {
	buf := make([]byte, 8)
	u := uint64(t.Unix())
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

// CronTriggerAble check cronOn is triggerable
// must be following format:
//             <ul>
//                 <li>expression represents a set of times, using 3 space-separated fields: Hour Day Month Year</li>
//                 <li>Hour must be is 0-23, can * / , - be used as special characters</li>
//                 <li>Day must be is 1-31, can * / , - be used as special characters</li>
//                 <li>Month must be is 1-12, can * / , - be used as special characters</li>
//                 <li>Year must be 4 digits, can * / , - be used as special characters</li>
//                 <li>Asterisk ( * ): The asterisk indicates that the expression will match for all values of the field;</li>
//                 <li>Slash ( / ): Slashes are used to describe increments of ranges</li>
//                 <li>Hyphen ( - ): Hyphens are used to define ranges. For example, 9-17 would indicate every hour between 9am and 5pm inclusive.</li>
//                 <li>Comma ( , ): Commas are used to separate items of a list. </li>
//                 </ul>
//      For example:
//      <ul>
//                 <li>Push on 2018-1-1 0:00:00: * * * 2018</li>
//                 <li>Push on 2018-12-1 0:00:00: * * 12 2018</li>
//                 <li>Push on 2018-12-11 0:00:00: * 11 12 2018</li>
//                 <li>Push on 2018-4-19 12:00:00: 12 19 4 2018</li>
//                 <li>Push on 4-19 12:00:00 of every year: 12 19 4 *</li>
//                 <li>Push on 4-19 12:00:00 of every 4 year: 12 19 4 */4</li>
//                 <li>Push on 4-19 12:00:00 of every year between 2018-2020: 12 19 4 2018-2020</li>
//                 <li>Push on 4-19 12:00:00 of every year can be divisible by 1010 between 2018-2020: 12 19 4 2018-2020/1010</li>
//                 <li>Push on 4-19 12:00:00 of year 2018,2020: 12 19 4 2018,2020</li>
//                 </ul>
func CronTriggerAble(cronOn string) bool {
	cronOnArr := strings.Split(cronOn, " ")
	if len(cronOnArr) != 4 {
		return false
	}
	now := time.Now().UTC()
	hour := now.Hour()
	day := now.Day()
	month := int(now.Month())
	year := now.Year()
	// var err error
	// fmt.Println(year, " ", month, " ", day, " ", hour)
	// fmt.Println(cronOnArr)
	for k, v := range cronOnArr {
		v = strings.TrimSpace(v)
		// fmt.Println("*k: ", k, " v: ", v)
		if v == "*" {
			continue
		}
		// fmt.Println("intk: ", k, " v: ", v)
		// is a int value
		if vi, err := strconv.ParseInt(v, 10, 64); err == nil {
			// fmt.Println("v: ", v, " vi: ", vi, " year: ", year)
			if k == 0 && int(vi) == hour {
				continue
			} else if k == 1 && int(vi) == day {
				continue
			} else if k == 2 && int(vi) == month {
				continue
			} else if k == 3 && int(vi) == year {
				continue
			} else {
				return false
			}
		}
		// fmt.Println("-k: ", k, " v: ", v)
		// meaning v has * , - / included
		if strings.Count(v, "-") == 1 && strings.Count(v, ",") == 0 && strings.Count(v, "/") == 0 && strings.Count(v, "*") == 0 {
			// only contain -
			vs := strings.Split(v, "-")
			if len(vs) == 2 {
				vs0, vs0Err := strconv.ParseInt(strings.TrimSpace(vs[0]), 10, 64)
				vs1, vs1Err := strconv.ParseInt(strings.TrimSpace(vs[1]), 10, 64)
				if vs0Err == nil && vs1Err == nil {
					if k == 0 && int(vs0) <= hour && int(vs1) >= hour {
						continue
					} else if k == 1 && int(vs0) <= day && int(vs1) >= day {
						continue
					} else if k == 2 && int(vs0) <= month && int(vs1) >= month {
						continue
					} else if k == 3 && int(vs0) <= year && int(vs1) >= year {
						continue
					}
				}
			}
		}
		// fmt.Println(",k: ", k, " v: ", v)
		if strings.Count(v, "-") == 0 && strings.Count(v, ",") > 0 && strings.Count(v, "/") == 0 && strings.Count(v, "*") == 0 {
			// only contain ,
			vs := strings.Split(v, ",")
			var vvs []int64
			for _, vv := range vs {
				vvsi, err := strconv.ParseInt(strings.TrimSpace(vv), 10, 64)
				if err != nil {
					return false
				}
				vvs = append(vvs, vvsi)
			}
			if k == 0 && FindInInt64Slice(vvs, int64(hour)) {
				continue
			} else if k == 1 && FindInInt64Slice(vvs, int64(day)) {
				continue
			} else if k == 2 && FindInInt64Slice(vvs, int64(month)) {
				continue
			} else if k == 3 && FindInInt64Slice(vvs, int64(year)) {
				continue
			}
		}
		// fmt.Println(",-k: ", k, " v: ", v)
		if strings.Count(v, "-") > 0 && strings.Count(v, ",") > 0 && strings.Count(v, "/") == 0 && strings.Count(v, "*") == 0 {
			// only contain - ,
			vs := strings.Split(v, ",")
			pass := false
			for _, l := range vs {
				l = strings.TrimSpace(l)
				if strings.Count(l, "-") > 0 {
					lsv := strings.Split(l, "-")
					if len(lsv) == 2 {
						vs0, vs0Err := strconv.ParseInt(strings.TrimSpace(lsv[0]), 10, 64)
						vs1, vs1Err := strconv.ParseInt(strings.TrimSpace(lsv[1]), 10, 64)
						if vs0Err == nil && vs1Err == nil {
							if k == 0 && int(vs0) <= hour && int(vs1) >= hour {
								pass = true
							} else if k == 1 && int(vs0) <= day && int(vs1) >= day {
								pass = true
							} else if k == 2 && int(vs0) <= month && int(vs1) >= month {
								pass = true
							} else if k == 3 && int(vs0) <= year && int(vs1) >= year {
								pass = true
							}
						}
					}
				} else {
					if lss, err := strconv.ParseInt(l, 10, 64); err == nil {
						if k == 0 && int(lss) == hour {
							pass = true
						} else if k == 1 && int(lss) == day {
							pass = true
						} else if k == 2 && int(lss) == month {
							pass = true
						} else if k == 3 && int(lss) == year {
							pass = true
						}
					}
				}
			}
			if pass {
				continue
			}
		}
		// fmt.Println("*/k: ", k, " v: ", v)
		if strings.Count(v, "-") == 0 && strings.Count(v, ",") == 0 && strings.Count(v, "/") == 1 && strings.Count(v, "*") == 1 {
			// only contain one / *, */20
			vs := strings.Split(v, "/")
			if strings.TrimSpace(vs[0]) != "*" || len(vs) != 2 {
				return false
			}
			vs2, err := strconv.ParseInt(strings.TrimSpace(vs[1]), 10, 64)
			if err != nil {
				return false
			}
			if k == 0 && (hour%int(vs2) == 0) {
				continue
			} else if k == 1 && (day%int(vs2) == 0) {
				continue
			} else if k == 2 && (month%int(vs2) == 0) {
				continue
			} else if k == 3 && (year%int(vs2) == 0) {
				continue
			}
		}
		return false
	}
	return true
}

// CronAbleFormat check cronOn format
func CronAbleFormat(cronOn string) bool {
	cronOnArr := strings.Split(cronOn, " ")
	if len(cronOnArr) != 4 {
		return false
	}
	for k, v := range cronOnArr {
		v = strings.TrimSpace(v)
		// fmt.Println("k: ", k, " v: ", v)
		if v == "*" {
			continue
		}
		// is a int value
		if vi, err := strconv.ParseInt(v, 10, 64); err == nil {
			if k == 0 && int(vi) >= 0 && int(vi) <= 23 {
				continue
			} else if k == 1 && int(vi) >= 1 && int(vi) <= 31 {
				continue
			} else if k == 2 && int(vi) >= 1 && int(vi) <= 12 {
				continue
			} else if k == 3 && int(vi) >= time.Now().UTC().Year() {
				continue
			}
		}
		// meaning v has * , - / included
		if strings.Count(v, "-") == 1 && strings.Count(v, ",") == 0 && strings.Count(v, "/") == 0 && strings.Count(v, "*") == 0 {
			// only contain -
			vs := strings.Split(v, "-")
			if len(vs) == 2 {
				vs0, vs0Err := strconv.ParseInt(strings.TrimSpace(vs[0]), 10, 64)
				vs1, vs1Err := strconv.ParseInt(strings.TrimSpace(vs[1]), 10, 64)
				if vs0Err == nil && vs1Err == nil {
					if k == 0 && int(vs0) <= 23 && int(vs1) >= 0 && int(vs1) <= 23 && int(vs0) >= 0 {
						continue
					} else if k == 1 && int(vs0) <= 31 && int(vs1) >= 1 && int(vs1) <= 31 && int(vs0) >= 1 {
						continue
					} else if k == 2 && int(vs0) <= 12 && int(vs1) >= 1 && int(vs1) <= 12 && int(vs0) >= 1 {
						continue
					} else if k == 3 && int(vs1) >= time.Now().UTC().Year() {
						continue
					}
				}
			}
		}
		if strings.Count(v, "-") == 0 && strings.Count(v, ",") > 0 && strings.Count(v, "/") == 0 && strings.Count(v, "*") == 0 {
			// only contain ,
			vs := strings.Split(v, ",")
			pass := true
			for _, vv := range vs {
				vvsi, err := strconv.ParseInt(strings.TrimSpace(vv), 10, 64)
				if err != nil {
					pass = false
				}
				if k == 0 && (vvsi > 23 || vvsi < 0) {
					pass = false
				} else if k == 1 && (vvsi > 31 || vvsi < 1) {
					pass = false
				} else if k == 2 && (vvsi > 12 || vvsi < 1) {
					pass = false
				}
			}
			if pass {
				continue
			}
		}
		if strings.Count(v, "-") > 0 && strings.Count(v, ",") > 0 && strings.Count(v, "/") == 0 && strings.Count(v, "*") == 0 {
			// only contain - ,
			vs := strings.Split(v, ",")
			pass := false
			for _, l := range vs {
				l = strings.TrimSpace(l)
				if strings.Count(l, "-") > 0 {
					lsv := strings.Split(l, "-")
					if len(lsv) == 2 {
						vs0, vs0Err := strconv.ParseInt(strings.TrimSpace(lsv[0]), 10, 64)
						vs1, vs1Err := strconv.ParseInt(strings.TrimSpace(lsv[1]), 10, 64)
						if vs0Err == nil && vs1Err == nil {
							if k == 0 && int(vs0) <= 23 && int(vs1) >= 0 && int(vs1) <= 23 && int(vs0) >= 0 {
								pass = true
							} else if k == 1 && int(vs0) <= 31 && int(vs1) >= 1 && int(vs1) <= 31 && int(vs0) >= 1 {
								pass = true
							} else if k == 2 && int(vs0) <= 12 && int(vs1) >= 1 && int(vs1) <= 12 && int(vs0) >= 1 {
								pass = true
							} else if k == 3 && int(vs0) < int(vs1) {
								pass = true
							}
						}
					}
				} else {
					if lss, err := strconv.ParseInt(l, 10, 64); err == nil {
						if k == 0 && int(lss) <= 23 && int(lss) >= 0 {
							pass = true
						} else if k == 1 && int(lss) <= 31 && int(lss) >= 1 {
							pass = true
						} else if k == 2 && int(lss) <= 12 && int(lss) >= 1 {
							pass = true
						} else if k == 3 {
							pass = true
						}
					}
				}
			}
			if pass {
				continue
			}
		}
		if strings.Count(v, "-") == 0 && strings.Count(v, ",") == 0 && strings.Count(v, "/") == 1 && strings.Count(v, "*") == 1 {
			// only contain one / *, */20
			vs := strings.Split(v, "/")
			if strings.TrimSpace(vs[0]) != "*" || len(vs) != 2 {
				return false
			}
			_, err := strconv.ParseInt(strings.TrimSpace(vs[1]), 10, 64)
			if err != nil {
				return false
			}
			continue
		}
		return false
	}
	return true
}
