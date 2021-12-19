package cron

import (
	"github.com/robfig/cron/v3"
)

func withRecover(job cron.Job) cron.Job {
	return cron.Recover(cron.DiscardLogger)(job)
}

func withSkipIfStillRunning(job cron.Job) cron.Job {
	return cron.SkipIfStillRunning(cron.DiscardLogger)(job)
}

func withDelayIfStillRunning(job cron.Job) cron.Job {
	return cron.DelayIfStillRunning(cron.DiscardLogger)(job)
}
