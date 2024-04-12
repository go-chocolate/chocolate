package parameter

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-chocolate/chocolate/pkg/toolkit/envutil"
)

type Parameters map[string]string

func (p Parameters) Get(key string, defaultValue ...string) string {
	if value, ok := p[key]; ok {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (p Parameters) Set(key, value string) {
	p[key] = value
}

func (p Parameters) Append(with Parameters) {
	for key, value := range with {
		p[key] = value
	}
}

func (p Parameters) Expand(template string, withEnv ...bool) (string, error) {
	return envutil.Expand(template, func(key string) (string, error) {
		if n := strings.Index(key, "("); n > 0 && key[len(key)-1] == ')' {
			funcName := key[:n]
			args := key[n+1 : len(key)-1]
			if f, ok := _functions_[funcName]; ok {
				return f(args)
			} else {
				return "", fmt.Errorf("function not found: %s", funcName)
			}
		}
		var def string
		if n := strings.Index(key, ":"); n > 0 {
			def = key[n+1:]
			key = key[:n]
		}
		if val, ok := p[key]; ok {
			return val, nil
		}
		if len(withEnv) > 0 && withEnv[0] {
			val, ok := os.LookupEnv(key)
			if ok {
				return val, nil
			}
		}
		return def, nil
	})
}

func (p Parameters) ExpandWithEnv(template string) (string, error) {
	return p.Expand(template, true)
}
