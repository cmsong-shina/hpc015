package hpc015

import (
	"fmt"
	"math"
	"time"
)

// counter help counting
//
// Due to sometime hpc015 send duplicated data,
// counter store recent data within 10 mins, and not count duplicated data.
//
// Use this just for a example,
// it store all count in one variable, not compitable with multiple device,
// and not thread safe and can cause panic.
type counter struct {
	count       int
	eventBuffer map[string]eventEntry
	// mux         *sync.Mutex
}

// Counter create new counter
//
// It run a go routine, which clear old events
func Counter(initCount int) *counter {
	counter := &counter{
		count:       initCount,
		eventBuffer: make(map[string]eventEntry),
		// mux:         &sync.Mutex{},
	}
	go counter.clearTicker()

	return counter
}

// Count a data
func (c *counter) Count(data *CacheData) {
	buf := make([]byte, 6, 6)
	buf[0] = data.Year
	buf[1] = data.Month
	buf[2] = data.Day
	buf[3] = data.Hour
	buf[4] = data.Minute
	buf[5] = data.Secound

	k := string(buf)

	_, ok := c.eventBuffer[k]
	if ok {
		if enableDebugMessage {
			fmt.Printf("- duplicated event: %d:%d:%d\n", buf[3], buf[4], buf[5])
		}
		return
	}

	c.count += int(data.DxIn)
	c.count -= int(data.Dxout)

	c.eventBuffer[k] = eventEntry{
		time.Now(),
		int(data.DxIn),
		int(data.Dxout),
	}
}

// Get current count
func (c *counter) Get() int {
	return c.count
}

// Set current count
func (c *counter) Set(count int) {
	c.count = count
}

// clearTicker excute clear every 1 min
func (c *counter) clearTicker() {
	t := time.NewTicker(time.Duration(time.Minute))
	for range t.C {
		c.clear()
	}
}

// clear events older than 10 mins
func (c *counter) clear() {
	deletedEntry := 0
	now := time.Now()

	for k, e := range c.eventBuffer {
		if math.Abs(e.Created.Sub(now).Minutes()) > 10 {
			delete(c.eventBuffer, k)
			deletedEntry++
		}
	}

	if enableDebugMessage && deletedEntry != 0 {
		fmt.Printf("- clear: %d deleted, %d remains\n", deletedEntry, len(c.eventBuffer))
	}
}

type eventEntry struct {
	Created time.Time
	DxIn    int
	Dxout   int
}
