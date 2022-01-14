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

import (
	"testing"
	"time"
)

// makeAfterFnWithLatch returns a fake alternative to time.After that blocks until
// the release function is called. The fake doesn't support having multiple concurrent
// calls to the After function, which is ok because Batcher should never do that.
func makeAfterFnWithLatch(t *testing.T) (func(time.Duration) <-chan time.Time, func()) {
	latch := make(chan time.Time, 1)
	afterFn := func(d time.Duration) <-chan time.Time {
		return latch
	}

	releaseFn := func() {
		select {
		case latch <- time.Now():
		default:
			t.Errorf("Previous batch still hasn't been released")
		}
	}
	return afterFn, releaseFn
}

func TestBatcher(t *testing.T) {
	interval := time.Duration(50 * time.Millisecond)

	afterFn, releaseBatch := makeAfterFnWithLatch(t)
	b := newBatcherForTest(interval, afterFn)

	waitersFinished := NewAtomicInt32(0)

	startWaiter := func(testcase string, want int) {
		go func() {
			id := b.Wait()
			if id != want {
				t.Errorf("%s: got %d, want %d", testcase, id, want)
			}
			waitersFinished.Add(1)
		}()
	}

	awaitVal := func(name string, val *AtomicInt32, expected int32) {
		for count := 0; val.Get() != expected; count++ {
			time.Sleep(50 * time.Millisecond)
			if count > 5 {
				t.Errorf("Timed out waiting for %s to be %v", name, expected)
				return
			}
		}
	}

	awaitBatch := func(name string, n int32) {
		// Wait for all the waiters to register
		awaitVal("Batcher.waiters for "+name, &b.waiters, n)
		// Release the batch and wait for the batcher to catch up.
		if waitersFinished.Get() != 0 {
			t.Errorf("Waiters finished before being released")
		}
		releaseBatch()
		awaitVal("Batcher.waiters for "+name, &b.waiters, 0)
		// Make sure the waiters actually run so they can verify their batch number.
		awaitVal("waitersFinshed for "+name, &waitersFinished, n)
		waitersFinished.Set(0)
	}

	// test single waiter
	startWaiter("single waiter", 1)
	awaitBatch("single waiter", 1)

	// multiple waiters all at once
	startWaiter("concurrent waiter", 2)
	startWaiter("concurrent waiter", 2)
	startWaiter("concurrent waiter", 2)
	awaitBatch("concurrent waiter", 3)

	startWaiter("more waiters", 3)
	startWaiter("more waiters", 3)
	startWaiter("more waiters", 3)
	startWaiter("more waiters", 3)
	startWaiter("more waiters", 3)
	awaitBatch("more waiters", 5)
}
