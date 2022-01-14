package sync2

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

// What's in a name? Channels have all you need to emulate a counting
// semaphore with a boatload of extra functionality. However, in some
// cases, you just want a familiar API.

import (
	"context"
	"time"
)

// Semaphore is a counting semaphore with the option to
// specify a timeout.
type Semaphore struct {
	slots   chan struct{}
	timeout time.Duration
}

// NewSemaphore creates a Semaphore. The count parameter must be a positive
// number. A timeout of zero means that there is no timeout.
func NewSemaphore(count int, timeout time.Duration) *Semaphore {
	sem := &Semaphore{
		slots:   make(chan struct{}, count),
		timeout: timeout,
	}
	for i := 0; i < count; i++ {
		sem.slots <- struct{}{}
	}
	return sem
}

// Acquire returns true on successful acquisition, and
// false on a timeout.
func (sem *Semaphore) Acquire() bool {
	if sem.timeout == 0 {
		<-sem.slots
		return true
	}
	tm := time.NewTimer(sem.timeout)
	defer tm.Stop()
	select {
	case <-sem.slots:
		return true
	case <-tm.C:
		return false
	}
}

// AcquireContext returns true on successful acquisition, and
// false on context expiry. Timeout is ignored.
func (sem *Semaphore) AcquireContext(ctx context.Context) bool {
	select {
	case <-sem.slots:
		return true
	case <-ctx.Done():
		return false
	}
}

// TryAcquire acquires a semaphore if it's immediately available.
// It returns false otherwise.
func (sem *Semaphore) TryAcquire() bool {
	select {
	case <-sem.slots:
		return true
	default:
		return false
	}
}

// Release releases the acquired semaphore. You must
// not release more than the number of semaphores you've
// acquired.
func (sem *Semaphore) Release() {
	sem.slots <- struct{}{}
}

// Size returns the current number of available slots.
func (sem *Semaphore) Size() int {
	return len(sem.slots)
}
