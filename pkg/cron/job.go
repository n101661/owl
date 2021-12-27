package cron

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type Values map[string]string

type Job interface {
	Name() string
	Run(context.Context, Values) error
}

func newJob(j Job) cron.Job {
	return cron.FuncJob(func() {
		startTime := time.Now()
		rid := xid.New().String()

		logger := zap.L().With(
			zap.String("request_id", rid),
			zap.String("name", j.Name()),
		)

		logger.Info("start the schedule")

		vs := Values{
			RequestID: rid,
		}

		err := j.Run(context.Background(), vs)
		ms := time.Now().Sub(startTime).Milliseconds()
		if err != nil {
			logger.Error("finish the schedule with error", zap.Int64("time_ms", ms), zap.Error(err))
		} else {
			logger.Info("finish the schedule", zap.Int64("time_ms", ms))
		}
	})
}
