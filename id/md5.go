package id

import (
	"crypto/md5"
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
)

// Md5String returns the MD5 checksum of v in string format.
func Md5String(v string) string {
	hasher := md5.New()
	hasher.Write([]byte(v))
	return hex.EncodeToString(hasher.Sum(nil))
}

func UUID(prefix string) string {
	return uuid.NewV5(uuid.NewV4(), prefix).String()
}
