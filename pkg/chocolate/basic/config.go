package basic

import (
	"github.com/go-chocolate/chocolate/pkg/database/orm"
	"github.com/go-chocolate/chocolate/pkg/kv"
	"github.com/go-chocolate/chocolate/pkg/telemetry"

	"github.com/go-chocolate/chocolate/pkg/chocolate/logs"
)

type Config struct {
	Database  orm.Config
	KV        kv.Config
	Telemetry telemetry.Config
	Logger    logs.Config
}
