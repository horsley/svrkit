package svrkit

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

//MD5Hash md5
func MD5Hash(src string) string {
	hash := md5.New()
	hash.Write([]byte(src))
	return hex.EncodeToString(hash.Sum(nil))
}

//SHA1Hash sha1
func SHA1Hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
