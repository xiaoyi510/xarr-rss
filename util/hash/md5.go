package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type Md5 struct {
}

func (this Md5) HashString(content string) string {
	return this.Hash([]byte(content))
}
func (this Md5) Hash(content []byte) string {
	hash := md5.New()
	hash.Write(content)
	b := hash.Sum([]byte{})
	return hex.EncodeToString(b)
}
