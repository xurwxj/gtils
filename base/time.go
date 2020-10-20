package base

import "time"

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
