package gscheduler

import (
	"sort"
	"time"
)

type Entry struct {
	Schedule Schedule
	Job      Job
	Next     time.Time
	Prev     time.Time
}

// 时间排序
type byTime []*Entry

func (b byTime) Len() int      { return len(b) }
func (b byTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// 时间最近的任务
func (b byTime) Less(i, j int) bool {

	if b[i].Next.IsZero() {
		return false
	}
	if b[j].Next.IsZero() {
		return true
	}
	return b[i].Next.Before(b[j].Next)
}

type Job interface {
	Run()
}

type Cron struct {
	entries []*Entry
	running bool
	add     chan *Entry
	stop    chan struct{}
}

func New() *Cron {
	return &Cron{
		stop: make(chan struct{}),
		add:  make(chan *Entry),
	}
}

func (c *Cron) Start() {
	c.running = true
	// 持续处理任务
	go c.run()
}

func (c *Cron) Add(s Schedule, j Job) {

	entry := &Entry{
		Schedule: s,
		Job:      j,
	}

	if !c.running {
		c.entries = append(c.entries, entry)
		return
	}
	c.add <- entry
}

// 停止
func (c *Cron) Stop() {

	if !c.running {
		return
	}
	c.running = false
	c.stop <- struct{}{}
}

var after = time.After

func (c *Cron) run() {

	var effective time.Time
	now := time.Now().Local()
	for _, e := range c.entries {
		e.Next = e.Schedule.Next(now)
	}

	for {
		sort.Sort(byTime(c.entries))
		if len(c.entries) > 0 {
			effective = c.entries[0].Next
		} else {
			effective = now.AddDate(15, 0, 0)
		}

		select {
		case now = <-after(effective.Sub(now)):
			// 相同时间的任务
			for _, entry := range c.entries {
				if entry.Next != effective {
					break
				}
				entry.Prev = now
				entry.Next = entry.Schedule.Next(now)
				go entry.Job.Run()
			}
		case e := <-c.add:
			e.Next = e.Schedule.Next(time.Now())
			c.entries = append(c.entries, e)
		case <-c.stop:
			return
		}
	}
}

func (c Cron) Entries() []*Entry {
	return c.entries
}
