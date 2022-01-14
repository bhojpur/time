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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSemaNoTimeout(t *testing.T) {
	s := NewSemaphore(1, 0)
	s.Acquire()
	released := false
	go func() {
		released = true
		s.Release()
	}()
	s.Acquire()
	assert.True(t, released)
}

func TestSemaTimeout(t *testing.T) {
	s := NewSemaphore(1, 1*time.Millisecond)
	s.Acquire()
	release := make(chan struct{})
	released := make(chan struct{})
	go func() {
		<-release
		s.Release()
		released <- struct{}{}
	}()
	assert.False(t, s.Acquire())
	release <- struct{}{}
	<-released
	assert.True(t, s.Acquire())
}

func TestSemaAcquireContext(t *testing.T) {
	s := NewSemaphore(1, 0)
	s.Acquire()
	release := make(chan struct{})
	released := make(chan struct{})
	go func() {
		<-release
		s.Release()
		released <- struct{}{}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.False(t, s.AcquireContext(ctx))
	release <- struct{}{}
	<-released
	assert.True(t, s.AcquireContext(context.Background()))
}

func TestSemaTryAcquire(t *testing.T) {
	s := NewSemaphore(1, 0)
	assert.True(t, s.TryAcquire())
	assert.False(t, s.TryAcquire())
	s.Release()
	assert.True(t, s.TryAcquire())
}
