package usecase

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduler_ScheduleAndStart(t *testing.T) {
	scheduler := NewScheduler()

	var counter int32

	err := scheduler.Schedule("test-task", 10*time.Millisecond, func(ctx context.Context) {
		atomic.AddInt32(&counter, 1)
	})

	if err != nil {
		t.Fatalf("unexpected error scheduling task: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler.Start(ctx)

	// Wait enough time for a few ticks
	time.Sleep(35 * time.Millisecond)

	scheduler.Stop()

	val := atomic.LoadInt32(&counter)
	// Expect around 3 ticks
	if val < 2 || val > 4 {
		t.Errorf("expected counter to be around 3, got %d", val)
	}
}

func TestScheduler_ScheduleAfterStart(t *testing.T) {
	scheduler := NewScheduler()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler.Start(ctx)

	var counter int32

	err := scheduler.Schedule("test-task-2", 10*time.Millisecond, func(ctx context.Context) {
		atomic.AddInt32(&counter, 1)
	})

	if err != nil {
		t.Fatalf("unexpected error scheduling task: %v", err)
	}

	time.Sleep(25 * time.Millisecond)

	scheduler.Stop()

	val := atomic.LoadInt32(&counter)
	if val < 1 || val > 3 {
		t.Errorf("expected counter to be around 2, got %d", val)
	}
}

func TestScheduler_Remove(t *testing.T) {
	scheduler := NewScheduler()

	var counter int32

	err := scheduler.Schedule("test-task-3", 10*time.Millisecond, func(ctx context.Context) {
		atomic.AddInt32(&counter, 1)
	})

	if err != nil {
		t.Fatalf("unexpected error scheduling task: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler.Start(ctx)

	time.Sleep(15 * time.Millisecond)

	err = scheduler.Remove("test-task-3")
	if err != nil {
		t.Fatalf("unexpected error removing task: %v", err)
	}

	// Wait to ensure no more ticks happen
	valBefore := atomic.LoadInt32(&counter)
	time.Sleep(20 * time.Millisecond)
	valAfter := atomic.LoadInt32(&counter)

	if valBefore != valAfter {
		t.Errorf("expected counter to not change after remove, but got %d -> %d", valBefore, valAfter)
	}

	scheduler.Stop()
}
