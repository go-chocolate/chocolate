package dependency

import (
	"github.com/go-chocolate/chocolate/pkg/chocolate/basic"

	"github.com/go-chocolate/chocolate/example/internal/app/config"
)

var instance *Dependency

type Dependency struct {
	basic.Dependency
	Config config.Config
}

func (dep *Dependency) setup(c config.Config) error {
	dep.Config = c
	if err := dep.Dependency.Setup(c.Config); err != nil {
		return err
	}
	return nil
}

func Setup(c config.Config) error {
	instance = new(Dependency)
	return instance.setup(c)
}
