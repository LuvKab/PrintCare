package main

import (
	"sync"
	"time"
)

type Scheduler struct {
	app     *App
	ticker  *time.Ticker
	done    chan struct{}
	mu      sync.Mutex
	running bool
}

func NewScheduler(app *App) *Scheduler {
	return &Scheduler{app: app}
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return
	}
	s.running = true
	s.done = make(chan struct{})
	s.ticker = time.NewTicker(15 * time.Minute)
	go s.loop()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	s.ticker.Stop()
	close(s.done)
}

func (s *Scheduler) Reset() {
	s.Stop()
	s.Start()
}

func (s *Scheduler) loop() {
	s.check()
	for {
		select {
		case <-s.done:
			return
		case <-s.ticker.C:
			s.check()
		}
	}
}

func (s *Scheduler) check() {
	cfg := s.app.config
	if !cfg.Enabled || cfg.PrinterName == "" || cfg.ImagePath == "" || cfg.IntervalDays <= 0 {
		return
	}

	interval := time.Duration(cfg.IntervalDays) * 24 * time.Hour
	if !cfg.LastPrint.IsZero() && time.Since(cfg.LastPrint) < interval {
		return
	}

	s.app.logger.Info("Scheduled print triggered, waiting 60s before printing")

	select {
	case <-time.After(60 * time.Second):
	case <-s.done:
		return
	}

	if err := PrintImage(cfg.PrinterName, cfg.ImagePath, cfg.PaperSource); err != nil {
		s.app.logger.Error("Scheduled print failed: %v", err)
		return
	}

	cfg.LastPrint = time.Now()
	saveConfig(cfg, s.app.logger)
	s.app.logger.Info("Scheduled print completed successfully")
}
