package envutil

import (
	"bytes"
	"os"
	"strings"
)

// ExpandEnv 环境变量占位符
// 格式： {NAME} (不带默认值)，{NAME:hello} (带默认值)
func ExpandEnv(s string) (string, error) {
	return expand(s, func(s string) (string, error) {
		return GetEnv(s), nil
	})
}

func Expand(s string, get func(string) (string, error)) (string, error) {
	return expand(s, get)
}

func expand(text string, f func(string) (string, error)) (string, error) {

	var buf []byte
	var key []byte
	var expanding bool
	var left = []byte{'$', '{'}
	var right = []byte{'}'}

	var i = 0
	for i < len(text) {
		if len(text) >= i+len(left) && bytes.Equal([]byte(text[i:i+len(left)]), left) {
			if expanding {
				buf = append(buf, left...)
			}
			buf = append(buf, key...)
			key = []byte{}
			expanding = true
			i += len(left)
			continue
		}
		if len(text) >= i+len(right) && expanding && bytes.Equal([]byte(text[i:i+len(right)]), right) {
			val, err := f(string(key))
			if err != nil {
				return "", err
			}
			buf = append(buf, val...)
			key = []byte{}
			expanding = false
			i += len(right)
			continue
		}
		if expanding {
			key = append(key, text[i])
		} else {
			buf = append(buf, text[i])
		}
		i++
	}
	if expanding {
		buf = append(buf, left...)
		buf = append(buf, key...)
	}
	return string(buf), nil
}

func GetEnv(key string) string {
	var k, def = key, ""
	if n := strings.Index(key, ":"); n > 0 {
		k = key[:n]
		def = key[n+1:]
	}
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
