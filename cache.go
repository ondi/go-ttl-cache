//
//
//

package cache

import (
	"time"

	"github.com/ondi/go-cache"
)

type Evict[Key_t comparable, Mapped_t any] func(ts time.Time, key Key_t, value Mapped_t)

func Drop[Key_t comparable, Mapped_t any](time.Time, Key_t, Mapped_t) {}

type Ttl_t[Mapped_t any] struct {
	ts    time.Time
	Value Mapped_t
}

type Cache_t[Key_t comparable, Mapped_t any] struct {
	c     *cache.Cache_t[Key_t, Ttl_t[Mapped_t]]
	ttl   time.Duration
	limit int
	evict Evict[Key_t, Mapped_t]
}

func New[Key_t comparable, Mapped_t any](limit int, ttl time.Duration, evict Evict[Key_t, Mapped_t]) (self *Cache_t[Key_t, Mapped_t]) {
	self = &Cache_t[Key_t, Mapped_t]{}
	self.c = cache.New[Key_t, Ttl_t[Mapped_t]]()
	if ttl < 0 {
		ttl = time.Duration(1<<63 - 1)
	}
	if limit < 0 {
		limit = 1<<63 - 1
	}
	self.ttl = ttl
	self.limit = limit
	self.evict = evict
	return
}

func (self *Cache_t[Key_t, Mapped_t]) flush(it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ts time.Time, keep int) bool {
	if self.c.Size() > keep || ts.After(it.Value.ts) {
		self.c.Remove(it.Key)
		self.evict(ts, it.Key, it.Value.Value)
		return true
	}
	return false
}

func (self *Cache_t[Key_t, Mapped_t]) Flush(ts time.Time) {
	for it := self.c.Front(); it != self.c.End() && self.flush(it, ts, self.limit); it = it.Next() {
	}
}

func (self *Cache_t[Key_t, Mapped_t]) FlushLimit(ts time.Time, limit int) {
	for it := self.c.Front(); it != self.c.End() && self.flush(it, ts, limit); it = it.Next() {
	}
}

func (self *Cache_t[Key_t, Mapped_t]) Create(ts time.Time, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.CreateBack(
		key,
		func(p *Ttl_t[Mapped_t]) {
			p.ts = ts.Add(self.ttl)
			value_init(&p.Value)
		},
		func(p *Ttl_t[Mapped_t]) {
			value_update(&p.Value)
		},
	)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Push(ts time.Time, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.PushBack(
		key,
		func(p *Ttl_t[Mapped_t]) {
			p.ts = ts.Add(self.ttl)
			value_init(&p.Value)
		},
		func(p *Ttl_t[Mapped_t]) {
			p.ts = ts.Add(self.ttl)
			value_update(&p.Value)
		},
	)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Update(ts time.Time, key Key_t, value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.FindBack(key)
	if ok {
		it.Value.ts = ts.Add(self.ttl)
		value_update(&it.Value.Value)
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Replace(ts time.Time, key Key_t, value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.Find(key)
	if ok {
		value_update(&it.Value.Value)
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Refresh(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.FindBack(key)
	if ok {
		it.Value.ts = ts.Add(self.ttl)
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Find(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.Find(key)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Remove(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, Ttl_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.c.Remove(key)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) LeastTs(ts time.Time) (time.Time, bool) {
	self.Flush(ts)
	if self.c.Size() > 0 {
		return self.c.Front().Value.ts, true
	}
	return time.Time{}, false
}

func (self *Cache_t[Key_t, Mapped_t]) Range(ts time.Time, f func(key Key_t, value Mapped_t) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key, it.Value.Value) == false {
			return
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) RangeTs(ts time.Time, f func(key Key_t, value Mapped_t, ts time.Time) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key, it.Value.Value, it.Value.ts) == false {
			return
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) Size(ts time.Time) int {
	self.Flush(ts)
	return self.c.Size()
}

func (self *Cache_t[Key_t, Mapped_t]) Limit() int {
	return self.limit
}

func (self *Cache_t[Key_t, Mapped_t]) TTL() time.Duration {
	return self.ttl
}
