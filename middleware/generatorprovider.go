package middleware

import (
	"github.com/labstack/echo/v4"
	"uidGenerator/generator"
	"uidGenerator/timeprovider"
)

func GeneratorProvider(workerId int64, provider timeprovider.TimeProvider) echo.MiddlewareFunc {
	workers := make(chan *generator.WorkerVariant, generator.ThreadCap)
	var i int64
	for i = 1; i <= generator.ThreadCap; i++ {
		worker := &generator.WorkerVariant{
			WorkerID:     workerId,
			ThreadId:     i,
			TimeProvider: provider,
		}
		workers <- worker
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			worker := <-workers
			c.Set("worker", worker)
			c.Logger().Debugf("worker %d", worker.WorkerID)
			defer func() {
				workers <- worker
			}()
			return next(c)
		}
	}
}
