//
//
//

package cache

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_ttl_cache1(t *testing.T) {
	c := NewSync(2, time.Second, Drop[int, int])

	ts := time.Now()
	c.Create(ts, 1, func(p *int) { *p = 1 }, func(p *int) {})
	c.Refresh(ts, 1)
	c.Create(ts, 2, func(p *int) { *p = 2 }, func(p *int) {})
	c.Refresh(ts, 1)
	c.Create(ts, 3, func(p *int) { *p = 3 }, func(p *int) {})
	c.Update(ts, 4, func(p *int) { *p = 4 })
	_, ok := c.Refresh(ts, 1)
	assert.Assert(t, ok == true)
	_, ok = c.Refresh(ts, 2)
	assert.Assert(t, ok == false)
	_, ok = c.Refresh(ts, 3)
	assert.Assert(t, ok == true)
}

func Test_ttl_cache2(t *testing.T) {
	c := NewSync(2, time.Second, Drop[int, int])

	ts := time.Now()
	c.Push(ts, 1, func(p *int) { *p = 10 }, func(p *int) {})
	c.Refresh(ts, 1)
	c.Push(ts, 2, func(p *int) { *p = 20 }, func(p *int) {})
	c.Refresh(ts, 1)
	c.Push(ts, 3, func(p *int) { *p = 30 }, func(p *int) {})
	c.Update(ts, 4, func(p *int) { *p = 40 })
	_, ok := c.Refresh(ts, 1)
	assert.Assert(t, ok == true)
	_, ok = c.Refresh(ts, 2)
	assert.Assert(t, ok == false)
	_, ok = c.Refresh(ts, 3)
	assert.Assert(t, ok == true)
}

func Example_ttl_cache3() {
	c := NewSync(10, time.Second, Drop[int, int])

	ts := time.Now()
	c.Create(ts, 1, func(p *int) { *p = 10 }, func(p *int) {})
	c.Push(ts, 2, func(p *int) { *p = 20 }, func(p *int) {})
	c.Create(ts, 3, func(p *int) { *p = 30 }, func(p *int) {})
	c.Push(ts, 4, func(p *int) { *p = 40 }, func(p *int) {})
	c.Create(ts, 5, func(p *int) { *p = 50 }, func(p *int) {})
	c.Push(ts, 6, func(p *int) { *p = 60 }, func(p *int) {})

	c.Update(ts, 1, func(p *int) { *p = *p + 100 })
	c.Update(ts, 7, func(p *int) { *p = *p + 100 })

	c.Range(
		ts,
		func(key int, value int) bool {
			fmt.Printf("%v %v\n", key, value)
			return true
		},
	)
	// Output:
	// 2 20
	// 3 30
	// 4 40
	// 5 50
	// 6 60
	// 1 110
}

func Example_ttl_cache4() {
	c := NewSyncVar(10, Drop[int, int])

	ts := time.Now()
	c.Create(ts, 1*time.Second, 1, func(p *int) { *p = 10 }, func(p *int) {})
	c.Push(ts, 1*time.Second, 2, func(p *int) { *p = 20 }, func(p *int) {})
	c.Create(ts, 1*time.Second, 3, func(p *int) { *p = 30 }, func(p *int) {})
	c.Push(ts, 1*time.Second, 4, func(p *int) { *p = 40 }, func(p *int) {})
	c.Create(ts, 1*time.Second, 5, func(p *int) { *p = 50 }, func(p *int) {})
	c.Push(ts, 1*time.Second, 6, func(p *int) { *p = 60 }, func(p *int) {})

	c.Update(ts, 1*time.Second, 1, func(p *int) { *p = *p + 100 })
	c.Update(ts, 1*time.Second, 7, func(p *int) { *p = *p + 100 })

	c.Range(
		ts,
		func(key int, value int) bool {
			fmt.Printf("%v %v\n", key, value)
			return true
		},
	)
	// Output:
	// 2 20
	// 3 30
	// 4 40
	// 5 50
	// 6 60
	// 1 110
}

func Example_ttl_cache5() {
	c := NewSyncVar(10, Drop[int, int])

	ts := time.Now()
	c.Create(ts, 1*time.Second, 1, func(p *int) { *p = 10 }, func(p *int) {})
	c.Push(ts, 2*time.Second, 2, func(p *int) { *p = 20 }, func(p *int) {})
	c.Create(ts, 3*time.Second, 3, func(p *int) { *p = 30 }, func(p *int) {})
	c.Push(ts, 4*time.Second, 4, func(p *int) { *p = 40 }, func(p *int) {})
	c.Create(ts, 5*time.Second, 5, func(p *int) { *p = 50 }, func(p *int) {})
	c.Push(ts, 6*time.Second, 6, func(p *int) { *p = 60 }, func(p *int) {})

	c.Range(
		ts,
		func(key int, value int) bool {
			fmt.Printf("%v %v\n", key, value)
			return true
		},
	)
	// Output:
	// 6 60
	// 5 50
	// 4 40
	// 3 30
	// 2 20
	// 1 10
}
