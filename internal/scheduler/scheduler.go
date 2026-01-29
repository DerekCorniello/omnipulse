// Package scheduler provides background task scheduling for data fetching.
package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/omnipulse/omnipulse/internal/config"
)

// Scheduler manages periodic background tasks.
type Scheduler struct {
	config   config.SchedulerConfig
	tasks    []*Task
	mu       sync.RWMutex
	running  bool
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// Task represents a scheduled task.
type Task struct {
	Name     string
	Interval time.Duration
	Fn       func(ctx context.Context) error
	ticker   *time.Ticker
	stopChan chan struct{}
}

// NewScheduler creates a new Scheduler.
func NewScheduler(cfg config.SchedulerConfig) *Scheduler {
	return &Scheduler{
		config:   cfg,
		tasks:    make([]*Task, 0),
		stopChan: make(chan struct{}),
	}
}

// AddTask adds a new task to the scheduler.
func (s *Scheduler) AddTask(name string, interval time.Duration, fn func(ctx context.Context) error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &Task{
		Name:     name,
		Interval: interval,
		Fn:       fn,
		stopChan: make(chan struct{}),
	}
	s.tasks = append(s.tasks, task)

	// If scheduler is already running, start the new task
	if s.running {
		s.startTask(task)
	}
}

// Start begins executing all scheduled tasks.
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}

	s.running = true
	for _, task := range s.tasks {
		s.startTask(task)
	}

	log.Printf("Scheduler started with %d tasks", len(s.tasks))
}

// startTask starts a single task's ticker loop.
func (s *Scheduler) startTask(task *Task) {
	task.ticker = time.NewTicker(task.Interval)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer task.ticker.Stop()

		// Run immediately on start
		if err := task.Fn(context.Background()); err != nil {
			log.Printf("Task %s error: %v", task.Name, err)
		}

		for {
			select {
			case <-task.ticker.C:
				log.Printf("Running task: %s", task.Name)
				if err := task.Fn(context.Background()); err != nil {
					log.Printf("Task %s error: %v", task.Name, err)
				}
			case <-task.stopChan:
				log.Printf("Task %s stopped", task.Name)
				return
			case <-s.stopChan:
				return
			}
		}
	}()
}

// Stop stops all scheduled tasks.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	// Signal all tasks to stop
	close(s.stopChan)
	for _, task := range s.tasks {
		close(task.stopChan)
	}

	// Wait for all tasks to finish
	s.wg.Wait()
	s.running = false

	log.Println("Scheduler stopped")
}

// IsRunning returns whether the scheduler is running.
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetTasks returns a list of task names and their intervals.
func (s *Scheduler) GetTasks() []TaskInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info := make([]TaskInfo, len(s.tasks))
	for i, task := range s.tasks {
		info[i] = TaskInfo{
			Name:     task.Name,
			Interval: task.Interval,
		}
	}
	return info
}

// TaskInfo contains information about a scheduled task.
type TaskInfo struct {
	Name     string        `json:"name"`
	Interval time.Duration `json:"interval"`
}

// RunTaskNow executes a specific task immediately.
func (s *Scheduler) RunTaskNow(ctx context.Context, name string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		if task.Name == name {
			return task.Fn(ctx)
		}
	}
	return nil
}
