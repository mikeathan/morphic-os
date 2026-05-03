package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ScheduledTask defines a background task that runs at a given interval.
type ScheduledTask struct {
	ID       string
	Interval time.Duration
	Action   func(ctx context.Context)
	ticker   *time.Ticker
	cancel   context.CancelFunc
}

// Scheduler defines the interface for scheduling tasks.
type Scheduler interface {
	Schedule(id string, interval time.Duration, action func(ctx context.Context)) error
	Remove(id string) error
	Start(ctx context.Context)
	Stop()
}

type taskScheduler struct {
	tasks map[string]*ScheduledTask
	mu    sync.Mutex
	ctx   context.Context
}

// NewScheduler creates a new task scheduler.
func NewScheduler() Scheduler {
	return &taskScheduler{
		tasks: make(map[string]*ScheduledTask),
	}
}

// Schedule adds a new task to the scheduler. If the scheduler is already started,
// the task starts running immediately.
func (s *taskScheduler) Schedule(id string, interval time.Duration, action func(ctx context.Context)) error {
	if interval <= 0 {
		return fmt.Errorf("invalid interval: must be > 0")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; exists {
		return fmt.Errorf("task with ID %s already exists", id)
	}

	task := &ScheduledTask{
		ID:       id,
		Interval: interval,
		Action:   action,
	}

	s.tasks[id] = task

	// If the scheduler is running, start the task immediately
	if s.ctx != nil {
		s.startTask(task)
	}

	return nil
}

// Remove stops and removes a task by its ID.
func (s *taskScheduler) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil // Or return error
	}

	if task.cancel != nil {
		task.cancel()
	}
	if task.ticker != nil {
		task.ticker.Stop()
	}

	delete(s.tasks, id)
	return nil
}

func (s *taskScheduler) startTask(task *ScheduledTask) {
	ctx, cancel := context.WithCancel(s.ctx)
	task.cancel = cancel
	task.ticker = time.NewTicker(task.Interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-task.ticker.C:
				task.Action(ctx)
			}
		}
	}()
}

// Start begins execution of all scheduled tasks.
func (s *taskScheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx != nil {
		return // already started
	}
	s.ctx = ctx

	for _, task := range s.tasks {
		s.startTask(task)
	}
}

// Stop stops all tasks and the scheduler itself.
func (s *taskScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if task.cancel != nil {
			task.cancel()
		}
		if task.ticker != nil {
			task.ticker.Stop()
		}
	}
	s.ctx = nil
}
