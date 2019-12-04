//
//
//

package cache

import "time"

import "github.com/ondi/go-cache"

type Mapped_t struct {
	Value interface{}
	ts time.Time
}

type Value_t struct {
	Key interface{}
	Value interface{}
}

type Evict func(interface{}) int

func Drop(interface{}) int {return 0}

type Cache_t struct {
	c * cache.Cache_t
	limit int
	ttl time.Duration
	evict Evict
}

func New(limit int, ttl time.Duration, evict Evict) (self * Cache_t) {
	self = &Cache_t{}
	self.c = cache.New()
	if ttl <= 0 {
		ttl = time.Duration(1 << 63 - 1)
	}
	if limit <= 0 {
		limit = 1 << 63 - 1
	}
	self.ttl = ttl
	self.limit = limit
	self.evict = evict
	return
}

func (self * Cache_t) flush(ts time.Time, it * cache.Value_t, keep int) bool {
	if self.c.Size() > keep || ts.After(it.Value().(Mapped_t).ts) {
		self.c.Remove(it.Key())
		self.evict(Value_t{Key: it.Key(), Value: it.Value().(Mapped_t).Value})
		return true
	}
	return false
}

func (self * Cache_t) Flush(ts time.Time) {
	for it := self.c.Front(); it != self.c.End() && self.flush(ts, it, self.limit); it = it.Next() {}
}

func (self * Cache_t) Create(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	var it * cache.Value_t
	it, ok = self.c.CreateBack(key, func() interface{} {return Mapped_t{Value: value(), ts: ts.Add(self.ttl)}})
	res = it.Value().(Mapped_t).Value
	self.Flush(ts)
	return
}

func (self * Cache_t) Push(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	var it * cache.Value_t
	if it, ok = self.c.PushBack(key, func() interface{} {return Mapped_t{Value: value(), ts: ts.Add(self.ttl)}}); !ok {
		it.Update(Mapped_t{Value: it.Value().(Mapped_t).Value, ts: ts.Add(self.ttl)})
	}
	res = it.Value().(Mapped_t).Value
	self.Flush(ts)
	return
}

func (self * Cache_t) Update(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	var it * cache.Value_t
	if it, ok = self.c.UpdateBack(key, func(prev interface{}) interface{} {return Mapped_t{Value: value(prev.(Mapped_t).Value), ts: ts.Add(self.ttl)}}); ok {
		res = it.Value().(Mapped_t).Value
	}
	self.Flush(ts)
	return
}

func (self * Cache_t) Refresh(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	var it * cache.Value_t
	if it, ok = self.c.Update(key, func(prev interface{}) interface{} {return Mapped_t{Value: value(prev.(Mapped_t).Value), ts: prev.(Mapped_t).ts}}); ok {
		res = it.Value().(Mapped_t).Value
	}
	self.Flush(ts)
	return
}

func (self * Cache_t) Get(ts time.Time, key interface{}) (interface{}, bool) {
	self.Flush(ts)
	if it, ok := self.c.FindBack(key); ok {
		it.Update(Mapped_t{Value: it.Value().(Mapped_t).Value, ts: ts.Add(self.ttl)})
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) Find(ts time.Time, key interface{}) (interface{}, bool) {
	self.Flush(ts)
	if it, ok := self.c.Find(key); ok {
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) Remove(ts time.Time, key interface{}) (interface {}, bool) {
	self.Flush(ts)
	if it, ok := self.c.Remove(key); ok {
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) LeastDiff(ts time.Time) (time.Duration, bool) {
	self.Flush(ts)
	if self.c.Size() > 0 {
		return self.c.Front().Value().(Mapped_t).ts.Sub(ts), true
	}
	return 0, false
}

func (self * Cache_t) Range(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key(), it.Value().(Mapped_t).Value) == false {
			return
		}
	}
}

func (self * Cache_t) RangeTs(ts time.Time, f func(key interface{}, value interface{}, ts time.Time) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key(), it.Value().(Mapped_t).Value, it.Value().(Mapped_t).ts) == false {
			return
		}
	}
}

func (self * Cache_t) Size() int {
	return self.c.Size()
}

func (self * Cache_t) Limit() int {
	return self.limit
}

func (self * Cache_t) TTL() time.Duration {
	return self.ttl
}
