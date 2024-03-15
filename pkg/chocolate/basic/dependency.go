package basic

import (
	"github.com/go-chocolate/chocolate/pkg/database/orm"
	"github.com/go-chocolate/chocolate/pkg/kv"
	"github.com/go-chocolate/chocolate/pkg/telemetry"
	"gorm.io/gorm"
)

type Dependency struct {
	DB        *gorm.DB
	KVStorage kv.Storage
}

func (dep *Dependency) Setup(c Config) error {
	var err error
	if dep.DB, err = orm.Open(c.Database); err != nil {
		return err
	}
	if dep.KVStorage, err = kv.New(c.KV); err != nil {
		return err
	}
	
	telemetry.Setup(c.Telemetry)

	return nil
}
