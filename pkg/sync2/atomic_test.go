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

	"gotest.tools/assert"
)

func TestAtomicInt32(t *testing.T) {
	i := NewAtomicInt32(1)
	assert.Equal(t, int32(1), i.Get())

	i.Set(2)
	assert.Equal(t, int32(2), i.Get())

	i.Add(1)
	assert.Equal(t, int32(3), i.Get())

	i.CompareAndSwap(3, 4)
	assert.Equal(t, int32(4), i.Get())

	i.CompareAndSwap(3, 5)
	assert.Equal(t, int32(4), i.Get())
}

func TestAtomicInt64(t *testing.T) {
	i := NewAtomicInt64(1)
	assert.Equal(t, int64(1), i.Get())

	i.Set(2)
	assert.Equal(t, int64(2), i.Get())

	i.Add(1)
	assert.Equal(t, int64(3), i.Get())

	i.CompareAndSwap(3, 4)
	assert.Equal(t, int64(4), i.Get())

	i.CompareAndSwap(3, 5)
	assert.Equal(t, int64(4), i.Get())
}

func TestAtomicFloat64(t *testing.T) {
	i := NewAtomicFloat64(1.0)
	assert.Equal(t, float64(1.0), i.Get())

	i.Set(2.0)
	assert.Equal(t, float64(2.0), i.Get())
	{
		swapped := i.CompareAndSwap(2.0, 4.0)
		assert.Equal(t, float64(4), i.Get())
		assert.Equal(t, true, swapped)
	}
	{
		swapped := i.CompareAndSwap(2.0, 5.0)
		assert.Equal(t, float64(4), i.Get())
		assert.Equal(t, false, swapped)
	}
}

func TestAtomicDuration(t *testing.T) {
	d := NewAtomicDuration(time.Second)
	assert.Equal(t, time.Second, d.Get())

	d.Set(time.Second * 2)
	assert.Equal(t, time.Second*2, d.Get())

	d.Add(time.Second)
	assert.Equal(t, time.Second*3, d.Get())

	d.CompareAndSwap(time.Second*3, time.Second*4)
	assert.Equal(t, time.Second*4, d.Get())

	d.CompareAndSwap(time.Second*3, time.Second*5)
	assert.Equal(t, time.Second*4, d.Get())
}

func TestAtomicString(t *testing.T) {
	var s AtomicString
	assert.Equal(t, "", s.Get())

	s.Set("a")
	assert.Equal(t, "a", s.Get())

	assert.Equal(t, false, s.CompareAndSwap("b", "c"))
	assert.Equal(t, "a", s.Get())

	assert.Equal(t, true, s.CompareAndSwap("a", "c"))
	assert.Equal(t, "c", s.Get())
}

func TestAtomicBool(t *testing.T) {
	b := NewAtomicBool(true)
	assert.Equal(t, true, b.Get())

	b.Set(false)
	assert.Equal(t, false, b.Get())

	b.Set(true)
	assert.Equal(t, true, b.Get())

	assert.Equal(t, false, b.CompareAndSwap(false, true))

	assert.Equal(t, true, b.CompareAndSwap(true, false))

	assert.Equal(t, true, b.CompareAndSwap(false, false))

	assert.Equal(t, true, b.CompareAndSwap(false, true))

	assert.Equal(t, true, b.CompareAndSwap(true, true))
}
