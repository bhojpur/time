package timeseries

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
	"math"
	"testing"
	"time"
)

func isNear(x *Float, y float64, tolerance float64) bool {
	return math.Abs(x.Value()-y) < tolerance
}

func isApproximate(x *Float, y float64) bool {
	return isNear(x, y, 1e-2)
}

func checkApproximate(t *testing.T, o Observable, y float64) {
	x := o.(*Float)
	if !isApproximate(x, y) {
		t.Errorf("Wanted %g, got %g", y, x.Value())
	}
}

func checkNear(t *testing.T, o Observable, y, tolerance float64) {
	x := o.(*Float)
	if !isNear(x, y, tolerance) {
		t.Errorf("Wanted %g +- %g, got %g", y, tolerance, x.Value())
	}
}

var baseTime = time.Date(2018, 3, 26, 0, 0, 0, 0, time.UTC)

func tu(s int64) time.Time {
	return baseTime.Add(time.Duration(s) * time.Second)
}

func tu2(s int64, ns int64) time.Time {
	return baseTime.Add(time.Duration(s)*time.Second + time.Duration(ns)*time.Nanosecond)
}

func TestBasicTimeSeries(t *testing.T) {
	ts := NewTimeSeries(NewFloat)
	fo := new(Float)
	*fo = Float(10)
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	checkApproximate(t, ts.Range(tu(0), tu(1)), 40)
	checkApproximate(t, ts.Total(), 40)
	ts.AddWithTime(fo, tu(3))
	ts.AddWithTime(fo, tu(3))
	ts.AddWithTime(fo, tu(3))
	checkApproximate(t, ts.Range(tu(0), tu(2)), 40)
	checkApproximate(t, ts.Range(tu(2), tu(4)), 30)
	checkApproximate(t, ts.Total(), 70)
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	checkApproximate(t, ts.Range(tu(0), tu(2)), 60)
	checkApproximate(t, ts.Range(tu(2), tu(4)), 30)
	checkApproximate(t, ts.Total(), 90)
	*fo = Float(100)
	ts.AddWithTime(fo, tu(100))
	checkApproximate(t, ts.Range(tu(99), tu(100)), 100)
	checkApproximate(t, ts.Range(tu(0), tu(4)), 36)
	checkApproximate(t, ts.Total(), 190)
	*fo = Float(10)
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	checkApproximate(t, ts.Range(tu(0), tu(4)), 44)
	checkApproximate(t, ts.Range(tu(37), tu2(100, 100e6)), 100)
	checkApproximate(t, ts.Range(tu(50), tu2(100, 100e6)), 100)
	checkApproximate(t, ts.Range(tu(99), tu2(100, 100e6)), 100)
	checkApproximate(t, ts.Total(), 210)

	for i, l := range ts.ComputeRange(tu(36), tu(100), 64) {
		if i == 63 {
			checkApproximate(t, l, 100)
		} else {
			checkApproximate(t, l, 0)
		}
	}

	checkApproximate(t, ts.Range(tu(0), tu(100)), 210)
	checkApproximate(t, ts.Range(tu(10), tu(100)), 100)

	for i, l := range ts.ComputeRange(tu(0), tu(100), 100) {
		if i < 10 {
			checkApproximate(t, l, 11)
		} else if i >= 90 {
			checkApproximate(t, l, 10)
		} else {
			checkApproximate(t, l, 0)
		}
	}
}

func TestFloat(t *testing.T) {
	f := Float(1)
	if g, w := f.String(), "1"; g != w {
		t.Errorf("Float(1).String = %q; want %q", g, w)
	}
	f2 := Float(2)
	var o Observable = &f2
	f.Add(o)
	if g, w := f.Value(), 3.0; g != w {
		t.Errorf("Float post-add = %v; want %v", g, w)
	}
	f.Multiply(2)
	if g, w := f.Value(), 6.0; g != w {
		t.Errorf("Float post-multiply = %v; want %v", g, w)
	}
	f.Clear()
	if g, w := f.Value(), 0.0; g != w {
		t.Errorf("Float post-clear = %v; want %v", g, w)
	}
	f.CopyFrom(&f2)
	if g, w := f.Value(), 2.0; g != w {
		t.Errorf("Float post-CopyFrom = %v; want %v", g, w)
	}
}

type mockClock struct {
	time time.Time
}

func (m *mockClock) Time() time.Time { return m.time }
func (m *mockClock) Set(t time.Time) { m.time = t }

const buckets = 6

var testResolutions = []time.Duration{
	10 * time.Second,  // level holds one minute of observations
	100 * time.Second, // level holds ten minutes of observations
	10 * time.Minute,  // level holds one hour of observations
}

// TestTimeSeries uses a small number of buckets to force a higher
// error rate on approximations from the timeseries.
type TestTimeSeries struct {
	timeSeries
}

func TestExpectedErrorRate(t *testing.T) {
	ts := new(TestTimeSeries)
	fake := new(mockClock)
	fake.Set(time.Now())
	ts.timeSeries.init(testResolutions, NewFloat, buckets, fake)
	for i := 1; i <= 61*61; i++ {
		fake.Set(fake.Time().Add(1 * time.Second))
		ob := Float(1)
		ts.AddWithTime(&ob, fake.Time())

		// The results should be accurate within one missing bucket (1/6) of the observations recorded.
		checkNear(t, ts.Latest(0, buckets), min(float64(i), 60), 10)
		checkNear(t, ts.Latest(1, buckets), min(float64(i), 600), 100)
		checkNear(t, ts.Latest(2, buckets), min(float64(i), 3600), 600)
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
