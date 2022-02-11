package cron

import (
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type DistributedTask struct {
	CheckInterval  time.Duration
	ActionInterval time.Duration
	LockKey        string
}

type CronJob struct {
	Name     string
	Interval time.Duration
	Function interface{}
	Params   []interface{}
}

type CronJobRunner struct {
	cr   *gocron.Scheduler
	Db   *gorm.DB
	Jobs []CronJob
}

func (c *CronJobRunner) Start() {
	c.cr = gocron.NewScheduler(time.UTC)

	for _, job := range c.Jobs {
		_, err := c.cr.Every(job.Interval).Do(job.Function, job.Params...)
		if err != nil {
			logrus.WithError(err).WithField("name", job.Name).Fatal("failed to start cron job")
		}
	}
	c.cr.StartAsync()
}

func (c *CronJobRunner) Stop() {
	c.cr.Stop()
}

func (c *CronJobRunner) Name() string {
	return "CronJobRunner"
}
