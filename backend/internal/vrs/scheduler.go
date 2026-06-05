package vrs

import (
	"context"
	"time"

	"github.com/Nishal77/resona/backend/pkg/config"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	cron    *cron.Cron
	service *Service
}

func NewScheduler(svc *Service) *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		service: svc,
	}
}

func (s *Scheduler) Start() {
	// VRS recalculation every 30 minutes
	s.cron.AddFunc("0,30 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		log.Info().Msg("vrs recalculation started")
		if err := s.service.RecalculateAll(ctx); err != nil {
			log.Error().Err(err).Msg("vrs recalculation failed")
			return
		}

		if err := s.service.UpdateTrendingTags(ctx); err != nil {
			log.Error().Err(err).Msg("trending tags update failed")
		}

		if err := s.service.CreateTrendingNotifications(ctx, config.App.VRSTrendingThreshold); err != nil {
			log.Error().Err(err).Msg("trending notifications failed")
		}

		log.Info().Msg("vrs recalculation complete")
	})

	// Snap of the week — every Monday at 00:00
	s.cron.AddFunc("0 0 * * 1", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := s.service.SetSnapOfWeek(ctx); err != nil {
			log.Error().Err(err).Msg("snap of week failed")
		} else {
			log.Info().Msg("snap of week updated")
		}
	})

	s.cron.Start()
	log.Info().Msg("vrs scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
