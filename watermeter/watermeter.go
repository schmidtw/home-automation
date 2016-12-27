// Package watermeter provide simple implementation of a water meter
package watermeter

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	time  time.Time
	total uint64
}

type Watermeter struct {
	Initial uint64 // 1/1000 gallons
	Timeout time.Duration
	Usage   func(gallons uint64)

	total  uint64
	events list.List
	mutex  sync.Mutex
}

func (w *Watermeter) Init() *Watermeter {
	w.total = w.Initial
	w.mutex = sync.Mutex{}
	w.events.Init()

	e := new(entry)
	e.time = time.Now()
	e.total = w.total
	w.events.PushFront(e)

	return w
}

func (w *Watermeter) GetFlow(duration time.Duration) float64 {
	now := time.Now()
	then := now.Add(-duration)

	end := entry{time: now, total: w.total}
	start := entry{time: now, total: w.total}

	w.mutex.Lock()

	item := w.events.Front()

	for nil != item {
		e := item.Value.(*entry)
		if then.Before(e.time) {
			start.time = e.time
			start.total = e.total
			item = item.Next()
		} else {
			item = nil
		}
	}
	w.mutex.Unlock()

	volume_delta := end.total - start.total
	return float64(volume_delta) / 1000 / duration.Minutes()
}

func (w *Watermeter) Gallons() uint64 {
	return w.total / 1000
}

func (w *Watermeter) Update(mGallons uint) {
	now := time.Now()
	prune := now.Add(-w.Timeout)

	w.mutex.Lock()
	before := w.total / 1000
	w.total += uint64(mGallons)
	after := w.total / 1000

	e := new(entry)
	e.time = now
	e.total = w.total
	w.events.PushFront(e)

	done := false
	for false == done {
		item := w.events.Back()
		e := item.Value.(*entry)
		if e.time.Before(prune) {
			w.events.Remove(item)
		} else {
			done = true
		}
	}

	w.mutex.Unlock()

	if (after - before) > 0 {
		if nil != w.Usage {
			go (w.Usage)(after)
		}
	}
}
