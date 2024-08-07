//
//
//

package cache

import (
	"math"
	"time"

	"github.com/ondi/go-cache"
)

type VarValue_t[Mapped_t any] struct {
	ts    time.Time
	Value Mapped_t
}

type CacheVar_t[Key_t comparable, Mapped_t any] struct {
	cx    *cache.Cache_t[Key_t, VarValue_t[Mapped_t]]
	limit int
	evict Evict[Key_t, Mapped_t]
}

func NewVar[Key_t comparable, Mapped_t any](limit int, evict Evict[Key_t, Mapped_t]) (self *CacheVar_t[Key_t, Mapped_t]) {
	self = &CacheVar_t[Key_t, Mapped_t]{}
	self.cx = cache.New[Key_t, VarValue_t[Mapped_t]]()
	if limit < 0 {
		limit = math.MaxInt
	}
	self.limit = limit
	self.evict = evict
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) flush(ts time.Time, keep int, it *cache.Value_t[Key_t, VarValue_t[Mapped_t]]) bool {
	if self.cx.Size() > keep || ts.Before(it.Value.ts) == false {
		self.cx.Remove(it.Key)
		self.evict(it.Key, it.Value.Value)
		return true
	}
	return false
}

func (self *CacheVar_t[Key_t, Mapped_t]) Flush(ts time.Time) {
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if self.flush(ts, self.limit, it) == false {
			break
		}
	}
}

func (self *CacheVar_t[Key_t, Mapped_t]) FlushLimit(ts time.Time, limit int) {
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if self.flush(ts, limit, it) == false {
			break
		}
	}
}

func (self *CacheVar_t[Key_t, Mapped_t]) find_place() {
	it1 := self.cx.Back()
	it2 := self.cx.End()
	for it3 := self.cx.Back().Prev(); it3 != self.cx.End() && it1.Value.ts.Before(it3.Value.ts); it3 = it3.Prev() {
		it2 = it3
	}
	if it2 != self.cx.End() {
		cache.CutList(it1)
		cache.SetPrev(it1, it2)
	}
}

func (self *CacheVar_t[Key_t, Mapped_t]) Create(ts time.Time, ttl time.Duration, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.CreateBack(
		key,
		func(p *VarValue_t[Mapped_t]) {
			p.ts = ts.Add(ttl)
			value_init(&p.Value)
		},
		func(p *VarValue_t[Mapped_t]) {
			value_update(&p.Value)
		},
	)
	if ok {
		self.find_place()
	}
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) Push(ts time.Time, ttl time.Duration, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.PushBack(
		key,
		func(p *VarValue_t[Mapped_t]) {
			p.ts = ts.Add(ttl)
			value_init(&p.Value)
		},
		func(p *VarValue_t[Mapped_t]) {
			p.ts = ts.Add(ttl)
			value_update(&p.Value)
		},
	)
	self.find_place()
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) Update(ts time.Time, ttl time.Duration, key Key_t, value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.FindBack(key)
	if ok {
		it.Value.ts = ts.Add(ttl)
		value_update(&it.Value.Value)
		self.find_place()
	}
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) Refresh(ts time.Time, ttl time.Duration, key Key_t) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.FindBack(key)
	if ok {
		it.Value.ts = ts.Add(ttl)
		self.find_place()
	}
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) Replace(ts time.Time, key Key_t, value_update func(*Mapped_t)) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.Find(key)
	if ok {
		value_update(&it.Value.Value)
	}
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) Find(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.Find(key)
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) Remove(ts time.Time, key Key_t) (it *cache.Value_t[Key_t, VarValue_t[Mapped_t]], ok bool) {
	self.Flush(ts)
	it, ok = self.cx.Remove(key)
	return
}

func (self *CacheVar_t[Key_t, Mapped_t]) LeastTs(ts time.Time) (time.Time, bool) {
	self.Flush(ts)
	if self.cx.Size() > 0 {
		return self.cx.Front().Value.ts, true
	}
	return time.Time{}, false
}

func (self *CacheVar_t[Key_t, Mapped_t]) Range(ts time.Time, f func(key Key_t, value Mapped_t) bool) {
	self.Flush(ts)
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if f(it.Key, it.Value.Value) == false {
			return
		}
	}
}

func (self *CacheVar_t[Key_t, Mapped_t]) RangeTs(ts time.Time, f func(key Key_t, value Mapped_t, ts time.Time) bool) {
	self.Flush(ts)
	for it := self.cx.Front(); it != self.cx.End(); it = it.Next() {
		if f(it.Key, it.Value.Value, it.Value.ts) == false {
			return
		}
	}
}

func (self *CacheVar_t[Key_t, Mapped_t]) Size(ts time.Time) int {
	self.Flush(ts)
	return self.cx.Size()
}

func (self *CacheVar_t[Key_t, Mapped_t]) Limit() int {
	return self.limit
}
