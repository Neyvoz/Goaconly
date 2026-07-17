package usecase

import (
	"context"
	"fmt"
	"goaconly/internal/domain"
	"goaconly/internal/worker"
	"sync"
	"time"
)

type SchedulerUseCase interface {
	Run(ctx context.Context)
	AddTarget(ctx context.Context, target domain.Target) error
	RemoveTarget(id int64) error
	UpdateInterval(id int64, d time.Duration) error
}
type WorkerPool interface {
	Submit(job worker.Job[domain.CheckJob]) bool
}
type scheduler struct {
	target map[int64]*domain.TickerEntry
	mu     sync.RWMutex
	cmdCh  chan domain.SchedulerCmd
	pool   WorkerPool
	repo   CheckResultRepository
}

func (s *scheduler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.stopAll()
			return
		case cmd := <-s.cmdCh:
			s.handleCmd(ctx, cmd)
		}
	}
}

func (s *scheduler) spawnTicker(ctx context.Context, t domain.Target, d time.Duration) {
	tickCtx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(d)
	s.target[t.ID] = &domain.TickerEntry{Ticker: ticker, Cancel: cancel, Target: t, Interval: d}
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-tickCtx.Done():
				return
			case <-ticker.C:
				s.pool.Submit(worker.Job[domain.CheckJob]{
					Payload: domain.CheckJob{Target: t},
				})
			}
		}
	}()
}

func (s *scheduler) stopAll() {
	for _, entry := range s.target {
		entry.Cancel()
		entry.Ticker.Stop()
	}
}

func (s *scheduler) handleCmd(ctx context.Context, cmd domain.SchedulerCmd) {
	switch cmd.Type {
	case domain.CmdAdd:
		if cmd.Target != nil {
			s.spawnTicker(ctx, *cmd.Target, time.Duration(cmd.Target.CheckInterval)*time.Minute)
		}
	case domain.CmdRemove:
		if entry, ok := s.target[cmd.TargetID]; ok {
			entry.Cancel()
			entry.Ticker.Stop()
			delete(s.target, cmd.TargetID)
		}
	case domain.CmdUpdateInterval:
		if entry, ok := s.target[cmd.TargetID]; ok {
			t := entry.Target
			entry.Cancel()
			entry.Ticker.Stop()
			delete(s.target, cmd.TargetID)
			s.spawnTicker(ctx, t, cmd.Interval)
		}
	}
}

func (s *scheduler) AddTarget(ctx context.Context, t domain.Target) error {
	select {
	case s.cmdCh <- domain.SchedulerCmd{Type: domain.CmdAdd, Target: &t}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *scheduler) RemoveTarget(id int64) error {
	select {
	case s.cmdCh <- domain.SchedulerCmd{Type: domain.CmdRemove, TargetID: id}:
		return nil
	default:
		return fmt.Errorf("scheduler: command channel full")
	}
}

func (s *scheduler) UpdateInterval(id int64, d time.Duration) error {
	select {
	case s.cmdCh <- domain.SchedulerCmd{Type: domain.CmdUpdateInterval, TargetID: id, Interval: d}:
		return nil
	default:
		return fmt.Errorf("scheduler: command channel full")
	}
}

func NewScheduler(pool WorkerPool, repo CheckResultRepository) SchedulerUseCase {
	return &scheduler{
		target: make(map[int64]*domain.TickerEntry),
		cmdCh:  make(chan domain.SchedulerCmd, 256),
		pool:   pool,
		repo:   repo,
	}
}
