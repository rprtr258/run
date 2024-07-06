package run

import (
	"context"
	"errors"
)

type Lifecycle struct {
	closers []func(context.Context) error
}

func New() *Lifecycle {
	return &Lifecycle{}
}

func (lc *Lifecycle) AddCE(closer func(context.Context) error) {
	lc.closers = append(lc.closers, closer)
}

func (lc *Lifecycle) AddE(closer func() error) {
	lc.AddCE(func(context.Context) error {
		return closer()
	})
}

func (lc *Lifecycle) Add(closer func()) {
	lc.AddCE(func(context.Context) error {
		closer()
		return nil
	})
}

func (lc *Lifecycle) Close(ctx context.Context) error {
	e := ""
	for i := len(lc.closers) - 1; i >= 0; i-- {
		if err := lc.closers[i](ctx); err != nil {
			e += err.Error()
		}
	}
	if e != "" {
		return errors.New(e)
	}
	return nil
}
