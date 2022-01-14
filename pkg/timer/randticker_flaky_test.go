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
	"testing"
	"time"
)

const (
	testDuration = 100 * time.Millisecond
	testVariance = 20 * time.Millisecond
)

func TestTick(t *testing.T) {
	tkr := NewRandTicker(testDuration, testVariance)
	for i := 0; i < 5; i++ {
		start := time.Now()
		end := <-tkr.C
		diff := start.Add(testDuration).Sub(end)
		tolerance := testVariance + 20*time.Millisecond
		if diff < -tolerance || diff > tolerance {
			t.Errorf("start: %v, end: %v, diff %v. Want <%v tolerenace", start, end, diff, tolerance)
		}
	}
	tkr.Stop()
	_, ok := <-tkr.C
	if ok {
		t.Error("Channel was not closed")
	}
}

func TestTickSkip(t *testing.T) {
	tkr := NewRandTicker(10*time.Millisecond, 1*time.Millisecond)
	time.Sleep(35 * time.Millisecond)
	end := <-tkr.C
	diff := time.Since(end)
	if diff < 20*time.Millisecond {
		t.Errorf("diff: %v, want >20ms", diff)
	}

	// This tick should be up-to-date
	end = <-tkr.C
	diff = time.Since(end)
	if diff > 1*time.Millisecond {
		t.Errorf("diff: %v, want <1ms", diff)
	}
	tkr.Stop()
}
