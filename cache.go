//
//
//

package cache

import (
	"time"

	"github.com/ondi/go-cache"
)

type Evict[Key_t comparable, Mapped_t any] func(key Key_t, value Mapped_t)

func Drop[Key_t comparable, Mapped_t any](Key_t, Mapped_t) {}

type mapped_t[Mapped_t any] struct {
	value Mapped_t
	ts    time.Time
}

type Cache_t[Key_t comparable, Mapped_t any] struct {
	c     *cache.Cache_t[Key_t, mapped_t[Mapped_t]]
	ttl   time.Duration
	limit int
	evict Evict[Key_t, Mapped_t]
}

func New[Key_t comparable, Mapped_t any](limit int, ttl time.Duration, evict Evict[Key_t, Mapped_t]) (self *Cache_t[Key_t, Mapped_t]) {
	self = &Cache_t[Key_t, Mapped_t]{}
	self.c = cache.New[Key_t, mapped_t[Mapped_t]]()
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

func (self *Cache_t[Key_t, Mapped_t]) flush(it *cache.Value_t[Key_t, mapped_t[Mapped_t]], ts time.Time, keep int) bool {
	if self.c.Size() > keep || ts.After(it.Value.ts) {
		self.c.Remove(it.Key)
		self.evict(it.Key, it.Value.value)
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

func (self *Cache_t[Key_t, Mapped_t]) Create(ts time.Time, key Key_t, value_new func() Mapped_t, value_update func(Mapped_t) Mapped_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.CreateBack(
		key,
		func() mapped_t[Mapped_t] {
			return mapped_t[Mapped_t]{value: value_new(), ts: ts.Add(self.ttl)}
		},
	)
	if !ok {
		it.Value = mapped_t[Mapped_t]{value: value_update(it.Value.value), ts: it.Value.ts}
	}
	res = it.Value.value
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Create2(ts time.Time, key Key_t, value_new func() (Mapped_t, error), value_update func(Mapped_t) (Mapped_t, error)) (res Mapped_t, ok bool, err error) {
	self.Flush(ts)
	it, ok := self.c.CreateBack2(
		key,
		func() *cache.Value_t[Key_t, mapped_t[Mapped_t]] {
			if res, err = value_new(); err != nil {
				return nil
			}
			return &cache.Value_t[Key_t, mapped_t[Mapped_t]]{Value: mapped_t[Mapped_t]{value: res, ts: ts.Add(self.ttl)}}
		},
	)
	if !ok {
		if res, err = value_update(it.Value.value); err != nil {
			return
		}
		it.Value = mapped_t[Mapped_t]{value: res, ts: it.Value.ts}
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Push(ts time.Time, key Key_t, value_new func() Mapped_t, value_update func(Mapped_t) Mapped_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.PushBack(
		key,
		func() mapped_t[Mapped_t] {
			return mapped_t[Mapped_t]{value: value_new(), ts: ts.Add(self.ttl)}
		},
	)
	if !ok {
		it.Value = mapped_t[Mapped_t]{value: value_update(it.Value.value), ts: ts.Add(self.ttl)}
	}
	res = it.Value.value
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Push2(ts time.Time, key Key_t, value_new func() (Mapped_t, error), value_update func(Mapped_t) (Mapped_t, error)) (res Mapped_t, ok bool, err error) {
	self.Flush(ts)
	it, ok := self.c.PushBack2(
		key,
		func() *cache.Value_t[Key_t, mapped_t[Mapped_t]] {
			if res, err = value_new(); err != nil {
				return nil
			}
			return &cache.Value_t[Key_t, mapped_t[Mapped_t]]{Value: mapped_t[Mapped_t]{value: res, ts: ts.Add(self.ttl)}}
		},
	)
	if !ok {
		if res, err = value_update(it.Value.value); err != nil {
			return
		}
		it.Value = mapped_t[Mapped_t]{value: res, ts: ts.Add(self.ttl)}
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Update(ts time.Time, key Key_t, value func(Mapped_t) Mapped_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.FindBack(key)
	if ok {
		res = value(it.Value.value)
		it.Value = mapped_t[Mapped_t]{value: res, ts: ts.Add(self.ttl)}
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Update2(ts time.Time, key Key_t, value func(Mapped_t) (Mapped_t, error)) (res Mapped_t, ok bool, err error) {
	self.Flush(ts)
	it, ok := self.c.FindBack(key)
	if ok {
		if res, err = value(it.Value.value); err != nil {
			return
		}
		it.Value = mapped_t[Mapped_t]{value: res, ts: ts.Add(self.ttl)}
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Replace(ts time.Time, key Key_t, value func(Mapped_t) Mapped_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.Find(key)
	if ok {
		res = value(it.Value.value)
		it.Value = mapped_t[Mapped_t]{value: res, ts: it.Value.ts}
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Replace2(ts time.Time, key Key_t, value func(Mapped_t) (Mapped_t, error)) (res Mapped_t, ok bool, err error) {
	self.Flush(ts)
	it, ok := self.c.Find(key)
	if ok {
		if res, err = value(it.Value.value); err != nil {
			return
		}
		it.Value = mapped_t[Mapped_t]{value: res, ts: it.Value.ts}
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Get(ts time.Time, key Key_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.FindBack(key);
	if ok {
		res = it.Value.value
		it.Value = mapped_t[Mapped_t]{value: res, ts: ts.Add(self.ttl)}
		return
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Find(ts time.Time, key Key_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.Find(key);
	if ok {
		res = it.Value.value
		return
	}
	return
}

func (self *Cache_t[Key_t, Mapped_t]) Remove(ts time.Time, key Key_t) (res Mapped_t, ok bool) {
	self.Flush(ts)
	it, ok := self.c.Remove(key)
	if ok {
		res = it.Value.value
		return
	}
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
		if f(it.Key, it.Value.value) == false {
			return
		}
	}
}

func (self *Cache_t[Key_t, Mapped_t]) RangeTs(ts time.Time, f func(key Key_t, value Mapped_t, ts time.Time) bool) {
	self.Flush(ts)
	for it := self.c.Front(); it != self.c.End(); it = it.Next() {
		if f(it.Key, it.Value.value, it.Value.ts) == false {
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
