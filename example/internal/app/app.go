package app

import (
	"context"

	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp"
	"github.com/go-chocolate/chocolate/pkg/chocolate/chocorpc"
	"github.com/go-chocolate/chocolate/pkg/chocolate/service"
	"github.com/go-chocolate/configuration/configuration"

	"github.com/go-chocolate/chocolate/example/internal/app/config"
	"github.com/go-chocolate/chocolate/example/internal/app/dependency"
	"github.com/go-chocolate/chocolate/example/internal/entrance/http"
	"github.com/go-chocolate/chocolate/example/internal/module"
)

var ctx, cancel = context.WithCancel(context.Background())

func Run() {
	var cfg config.Config
	if err := configuration.Load(&cfg); err != nil {
		panic(err)
	}
	if err := dependency.Setup(cfg); err != nil {
		panic(err)
	}
	if err := module.Setup(); err != nil {
		panic(err)
	}

	httpsrv := chocohttp.NewServer(cfg.HTTP)
	httpsrv.SetRouter(http.Router())

	rpcsrv := chocorpc.NewServer(cfg.RPC)

	group := service.Group(httpsrv, rpcsrv)

	if err := group.Run(ctx); err != nil {
		panic(err)
	}

}

func Shutdown() {
	cancel()

}
