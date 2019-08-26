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

type Evict interface {
	PushBackNoWait(interface{}) bool
}

type Cache_t struct {
	c * cache.Cache_t
	limit int
	ttl time.Duration
	evict Evict
}

type Evict_t []Value_t

func (self * Evict_t) PushBackNoWait(value interface{}) bool {
	*self = append(*self, value.(Value_t))
	return true
}

type Drop_t struct {}

func (Drop_t) PushBackNoWait(interface{}) bool {
	return true
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

func (self * Cache_t) flush(it * cache.Value_t, ts time.Time, keep int) bool {
	if self.c.Size() > keep || ts.Sub(it.Value().(Mapped_t).ts) > self.ttl {
		self.c.Remove(it.Key())
		self.evict.PushBackNoWait(Value_t{Key: it.Key(), Value: it.Value().(Mapped_t).Value})
		return true
	}
	return false
}

func (self * Cache_t) Flush(ts time.Time) {
	for it := self.c.Back(); it != self.c.End() && self.flush(it, ts, self.limit); it = it.Prev() {}
}

func (self * Cache_t) Create(ts time.Time, key interface{}, value func() interface{}) (interface{}, bool) {
	it, ok := self.c.CreateFront(key, func() interface{} {return Mapped_t{Value: value(), ts: ts}})
	if ok {
		self.Flush(ts)
	}
	return it.Value().(Mapped_t).Value, ok
}

func (self * Cache_t) Push(ts time.Time, key interface{}, value func() interface{}) (interface{}, bool) {
	it, ok := self.c.PushFront(key, func() interface{} {return Mapped_t{Value: value(), ts: ts}})
	if ok {
		self.Flush(ts)
	}
	return it.Value().(Mapped_t).Value, ok
}

func (self * Cache_t) Update(ts time.Time, key interface{}, value interface{}) (interface{}, bool) {
	it, ok := self.c.UpdateFront(key, Mapped_t{Value: value, ts: ts})
	if ok {
		self.Flush(ts)
	}
	return it.Value().(Mapped_t).Value, ok
}

func (self * Cache_t) Get(ts time.Time, key interface{}) (interface{}, bool) {
	self.Flush(ts)
	if it := self.c.FindFront(key); it != self.c.End() {
		it.Update(Mapped_t{Value: it.Value().(Mapped_t).Value, ts: ts})
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) Find(key interface{}) (interface{}, bool) {
	if it := self.c.Find(key); it != self.c.End() {
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) Remove(key interface{}) {
	self.c.Remove(key)
}

func (self * Cache_t) LeastTs() (time.Time, bool) {
	if self.c.Size() > 0 {
		return self.c.Back().Value().(Mapped_t).ts, true
	}
	return time.Time{}, false
}

func (self * Cache_t) Range(f func(key interface{}, value interface{}) bool) {
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key(), it.Value().(Mapped_t).Value) == false {
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
