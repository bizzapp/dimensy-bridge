package jobs

import (
	"context"
	"dimensy-bridge/internal/config"
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
	deps *config.AppDependencies
}

// NewScheduler buat instance baru scheduler
func NewScheduler(deps *config.AppDependencies) *Scheduler {
	c := cron.New(cron.WithSeconds()) // pakai format cron dengan detik
	return &Scheduler{
		cron: c,
		deps: deps,
	}
}

// Register semua cron jobs di sini
func (s *Scheduler) Register() error {
	// contoh job tiap 10 detik
	_, err := s.cron.AddFunc("*/100000000 * * * * *", func() {
		log.Println("🔄 Running sample job: cek sesuatu di database...")
		// kamu bisa akses s.deps untuk pakai service/DB
	})
	if err != nil {
		return err
	}

	// contoh job tiap jam 2 pagi
	_, err = s.cron.AddFunc("0 0 2 * * *", func() {
		log.Println("🌙 Running nightly cleanup job...")
	})
	if err != nil {
		return err
	}

	return nil
}

// Start jalanin cron
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop berhentiin cron dengan context timeout
func (s *Scheduler) Stop(ctx context.Context) error {
	// Graceful stop
	s.cron.Stop()
	log.Println("✅ Scheduler stopped")
	return nil
}
