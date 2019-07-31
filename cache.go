//
//
//

package cache

import "sync"
import "time"

import "github.com/ondi/go-cache"

type Cache_t struct {
	mx sync.Mutex
	c * cache.Cache_t
	limit int
	ttl time.Duration
}

type Mapped_t struct {
	Value interface{}
	Ts time.Time
}

type Value_t struct {
	Key interface{}
	Value interface{}
}

type Evict interface {
	Evict(Value_t) bool
}

type Evict_t []Value_t

func (self * Evict_t) Evict(value Value_t) bool {
	*self = append(*self, value)
	return true
}

type Drop_t struct {}

func (Drop_t) Evict(Value_t) bool {
	return true
}

func New(limit int, ttl time.Duration) (self * Cache_t) {
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
	return
}

func (self * Cache_t) __evict(it * cache.Value_t, ts time.Time, keep int, evicted Evict) bool {
	if self.c.Size() > keep || ts.Sub(it.Value().(Mapped_t).Ts) > self.ttl {
		self.c.Remove(it.Key())
		evicted.Evict(Value_t{Key: it.Key(), Value: it.Value().(Mapped_t).Value})
		return true
	}
	return false
}

func (self * Cache_t) __flush(ts time.Time, keep int, evicted Evict) {
	for it := self.c.Back(); it != self.c.End() && self.__evict(it, ts, keep, evicted); it = it.Prev() {}
}

func (self * Cache_t) Flush(ts time.Time, evicted Evict) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.__flush(ts, self.limit, evicted)
}

func (self * Cache_t) Create(ts time.Time, key interface{}, value interface{}, evicted Evict) (ok bool) {
	self.mx.Lock()
	defer self.mx.Unlock()
	if _, ok = self.c.CreateFront(key, Mapped_t{Value: value, Ts: ts}); ok {
		self.__flush(ts, self.limit, evicted)
	}
	return
}

func (self * Cache_t) Update(ts time.Time, key interface{}, value interface{}, evicted Evict) (ok bool) {
	self.mx.Lock()
	defer self.mx.Unlock()
	var it * cache.Value_t
	if it, ok = self.c.PushFront(key, Mapped_t{Value: value, Ts: ts}); ok {
		self.__flush(ts, self.limit, evicted)
	} else {
		it.Update(Mapped_t{Value: value, Ts: ts})
	}
	return
}

func (self * Cache_t) Get(ts time.Time, key interface{}, evicted Evict) (interface{}, bool) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.__flush(ts, self.limit, evicted)
	if it := self.c.FindFront(key); it != self.c.End() {
		it.Update(Mapped_t{Value: it.Value().(Mapped_t).Value, Ts: ts})
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) Find(ts time.Time, key interface{}, evicted Evict) (interface{}, bool) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.__flush(ts, self.limit, evicted)
	if it := self.c.Find(key); it != self.c.End() {
		return it.Value().(Mapped_t).Value, true
	}
	return nil, false
}

func (self * Cache_t) Remove(key interface{}) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.c.Remove(key)
}

func (self * Cache_t) LeastTs() (time.Time, bool) {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.c.Size() > 0 {
		return self.c.Back().Value().(Mapped_t).Ts, true
	}
	return time.Time{}, false
}

func (self * Cache_t) Range(f func(key interface{}, value interface{}) bool) {
	self.mx.Lock()
	defer self.mx.Unlock()
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key(), it.Value().(Mapped_t).Value) == false {
			return
		}
	}
}

func (self * Cache_t) Size() int {
	self.mx.Lock()
	defer self.mx.Unlock()
	return self.c.Size()
}

func (self * Cache_t) Limit() int {
	return self.limit
}


func (self * Cache_t) TTL() time.Duration {
	return self.ttl
}
