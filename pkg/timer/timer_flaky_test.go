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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/time/pkg/sync2"
)

const (
	half    = 50 * time.Millisecond
	quarter = 25 * time.Millisecond
	tenth   = 10 * time.Millisecond
)

var numcalls sync2.AtomicInt64

func f() {
	numcalls.Add(1)
}

func TestWait(t *testing.T) {
	numcalls.Set(0)
	timer := NewTimer(quarter)
	assert.False(t, timer.Running())
	timer.Start(f)
	defer timer.Stop()
	assert.True(t, timer.Running())
	time.Sleep(tenth)
	assert.Equal(t, int64(0), numcalls.Get())
	time.Sleep(quarter)
	assert.Equal(t, int64(1), numcalls.Get())
	time.Sleep(quarter)
	assert.Equal(t, int64(2), numcalls.Get())
}

func TestReset(t *testing.T) {
	numcalls.Set(0)
	timer := NewTimer(half)
	timer.Start(f)
	defer timer.Stop()
	timer.SetInterval(quarter)
	time.Sleep(tenth)
	assert.Equal(t, int64(0), numcalls.Get())
	time.Sleep(quarter)
	assert.Equal(t, int64(1), numcalls.Get())
}

func TestIndefinite(t *testing.T) {
	numcalls.Set(0)
	timer := NewTimer(0)
	timer.Start(f)
	defer timer.Stop()
	timer.TriggerAfter(quarter)
	time.Sleep(tenth)
	assert.Equal(t, int64(0), numcalls.Get())
	time.Sleep(quarter)
	assert.Equal(t, int64(1), numcalls.Get())
}
