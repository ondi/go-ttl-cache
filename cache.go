//
//
//

package cache

import (
	"time"

	"github.com/ondi/go-cache"
)

type Evict func(key interface{}, value interface{})

func Drop(interface{}, interface{}) {}

type mapped_t struct {
	value interface{}
	ts    time.Time
}

type Cache_t struct {
	c     *cache.Cache_t
	ttl   time.Duration
	limit int
	evict Evict
}

func New(limit int, ttl time.Duration, evict Evict) (self *Cache_t) {
	self = &Cache_t{}
	self.c = cache.New()
	if ttl <= 0 {
		ttl = time.Duration(1<<63 - 1)
	}
	if limit <= 0 {
		limit = 1<<63 - 1
	}
	self.ttl = ttl
	self.limit = limit
	self.evict = evict
	return
}

func (self *Cache_t) flush(it *cache.Value_t, ts time.Time, keep int) bool {
	if self.c.Size() > keep || ts.After(it.Value.(mapped_t).ts) {
		self.c.Remove(it.Key)
		self.evict(it.Key, it.Value.(mapped_t).value)
		return true
	}
	return false
}

func (self *Cache_t) Flush(ts time.Time) {
	for it := self.c.Front(); it != self.c.End() && self.flush(it, ts, self.limit); it = it.Next() {
	}
}

func (self *Cache_t) FlushLimit(ts time.Time, limit int) {
	for it := self.c.Front(); it != self.c.End() && self.flush(it, ts, limit); it = it.Next() {
	}
}

func (self *Cache_t) Create(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	it, ok = self.c.CreateBack(
		key,
		func() interface{} {
			return mapped_t{value: value(), ts: ts.Add(self.ttl)}
		},
	)
	res = it.Value.(mapped_t).value
	self.Flush(ts)
	return
}

func (self *Cache_t) Create2(ts time.Time, key interface{}, value func() (interface{}, error)) (res interface{}, ok bool, err error) {
	var it *cache.Value_t
	it, ok, err = self.c.CreateBack2(
		key,
		func() (v interface{}, err error) {
			if v, err = value(); err != nil {
				return
			}
			return mapped_t{value: v, ts: ts.Add(self.ttl)}, nil
		},
	)
	res = it.Value.(mapped_t).value
	self.Flush(ts)
	return
}

func (self *Cache_t) Write(ts time.Time, key interface{}, value_new func() interface{}, value_update func(interface{}) interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	it, ok = self.c.PushBack(
		key,
		func() interface{} {
			return mapped_t{value: value_new(), ts: ts.Add(self.ttl)}
		},
	)
	if !ok {
		it.Value = mapped_t{value: value_update(it.Value.(mapped_t).value), ts: ts.Add(self.ttl)}
	}
	res = it.Value.(mapped_t).value
	self.Flush(ts)
	return
}

func (self *Cache_t) Write2(ts time.Time, key interface{}, value_new func() (interface{}, error), value_update func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	var it *cache.Value_t
	it, ok, err = self.c.PushBack2(
		key,
		func() (v interface{}, err error) {
			if v, err = value_new(); err != nil {
				return
			}
			return mapped_t{value: v, ts: ts.Add(self.ttl)}, nil
		},
	)
	if !ok {
		if res, err = value_update(it.Value.(mapped_t).value); err != nil {
			return
		}
		it.Value = mapped_t{value: res, ts: ts.Add(self.ttl)}
	} else {
		res = it.Value.(mapped_t).value
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Update(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		res = value(it.Value.(mapped_t).value)
		it.Value = mapped_t{value: res, ts: ts.Add(self.ttl)}
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Update2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		if res, err = value(it.Value.(mapped_t).value); err != nil {
			return
		}
		it.Value = mapped_t{value: res, ts: ts.Add(self.ttl)}
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Replace(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	if it, ok = self.c.Find(key); ok {
		res = value(it.Value.(mapped_t).value)
		it.Value = mapped_t{value: res, ts: it.Value.(mapped_t).ts}
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Replace2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	var it *cache.Value_t
	if it, ok = self.c.Find(key); ok {
		if res, err = value(it.Value.(mapped_t).value); err != nil {
			return
		}
		it.Value = mapped_t{value: res, ts: it.Value.(mapped_t).ts}
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Get(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.Flush(ts)
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		res = it.Value.(mapped_t).value
		it.Value = mapped_t{value: res, ts: ts.Add(self.ttl)}
		return
	}
	return nil, false
}

func (self *Cache_t) Find(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.Flush(ts)
	var it *cache.Value_t
	if it, ok = self.c.Find(key); ok {
		res = it.Value.(mapped_t).value
		return
	}
	return
}

func (self *Cache_t) Remove(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.Flush(ts)
	var it *cache.Value_t
	if it, ok = self.c.Remove(key); ok {
		res = it.Value.(mapped_t).value
		return
	}
	return
}

func (self *Cache_t) LeastDiff(ts time.Time) (time.Duration, bool) {
	self.Flush(ts)
	if self.c.Size() > 0 {
		return self.c.Front().Value.(mapped_t).ts.Sub(ts), true
	}
	return self.ttl, false
}

func (self *Cache_t) Range(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key, it.Value.(mapped_t).value) == false {
			return
		}
	}
}

func (self *Cache_t) RangeTs(ts time.Time, f func(key interface{}, value interface{}, diff time.Duration) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key, it.Value.(mapped_t).value, it.Value.(mapped_t).ts.Sub(ts)) == false {
			return
		}
	}
}

func (self *Cache_t) Size() int {
	return self.c.Size()
}

func (self *Cache_t) Limit() int {
	return self.limit
}

func (self *Cache_t) TTL() time.Duration {
	return self.ttl
}
