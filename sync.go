//
//
//

package cache

import "sync"
import "time"

type SyncCache_t struct {
	mx sync.Mutex
	Cache_t
}

func NewSync(limit int, ttl time.Duration, evict Evict) (self * SyncCache_t) {
	self = &SyncCache_t{}
	self.Cache_t = *New(limit, ttl, evict)
	return
}

func (self * SyncCache_t) Flush(ts time.Time) {
	self.mx.Lock()
	self.Cache_t.Flush(ts)
	self.mx.Unlock()
}

func (self * SyncCache_t) Remove(key interface{}) {
	self.mx.Lock()
	self.Cache_t.Remove(key)
	self.mx.Unlock()
}

func (self * SyncCache_t) Create(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Create(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Push(ts time.Time, key interface{}, value func () interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Push(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Update(ts time.Time, key interface{}, value interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Update(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Get(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Get(ts, key)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Find(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Find(ts, key)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) LeastTs(ts time.Time) (t time.Time, ok bool) {
	self.mx.Lock()
	t, ok = self.Cache_t.LeastTs(ts)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Range(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.mx.Lock()
	self.Cache_t.Range(ts, f)
	self.mx.Unlock()
}

func (self * SyncCache_t) Size() (res int) {
	self.mx.Lock()
	res = self.Cache_t.Size()
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Limit() int {
	return self.Cache_t.Limit()
}

func (self * SyncCache_t) TTL() time.Duration {
	return self.Cache_t.TTL()
}
