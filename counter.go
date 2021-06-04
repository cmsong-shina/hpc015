package hpc015

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	simple = "2006-01-02 15:04:05"
)

// counter help counting
//
// counter provide information about number of in/out, occupants
//
// Due to sometime hpc015 send duplicated data,
// counter store recent data within 10 mins, and not count duplicated data.
//
// Use this just for a example,
// it store all count in one variable, not compitable with multiple device.
type counter struct {
	in          int
	out         int
	eventBuffer map[string]*eventEntry
	mux         *sync.Mutex
}

// Counter create new counter
//
// It run a go routine, which clear old events
func Counter() *counter {
	counter := &counter{
		eventBuffer: make(map[string]*eventEntry),
		mux:         &sync.Mutex{},
	}
	go counter.clearTicker()

	return counter
}

// Count a data
// If data is duplicated, return nil, Otherwise return eventEntry.
func (c *counter) Count(data *CacheData) *eventEntry {
	buf := make([]byte, 6, 6)
	buf[0] = data.Year
	buf[1] = data.Month
	buf[2] = data.Day
	buf[3] = data.Hour
	buf[4] = data.Minute
	buf[5] = data.Secound

	key := string(buf)
	ee := &eventEntry{
		time.Now(),
		time.Date(int(data.Year)+2000, time.Month(data.Month), int(data.Day), int(data.Hour), int(data.Minute), int(data.Secound), 0, time.Local),
		int(data.DxIn),
		int(data.DxOut),
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.eventBuffer[key]
	if ok {
		if EnableDebugMessage {
			fmt.Printf("- duplicated(%s)\n", ee.EventTime.Format(simple))
		}
		return nil
	}
	c.in += int(data.DxIn)
	c.out += int(data.DxOut)
	c.eventBuffer[key] = ee

	if EnableDebugMessage {
		fmt.Printf("- summary(%v): {in: %d, out: %d, current: %d}\n", ee.EventTime.Format(simple), ee.DxIn, ee.DxOut, c.GetOccupants())
	}
	return ee
}

// GetOccupants current count
func (c *counter) GetOccupants() int {
	return c.in - c.out
}

// GetInOut
func (c *counter) GetInOut() (int, int) {
	return c.in, c.out
}

// Set current count
func (c *counter) Set(num int) {
	c.in = num
	c.out = 0
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

	c.mux.Lock()
	defer c.mux.Unlock()

	for k, e := range c.eventBuffer {
		if math.Abs(e.Created.Sub(now).Minutes()) > 10 {
			delete(c.eventBuffer, k)
			deletedEntry++
		}
	}

	if EnableDebugMessage && deletedEntry != 0 {
		fmt.Printf("- clear: %d deleted, %d remains\n", deletedEntry, len(c.eventBuffer))
	}
}

type eventEntry struct {
	Created   time.Time
	EventTime time.Time
	DxIn      int
	DxOut     int
}
