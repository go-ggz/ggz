package base

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// EncodeMD5 encodes string to md5 hex value.
func EncodeMD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// EncodeSha1 string to sha1 hex value.
func EncodeSha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// HashEmail hashes email address to MD5 string.
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	return EncodeMD5(strings.ToLower(strings.TrimSpace(email)))
}
