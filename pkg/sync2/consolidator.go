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
	"sync"
	"sync/atomic"

	cache "github.com/bhojpur/cache/pkg/engine"
)

// Consolidator consolidates duplicate queries from executing simulaneously
// and shares results between them.
type Consolidator struct {
	*ConsolidatorCache

	mu      sync.Mutex
	queries map[string]*Result
}

// NewConsolidator creates a new Consolidator
func NewConsolidator() *Consolidator {
	return &Consolidator{
		queries:           make(map[string]*Result),
		ConsolidatorCache: NewConsolidatorCache(1000),
	}
}

// Result is a wrapper for result of a query.
type Result struct {
	// executing is used to block additional requests.
	// The original request holds a write lock while additional ones are blocked
	// on acquiring a read lock (see Wait() below.)
	executing    sync.RWMutex
	consolidator *Consolidator
	query        string
	Result       interface{}
	Err          error
}

// Create adds a query to currently executing queries and acquires a
// lock on its Result if it is not already present. If the query is
// a duplicate, Create returns false.
func (co *Consolidator) Create(query string) (r *Result, created bool) {
	co.mu.Lock()
	defer co.mu.Unlock()
	if r, ok := co.queries[query]; ok {
		return r, false
	}
	r = &Result{consolidator: co, query: query}
	r.executing.Lock()
	co.queries[query] = r
	return r, true
}

// Broadcast removes the entry from current queries and releases the
// lock on its Result. Broadcast should be invoked when original
// query completes execution.
func (rs *Result) Broadcast() {
	rs.consolidator.mu.Lock()
	defer rs.consolidator.mu.Unlock()
	delete(rs.consolidator.queries, rs.query)
	rs.executing.Unlock()
}

// Wait waits for the original query to complete execution. Wait should
// be invoked for duplicate queries.
func (rs *Result) Wait() {
	rs.consolidator.Record(rs.query)
	rs.executing.RLock()
}

// ConsolidatorCache is a thread-safe object used for counting how often recent
// queries have been consolidated.
// It is also used by the txserializer package to count how often transactions
// have been queued and had to wait because they targeted the same row (range).
type ConsolidatorCache struct {
	*cache.LRUCache
}

// NewConsolidatorCache creates a new cache with the given capacity.
func NewConsolidatorCache(capacity int64) *ConsolidatorCache {
	return &ConsolidatorCache{cache.NewLRUCache(capacity, func(_ interface{}) int64 {
		return 1
	})}
}

// Record increments the count for "query" by 1.
// If it's not in the cache yet, it will be added.
func (cc *ConsolidatorCache) Record(query string) {
	if v, ok := cc.Get(query); ok {
		v.(*ccount).add(1)
	} else {
		c := ccount(1)
		cc.Set(query, &c)
	}
}

// ConsolidatorCacheItem is a wrapper for the items in the consolidator cache
type ConsolidatorCacheItem struct {
	Query string
	Count int64
}

// Items returns the items in the cache as an array of String, int64 structs
func (cc *ConsolidatorCache) Items() []ConsolidatorCacheItem {
	items := cc.LRUCache.Items()
	ret := make([]ConsolidatorCacheItem, len(items))
	for i, v := range items {
		ret[i] = ConsolidatorCacheItem{Query: v.Key, Count: v.Value.(*ccount).get()}
	}
	return ret
}

// ccount elements are used with a cache.LRUCache object to track if another
// request for the same query is already in progress.
type ccount int64

func (cc *ccount) add(n int64) int64 {
	return atomic.AddInt64((*int64)(cc), n)
}

func (cc *ccount) get() int64 {
	return atomic.LoadInt64((*int64)(cc))
}
