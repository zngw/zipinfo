package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Str(txt string) string {
	m := md5.New()
	m.Write([]byte(txt))
	return hex.EncodeToString(m.Sum(nil))
}
