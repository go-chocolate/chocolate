package service

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Service interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Services []Service

func (s Services) Run(ctx context.Context) error {
	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	for i := range s {
		group.Go(s.run(ctx, i))
	}
	return group.Wait()
}

func (s Services) run(ctx context.Context, i int) func() error {
	return func() error {
		return s[i].Run(ctx)
	}
}

func (s Services) Shutdown(ctx context.Context) (err error) {
	for i := range s {
		if e := s[i].Shutdown(ctx); e != nil {
			err = e
		}
	}
	return
}

func Group(services ...Service) Service {
	return Services(services)
}
