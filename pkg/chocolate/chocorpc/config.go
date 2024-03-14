package chocorpc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type ByteText string

func (b ByteText) Value() int {
	text := strings.ToLower(string(b))
	var num []byte
	var unit string
	for i, v := range text {
		if v >= '0' && v <= '9' {
			num = append(num, byte(v))
		} else {
			unit = text[i:]
			break
		}
	}
	if len(num) == 0 {
		return 0
	}
	val, _ := strconv.Atoi(string(num))
	switch unit {
	case "", "b":
		return val
	case "k", "kb":
		return val * 1024
	case "m", "mb":
		return val * 1024 * 1024
	case "g", "gb":
		return val * 1024 * 1024 * 1024
	case "t", "tb":
		return val * 1024 * 1024 * 1024 * 1024
	case "p", "pb":
		return val * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		panic(fmt.Errorf("invalid byte text '%s'", b))
	}
}

type Config struct {
	Name           string //服务名
	Addr           string //监听地址加端口
	Timeout        string
	MaxRecvMsgSize ByteText
	MaxSendMsgSize ByteText
	//Logger         LoggerConfig
}

func (c *Config) GetTimeout() time.Duration {
	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		panic(fmt.Errorf("chocorpc.Config.Timeout:cannot decode string value '%s' to duration, %v", c.Timeout, err))
	}
	return timeout
}

type LoggerConfig struct {
	Enable bool
}

func (c *Config) apply() []grpc.ServerOption {
	var options []grpc.ServerOption
	if c.Timeout != "" {
		options = append(options, grpc.ConnectionTimeout(c.GetTimeout()))
	}
	if c.MaxRecvMsgSize != "" {
		options = append(options, grpc.MaxRecvMsgSize(c.MaxRecvMsgSize.Value()))
	}
	if c.MaxSendMsgSize != "" {
		options = append(options, grpc.MaxSendMsgSize(c.MaxSendMsgSize.Value()))
	}
	return options
}
