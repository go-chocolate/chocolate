package config

import (
	"github.com/go-chocolate/chocolate/pkg/chocolate/basic"
	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp"
	"github.com/go-chocolate/chocolate/pkg/chocolate/chocorpc"
)

type Config struct {
	basic.Config
	HTTP chocohttp.Config
	RPC  chocorpc.Config
}
