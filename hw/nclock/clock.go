package nclock

import (
	"io"
	"sync"
	"time"
)

type Clock struct {
	l  sync.RWMutex
	c  *sync.Cond
	hz int64

	stop chan struct{}
}

func New(hz int64) *Clock {
	c := &Clock{hz: hz, stop: make(chan struct{})}
	c.c = sync.NewCond(c.l.RLocker())

	go c.run()

	return c
}

func (c *Clock) run() {
	t := time.NewTicker(time.Second / c.hz)

	for {
		select {
		case <-t.C:
			c.c.Broadcast()
		case <-c.stop:
			c.l.Lock()
			defer c.l.Unlock()
			c.c.Broadcast()
			c.c = nil
			return
		}
	}
}

func (c *Clock) Wait() error {
	c.l.RLock()
	defer c.l.RUnlock()

	if c.c == nil {
		return io.EOF
	}

	c.c.Wait()
	return nil
}
