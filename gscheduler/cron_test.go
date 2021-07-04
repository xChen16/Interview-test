package gscheduler

import (
	"testing"
	"time"
)

const OneMinute = 1*time.Minute + 10*time.Second

func TestNoEntries(t *testing.T) {
	cron := New()
	cron.Start()

	select {
	case <-time.After(OneMinute):
		t.FailNow()
	case <-stop(cron):
	}
}

// 模拟延迟，触发任务
func TestNoPhantomJobs(t *testing.T) {
	entry := 0

	after = func(d time.Duration) <-chan time.Time {
		entry++
		return time.After(d)
	}
	defer func() {
		after = time.After
	}()

	cron := New()
	cron.Start()
	defer cron.Stop()

	time.Sleep(1 * time.Second)

	if entry > 1 {
		t.Errorf("phantom job had run %d time(s).", entry)
	}
}

// 运行时添加任务
func TestAddBeforeRun(t *testing.T) {
	done := make(chan struct{})
	cron := New()
	cron.Start()
	defer cron.Stop()

	select {
	case <-time.After(OneMinute):
		t.FailNow()
	case <-done:
	}
}

func stop(c *Cron) <-chan bool {
	done := make(chan bool)
	go func() {
		c.Stop()
		done <- true
	}()
	return done
}
