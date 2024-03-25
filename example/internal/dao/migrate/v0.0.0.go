package migrate

import (
	"gorm.io/gorm"

	"github.com/go-chocolate/chocolate/example/internal/dao/model"
)

var (
	v0_0_0 = migration{
		version: "v0.0.0",
		migrate: func(db *gorm.DB) error {
			m := db.Migrator()
			if err := m.CreateTable(&model.User{}); err != nil {
				return err
			}
			return nil
		},
		rollback: func(db *gorm.DB) error {
			m := db.Migrator()
			if err := m.DropTable(&model.User{}); err != nil {
				return err
			}
			return nil
		},
	}
)
