package pkg

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5encode(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
