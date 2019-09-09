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

func (self * SyncCache_t) Remove(key interface{}) (ok bool) {
	self.mx.Lock()
	ok = self.Cache_t.Remove(key)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) CreateFront(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.CreateFront(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) CreateBack(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.CreateBack(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) PushFront(ts time.Time, key interface{}, value func () interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.PushFront(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) PushBack(ts time.Time, key interface{}, value func () interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.PushBack(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) UpdateFront(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.UpdateFront(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) UpdateBack(ts time.Time, key interface{}, value func() interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.UpdateBack(ts, key, value)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) FindFront(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.FindFront(ts, key)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) FindBack(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.FindBack(ts, key)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) Find(ts time.Time, key interface{}) (res interface{}, ok bool) {
	self.mx.Lock()
	res, ok = self.Cache_t.Find(ts, key)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) FrontTs(ts time.Time) (t time.Time, ok bool) {
	self.mx.Lock()
	t, ok = self.Cache_t.FrontTs(ts)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) BackTs(ts time.Time) (t time.Time, ok bool) {
	self.mx.Lock()
	t, ok = self.Cache_t.BackTs(ts)
	self.mx.Unlock()
	return
}

func (self * SyncCache_t) RangeFrontBack(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.mx.Lock()
	self.Cache_t.RangeFrontBack(ts, f)
	self.mx.Unlock()
}

func (self * SyncCache_t) RangeBackFront(ts time.Time, f func(key interface{}, value interface{}) bool) {
	self.mx.Lock()
	self.Cache_t.RangeBackFront(ts, f)
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
