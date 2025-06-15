package monitor

import (
	"sync"
)

type Monitor struct {
	mu sync.Mutex
}

type Condition struct {
	mon   *Monitor
	queue []chan struct{}
}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (m *Monitor) Lock() {
	m.mu.Lock()
}

func (m *Monitor) Unlock() {
	m.mu.Unlock()
}

func (m *Monitor) NewCondition() *Condition {
	return &Condition{
		mon:   m,
		queue: make([]chan struct{}, 0),
	}
}

func (c *Condition) Signal() {
	if len(c.queue) > 0 {
		//pobiera kanał oczekującej go-routyny
		ch := c.queue[0]
		//usuwamy ten kanał z kolejki
		c.queue = c.queue[1:]
		//odblokuj tę go-routine
		ch <- struct{}{}
	}
}

func (c *Condition) Wait() {
	ch := make(chan struct{})
	c.queue = append(c.queue, ch)
	c.mon.mu.Unlock()
	<-ch //odbiór z kanału (blokujący)
	c.mon.mu.Lock()
}

func (c *Condition) SignalAll() {
	for len(c.queue) > 0 {
		c.Signal()
	}
	c.queue = make([]chan struct{}, 0)
}
