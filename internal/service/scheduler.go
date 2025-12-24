package service

import (
	"netdash/internal/logger"
	"netdash/internal/repository"

	"github.com/robfig/cron/v3"
)

type SchedulerService interface {
	Start()
	UpdateSchedule()
}

type schedulerService struct {
	cron  *cron.Cron
	repo  repository.Repository
	speed SpeedtestService
	jobID cron.EntryID
}

func NewSchedulerService(repo repository.Repository, speed SpeedtestService) SchedulerService {
	return &schedulerService{
		cron:  cron.New(),
		repo:  repo,
		speed: speed,
	}
}

func (s *schedulerService) Start() {
	s.cron.Start()
	s.UpdateSchedule()
}

func (s *schedulerService) UpdateSchedule() {
	if s.jobID != 0 {
		s.cron.Remove(s.jobID)
		s.jobID = 0
	}

	conf, err := s.repo.GetConfig()
	if err != nil {
		return
	}

	if conf.CronSchedule == "manual" || conf.CronSchedule == "" {
		logger.Log("CRON", "Automated schedule disabled")
		return
	}

	id, err := s.cron.AddFunc(conf.CronSchedule, func() {
		logger.Log("CRON", "Triggering scheduled speedtest")
		_, _ = s.speed.Run()
	})

	if err != nil {
		logger.Log("CRON", "Error parsing schedule '%s': %v", conf.CronSchedule, err)
	} else {
		s.jobID = id
		logger.Log("CRON", "Schedule updated: %s", conf.CronSchedule)
	}
}