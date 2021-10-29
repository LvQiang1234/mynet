package base

import (
	"crypto/md5"
	"fmt"
	"io"
)

// md5加密字符串
func MD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
