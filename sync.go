//
//
//

package cache

import (
	"sync"
	"time"
)

type SyncCache_t[Key_t comparable, Mapped_t any] struct {
	mx sync.Mutex
	cx *Cache_t[Key_t, Mapped_t]
}

func NewSync[Key_t comparable, Mapped_t any](limit int, ttl time.Duration, evict Evict[Key_t, Mapped_t]) *SyncCache_t[Key_t, Mapped_t] {
	return &SyncCache_t[Key_t, Mapped_t]{
		cx: New(limit, ttl, evict),
	}
}

func (self *SyncCache_t[Key_t, Mapped_t]) Flush(ts time.Time) {
	self.mx.Lock()
	self.cx.Flush(ts)
	self.mx.Unlock()
}

func (self *SyncCache_t[Key_t, Mapped_t]) Create(ts time.Time, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Create(ts, key, value_init, value_update)
	res = it.Value.Value
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) CreateTtl(ts time.Time, ttl time.Duration, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.CreateTtl(ts, ttl, key, value_init, value_update)
	res = it.Value.Value
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Push(ts time.Time, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Push(ts, key, value_init, value_update)
	res = it.Value.Value
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) PushTtl(ts time.Time, ttl time.Duration, key Key_t, value_init func(*Mapped_t), value_update func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.PushTtl(ts, ttl, key, value_init, value_update)
	res = it.Value.Value
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Update(ts time.Time, key Key_t, value func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Update(ts, key, value)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) UpdateTtl(ts time.Time, ttl time.Duration, key Key_t, value func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.UpdateTtl(ts, ttl, key, value)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Refresh(ts time.Time, key Key_t) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Refresh(ts, key)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) RefreshTtl(ts time.Time, ttl time.Duration, key Key_t) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.RefreshTtl(ts, ttl, key)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Replace(ts time.Time, key Key_t, value func(*Mapped_t)) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Replace(ts, key, value)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Find(ts time.Time, key Key_t) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Find(ts, key)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Remove(ts time.Time, key Key_t) (res Mapped_t, ok bool) {
	self.mx.Lock()
	it, ok := self.cx.Remove(ts, key)
	if ok {
		res = it.Value.Value
	}
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) LeastTs(ts time.Time) (diff time.Time, ok bool) {
	self.mx.Lock()
	diff, ok = self.cx.LeastTs(ts)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Range(ts time.Time, f func(key Key_t, value Mapped_t) bool) {
	self.mx.Lock()
	self.cx.Range(ts, f)
	self.mx.Unlock()
}

func (self *SyncCache_t[Key_t, Mapped_t]) RangeTs(ts time.Time, f func(key Key_t, value Mapped_t, ts time.Time) bool) {
	self.mx.Lock()
	self.cx.RangeTs(ts, f)
	self.mx.Unlock()
}

func (self *SyncCache_t[Key_t, Mapped_t]) Size(ts time.Time) (res int) {
	self.mx.Lock()
	res = self.cx.Size(ts)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t[Key_t, Mapped_t]) Limit() int {
	return self.cx.Limit()
}

func (self *SyncCache_t[Key_t, Mapped_t]) TTL() time.Duration {
	return self.cx.TTL()
}
