package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type CronManager struct {
	scheduler *gocron.Scheduler
	logger    *zap.SugaredLogger
	enabled   bool
}

func NewCronManager(enabled bool, logger *zap.SugaredLogger) *CronManager {
	return &CronManager{
		scheduler: gocron.NewScheduler(time.UTC),
		logger:    logger,
		enabled:   enabled,
	}
}

type CronJob string

const (
	HourlyJob CronJob = "HOURLY"
	DailyJob  CronJob = "DAILY"
)

func (cm *CronManager) RegisterJob(interval CronJob, task func(), timeUTC ...string) error {
	if !cm.enabled {
		cm.logger.Info("Cron jobs are disabled; skipping job registration.")
		return nil
	}

	switch interval {
	case HourlyJob:
		_, err := cm.scheduler.Every(1).Hour().Do(task)
		if err != nil {
			return err
		}
	case DailyJob:
		if len(timeUTC) > 0 {
			_, err := cm.scheduler.Every(1).Day().At(timeUTC[0]).Do(task)
			if err != nil {
				return err
			}
			cm.logger.Infof("Job registered for interval: %s at %v", interval, timeUTC[0])
		} else {
			_, err := cm.scheduler.Every(1).Day().At("00:00").Do(task)
			if err != nil {
				return err
			}
		}
	default:
		cm.logger.Errorf("unsupported interval: %s", interval)
		return fmt.Errorf("unsupported interval: %s", interval)
	}

	cm.logger.Infof("Job registered for interval: %s at %v", interval, timeUTC)
	return nil
}

func (cm *CronManager) Start() {
	if !cm.enabled {
		cm.logger.Info("Cron jobs are disabled; not starting the cron manager.")
		return
	}
	cm.logger.Info("Starting the cron manager...")
	cm.scheduler.StartAsync()
}

func (cm *CronManager) Stop() {
	if !cm.enabled {
		cm.logger.Info("Cron jobs are disabled; nothing to stop.")
		return
	}
	cm.logger.Info("Stopping the cron manager...")
	cm.scheduler.Stop()
}
