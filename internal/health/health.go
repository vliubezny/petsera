package health

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Checker interface {
	Check(ctx context.Context) error
}

type CheckFunc func(ctx context.Context) error

func (f CheckFunc) Check(ctx context.Context) error {
	return f(ctx)
}

func SetupChecks(checkers ...Checker) CheckFunc {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		gr, ctx := errgroup.WithContext(ctx)

		for i := range checkers {
			c := checkers[i]
			gr.Go(func() error {
				if err := c.Check(ctx); err != nil {
					logrus.WithError(err).Error("health check failed")
					return err
				}
				return nil
			})
		}

		if err := gr.Wait(); err != nil {
			return err
		}

		return nil
	}
}
