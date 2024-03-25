package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var (
	migrations = []*gormigrate.Migration{
		{
			ID:       v0_0_0.version,
			Migrate:  v0_0_0.migrate,
			Rollback: v0_0_0.rollback,
		},
	}

	options = gormigrate.DefaultOptions
)

type migration struct {
	version  string
	migrate  func(db *gorm.DB) error
	rollback func(db *gorm.DB) error
}

func Migrate(db *gorm.DB, ids ...string) error {
	m := gormigrate.New(db, options, migrations)
	if len(ids) > 0 {
		return m.MigrateTo(ids[0])
	} else {
		return m.Migrate()
	}
}

func Rollback(db *gorm.DB, ids ...string) error {
	m := gormigrate.New(db, options, migrations)
	if len(ids) > 0 {
		return m.RollbackTo(ids[0])
	} else {
		return m.RollbackLast()
	}
}
