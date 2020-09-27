//
//
//

package cache

import (
	"time"

	"github.com/ondi/go-cache"
)

type Mapped_t struct {
	Value interface{}
	ts    time.Time
}

type Value_t struct {
	Key   interface{}
	Value interface{}
}

type Evict func(interface{}) int

func Drop(interface{}) int { return 0 }

type Cache_t struct {
	c     *cache.Cache_t
	limit int
	ttl   time.Duration
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

func (self *Cache_t) flush(ts time.Time, it *cache.Value_t, keep int) bool {
	if self.c.Size() > keep || ts.After(it.Value().(Mapped_t).ts) {
		self.c.Remove(it.Key())
		self.evict(Value_t{Key: it.Key(), Value: it.Value().(Mapped_t).Value})
		return true
	}
	return false
}

func (self *Cache_t) Flush(ts time.Time) {
	for it := self.c.Front(); it != self.c.End() && self.flush(ts, it, self.limit); it = it.Next() {
	}
}

func (self *Cache_t) Create(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	it, ok = self.c.CreateBack(
		key,
		func() interface{} {
			return Mapped_t{Value: value(), ts: ts.Add(self.ttl)}
		},
	)
	res = it.Value().(Mapped_t).Value
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
			return Mapped_t{Value: v, ts: ts.Add(self.ttl)}, nil
		},
	)
	res = it.Value().(Mapped_t).Value
	self.Flush(ts)
	return
}

func (self *Cache_t) Write(ts time.Time, key interface{}, value_new func() interface{}, value_update func(interface{}) interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	it, ok = self.c.PushBack(
		key,
		func() interface{} {
			return Mapped_t{Value: value_new(), ts: ts.Add(self.ttl)}
		},
	)
	if !ok {
		it.Update(Mapped_t{Value: it.Value().(Mapped_t).Value, ts: ts.Add(self.ttl)})
	}
	res = it.Value().(Mapped_t).Value
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
			return Mapped_t{Value: v, ts: ts.Add(self.ttl)}, nil
		},
	)
	if !ok {
		if res, err = value_update(it.Value().(Mapped_t).Value); err != nil {
			return
		}
		it.Update(Mapped_t{Value: res, ts: ts.Add(self.ttl)})
	} else {
		res = it.Value().(Mapped_t).Value
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Update(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		res = value(it.Value().(Mapped_t).Value)
		it.Update(Mapped_t{Value: res, ts: ts.Add(self.ttl)})
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Update2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		if res, err = value(it.Value().(Mapped_t).Value); err != nil {
			return
		}
		it.Update(Mapped_t{Value: res, ts: ts.Add(self.ttl)})
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Refresh(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		res = value(it.Value().(Mapped_t).Value)
		it.Update(Mapped_t{Value: res, ts: it.Value().(Mapped_t).ts})
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Refresh2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	var it *cache.Value_t
	if it, ok = self.c.FindBack(key); ok {
		if res, err = value(it.Value().(Mapped_t).Value); err != nil {
			return
		}
		it.Update(Mapped_t{Value: res, ts: it.Value().(Mapped_t).ts})
	}
	self.Flush(ts)
	return
}

func (self *Cache_t) Get(ts time.Time, key interface{}) (interface{}, bool) {
	self.Flush(ts)
	if it, ok := self.c.FindBack(key); ok {
		it.Update(Mapped_t{Value: it.Value().(Mapped_t).Value, ts: ts.Add(self.ttl)})
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self *Cache_t) Find(ts time.Time, key interface{}) (interface{}, bool) {
	self.Flush(ts)
	if it, ok := self.c.Find(key); ok {
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self *Cache_t) Remove(ts time.Time, key interface{}) (interface{}, bool) {
	self.Flush(ts)
	if it, ok := self.c.Remove(key); ok {
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self *Cache_t) LeastDiff(ts time.Time) (time.Duration, bool) {
	self.Flush(ts)
	if self.c.Size() > 0 {
		return self.c.Front().Value().(Mapped_t).ts.Sub(ts), true
	}
	return 0, false
}

func (self *Cache_t) Range(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key(), it.Value().(Mapped_t).Value) == false {
			return
		}
	}
}

func (self *Cache_t) RangeTs(ts time.Time, f func(key interface{}, value interface{}, diff time.Duration) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key(), it.Value().(Mapped_t).Value, it.Value().(Mapped_t).ts.Sub(ts)) == false {
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
