//
//
//

package cache

import (
	"math"
	"time"

	"github.com/ondi/go-cache"
)

type FixValue_t[Mapped_t any] struct {
	ts    time.Time
	Value Mapped_t
}

type Cache_t[Key_t comparable, Mapped_t any] struct {
	cx    *cache.Cache_t[Key_t, FixValue_t[Mapped_t]]
	ttl   time.Duration
	limit int
	evict Evict[Key_t, Mapped_t]
}

func New[Key_t comparable, Mapped_t any](limit int, ttl time.Duration, evict Evict[Key_t, Mapped_t]) (self *Cache_t[Key_t, Mapped_t]) {
	self = &Cache_t[Key_t, Mapped_t]{}
	self.cx = cache.New[Key_t, FixValue_t[Mapped_t]]()
	if ttl < 0 {
		ttl = time.Duration(math.MaxInt64)
	}
	if limit < 0 {
		limit = math.MaxInt
	}
	self.ttl = ttl
	self.limit = limit
	self.evict = evict
	return
}

func (self *Cache_t[Key_t, Mapped_t]) flush(ts time.Time, keep int, it *cache.Value_t[Key_t, FixValue_t[Mapped_t]]) bool {
	if self.cx.Size() > keep || ts.Before(it.Value.ts) == false {
		self.cx.Remove(it.Key)
		self.evict(it.Key, it.Value.Value)
		return true
	}
	return false
}

func (self *Cache_t[Key_t, Mapped_t]) Flush(ts time.Time) {
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if self.flush(ts, self.limit, it) == false {
			break
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) FlushLimit(ts time.Time, limit int) {
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if self.flush(ts, limit, it) == false {
			break
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) Create(ts time.Time, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.CreateBack(
		key,
		func(p *FixValue_t[Mapped_t]) {
			p.ts = ts.Add(self.ttl)
			value_init(&p.Value)
		},
		func(p *FixValue_t[Mapped_t]) {
			value_update(&p.Value)
		},
	)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Push(ts time.Time, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.PushBack(
		key,
		func(p *FixValue_t[Mapped_t]) {
			p.ts = ts.Add(self.ttl)
			value_init(&p.Value)
		},
		func(p *FixValue_t[Mapped_t]) {
			p.ts = ts.Add(self.ttl)
			value_update(&p.Value)
		},
	)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Update(ts time.Time, key Key_t, value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.FindBack(key)
	if ok {
		it.Value.ts = ts.Add(self.ttl)
		value_update(&it.Value.Value)
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Refresh(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.FindBack(key)
	if ok {
		it.Value.ts = ts.Add(self.ttl)
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Replace(ts time.Time, key Key_t, value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.Find(key)
	if ok {
		value_update(&it.Value.Value)
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Find(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.Find(key)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Remove(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, FixValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.Remove(key)
	return
}

func (self *Cache_t[Key_t, Mapped_t]) LeastTs(ts time.Time) (time.Time, bool) {
	self.Flush(ts)
	if self.cx.Size() > 0 {
		return self.cx.Front().Value.ts, true
	}
	return time.Time{}, false
}

func (self *Cache_t[Key_t, Mapped_t]) Range(ts time.Time, f func(key Key_t, value Mapped_t) bool) {
	self.Flush(ts)
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if f(it.Key, it.Value.Value) == false {
			return
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) RangeTs(ts time.Time, f func(key Key_t, value Mapped_t, ts time.Time) bool) {
	self.Flush(ts)
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if f(it.Key, it.Value.Value, it.Value.ts) == false {
			return
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) Size(ts time.Time) int {
	self.Flush(ts)
	return self.cx.Size()
}

func (self *Cache_t[Key_t, Mapped_t]) Limit() int {
	return self.limit
}

func (self *Cache_t[Key_t, Mapped_t]) TTL() time.Duration {
	return self.ttl
}
