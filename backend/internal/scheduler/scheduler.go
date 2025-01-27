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
}

func NewCronManager(logger *zap.SugaredLogger) *CronManager {
	return &CronManager{
		scheduler: gocron.NewScheduler(time.UTC),
		logger:    logger,
	}
}

type CronJob string

const (
	HourlyJob CronJob = "HOURLY"
	DailyJob  CronJob = "DAILY"
)

func (cm *CronManager) RegisterJob(interval CronJob, task func(), timeUTC ...string) error {
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
	cm.logger.Info("Starting the cron manager...")
	cm.scheduler.StartAsync()
}

func (cm *CronManager) Stop() {
	cm.logger.Info("Stopping the cron manager...")
	cm.scheduler.Stop()
}
