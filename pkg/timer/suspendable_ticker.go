package timer

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"sync/atomic"
	"time"
)

// SuspendableTicker is similar to time.Ticker, but also offers Suspend() and Resume() functions.
// While the ticker is suspended, nothing comes from the time channel C
type SuspendableTicker struct {
	ticker *time.Ticker
	// C is user facing
	C chan time.Time

	suspended int64
}

// NewSuspendableTicker creates a new suspendable ticker, indicating whether the ticker should start
// suspendable or running
func NewSuspendableTicker(d time.Duration, initiallySuspended bool) *SuspendableTicker {
	s := &SuspendableTicker{
		ticker: time.NewTicker(d),
		C:      make(chan time.Time),
	}
	if initiallySuspended {
		s.suspended = 1
	}
	go s.loop()
	return s
}

// Suspend stops sending time events on the channel C
// time events sent during suspended time are lost
func (s *SuspendableTicker) Suspend() {
	atomic.StoreInt64(&s.suspended, 1)
}

// Resume re-enables time events on channel C
func (s *SuspendableTicker) Resume() {
	atomic.StoreInt64(&s.suspended, 0)
}

// Stop completely stops the timer, like time.Timer
func (s *SuspendableTicker) Stop() {
	s.ticker.Stop()
}

// TickNow generates a tick at this point in time. It may block
// if nothing consumes the tick.
func (s *SuspendableTicker) TickNow() {
	if atomic.LoadInt64(&s.suspended) == 0 {
		// not suspended
		s.C <- time.Now()
	}
}

func (s *SuspendableTicker) loop() {
	for t := range s.ticker.C {
		if atomic.LoadInt64(&s.suspended) == 0 {
			// not suspended
			s.C <- t
		}
	}
}
