package args

import (
	"os"
	"strings"
)

type KV struct {
	Key string
	Val string
}

type Options []*KV

func (o Options) KVs() KVs {
	m := make(map[string][]string)
	for _, v := range o {
		m[v.Key] = append(m[v.Key], v.Val)
	}
	return m
}

type KVs map[string][]string

func (kv KVs) Get(key string, defaults ...string) string {
	if val := kv[key]; len(val) > 0 {
		return val[0]
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (kv KVs) GetAlias(alias []string, defaults ...string) string {
	for _, k := range alias {
		if val := kv[k]; len(val) > 0 {
			return val[0]
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (kv KVs) Exist(alias ...string) bool {
	for _, k := range alias {
		if _, ok := kv[k]; ok {
			return true
		}
	}
	return false
}

type Command struct {
	Command string
	Options KVs
}

func Parse() *Command {
	n := len(os.Args)
	if n < 2 {
		return &Command{}
	}
	cmd := new(Command)
	command := os.Args[1]
	if command != "" && command[0] != '-' {
		cmd.Command = command
	}
	options := Options{}

	i := 2
	for i < n {
		key := strings.TrimSpace(os.Args[i])
		i++
		if key == "" || key[0] != '-' {
			continue
		}
		opt := &KV{Key: key[1:]}
		options = append(options, opt)
		if i >= n {
			break
		}
		val := strings.TrimSpace(os.Args[i])
		i++
		if val != "" && val[0] == '-' {
			options = append(options, &KV{Key: val[1:]})
		} else {
			opt.Val = val
		}
	}
	cmd.Options = options.KVs()
	return cmd
}
