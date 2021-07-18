package cron

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

// Hub definition
type Hub struct {
	mu    sync.Mutex
	crons map[config]*cron.Cron
}

// NewHub creates a new hub.
func NewHub() *Hub {
	return &Hub{
		crons: make(map[config]*cron.Cron),
	}
}

// AddSchedule adds a schedule.
func (h *Hub) AddSchedule(s *Schedule) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	cfg := s.cronConfig()
	c := h.crons[cfg]
	if c == nil {
		c = cfg.build()
		h.crons[cfg] = c
	}

	if _, err := c.AddJob(s.cron.Express, s.job()); err != nil {
		return err
	}
	return nil
}

// Run starts all of schedules.
func (h *Hub) Start() error {
	for _, c := range h.crons {
		c.Start()
	}
	return nil
}

// Close shuts schedules down. Force to close the hub
// by the context with timeout.
func (h *Hub) Close(ctx context.Context) error {
	var g sync.WaitGroup
	for _, c := range h.crons {
		g.Add(1)
		go func(c *cron.Cron) {
			ctx := c.Stop()
			<-ctx.Done()
			g.Done()
		}(c)
	}
	closed := make(chan struct{})
	go func() {
		g.Wait()
		close(closed)
	}()

	select {
	case <-closed:
		return nil
	case <-ctx.Done():
		return ErrForceClose
	}
}
