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
	"math/rand"
	"time"
)

// RandTicker is just like time.Ticker, except that
// it adds randomness to the events.
type RandTicker struct {
	C    <-chan time.Time
	done chan struct{}
}

// NewRandTicker creates a new RandTicker. d is the duration,
// and variance specifies the variance. The ticker will tick
// every d +/- variance.
func NewRandTicker(d, variance time.Duration) *RandTicker {
	c := make(chan time.Time, 1)
	done := make(chan struct{})
	go func() {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			vr := time.Duration(rnd.Int63n(int64(2*variance)) - int64(variance))
			tmr := time.NewTimer(d + vr)
			select {
			case <-tmr.C:
				select {
				case c <- time.Now():
				default:
				}
			case <-done:
				tmr.Stop()
				close(c)
				return
			}
		}
	}()
	return &RandTicker{
		C:    c,
		done: done,
	}
}

// Stop stops the ticker and closes the underlying channel.
func (tkr *RandTicker) Stop() {
	close(tkr.done)
}
