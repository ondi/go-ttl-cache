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
	Cache_t
}

func NewSync(limit int, ttl time.Duration, evict Evict) (self *SyncCache_t) {
	self = &SyncCache_t{}
	self.Cache_t = *New(limit, ttl, evict)
	return
}

func (self *SyncCache_t) Flush(ts time.Time) {
	self.mx.Lock()
	self.Cache_t.Flush(ts)
	self.mx.Unlock()
}

func (self *SyncCache_t) Create(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Create(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Create2(ts time.Time, key interface{}, value func() (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.Cache_t.Create2(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Write(ts time.Time, key interface{}, value_new func() interface{}, value_update func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Write(ts, key, value_new, value_update)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Write2(ts time.Time, key interface{}, value_new func() (interface{}, error), value_update func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.Cache_t.Write2(ts, key, value_new, value_update)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Update(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Update(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Update2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.Cache_t.Update2(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Replace(ts time.Time, key interface{}, value func(interface{}) interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Replace(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Replace2(ts time.Time, key interface{}, value func(interface{}) (interface{}, error)) (res interface{}, ok bool, err error) {
	self.mx.Lock()
	res, ok, err = self.Cache_t.Replace2(ts, key, value)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Get(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Get(ts, key)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Find(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Find(ts, key)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Remove(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Remove(ts, key)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) LeastDiff(ts time.Time) (diff time.Duration, ok bool) {
	self.mx.Lock()
	diff, ok = self.Cache_t.LeastDiff(ts)
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Range(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.mx.Lock()
	self.Cache_t.Range(ts, f)
	self.mx.Unlock()
}

func (self *SyncCache_t) RangeTs(ts time.Time, f func(key interface{}, value interface{}, diff time.Duration) bool) {
	self.mx.Lock()
	self.Cache_t.RangeTs(ts, f)
	self.mx.Unlock()
}

func (self *SyncCache_t) Size() (res int) {
	self.mx.Lock()
	res = self.Cache_t.Size()
	self.mx.Unlock()
	return
}

func (self *SyncCache_t) Limit() int {
	return self.Cache_t.Limit()
}

func (self *SyncCache_t) TTL() time.Duration {
	return self.Cache_t.TTL()
}
