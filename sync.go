//
//
//

package cache

import (
	"sync"
	"time"
)

type SyncCache_t struct {
	mx sync.Mutex
	cx *Cache_t
}

func NewSync(limit int, ttl time.Duration, evict Evict) *SyncCache_t {
	return &SyncCache_t{
		cx: New(limit, ttl, evict),
	}
}

func (self *SyncCache_t) Flush(ts time.Time) {
	self.mx.Lock()
	self.cx.Flush(ts)
	self.mx.Unlock()
}

func (self *SyncCache_t) Create(ts time.Time, key interface{}, value_new func() interface{}, value_update func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Create(ts, key, value_new, value_update)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Create2(ts time.Time, key interface{}, value_new func() (interface{}, error), value_update func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.cx.Create2(ts, key, value_new, value_update)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Push(ts time.Time, key interface{}, value_new func() interface{}, value_update func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Push(ts, key, value_new, value_update)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Push2(ts time.Time, key interface{}, value_new func() (interface{}, error), value_update func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.cx.Push2(ts, key, value_new, value_update)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Update(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Update(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Update2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.cx.Update2(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Replace(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Replace(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Replace2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.cx.Replace2(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Get(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Get(ts, key)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Find(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Find(ts, key)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Remove(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.cx.Remove(ts, key)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) LeastTs(ts time.Time) (diff time.Time, ok bool) {
	self.mx.Lock()
	diff, ok = self.cx.LeastTs(ts)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Range(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.mx.Lock()
	self.cx.Range(ts, f)
	self.mx.Unlock()
}

func (self *SyncCache_t) RangeTs(ts time.Time, f func(key interface{}, value interface{}, ts time.Time) bool) {
	self.mx.Lock()
	self.cx.RangeTs(ts, f)
	self.mx.Unlock()
}

func (self *SyncCache_t) Size(ts time.Time) (res int) {
	self.mx.Lock()
	res = self.cx.Size(ts)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Limit() int {
	return self.cx.Limit()
}

func (self *SyncCache_t) TTL() time.Duration {
	return self.cx.TTL()
}
