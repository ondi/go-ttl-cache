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

	c.Create(time.Now(), 1, func(p *int) { *p = 1 }, func(p *int) {})
	c.Refresh(time.Now(), 1)
	c.Create(time.Now(), 2, func(p *int) { *p = 2 }, func(p *int) {})
	c.Refresh(time.Now(), 1)
	c.Create(time.Now(), 3, func(p *int) { *p = 3 }, func(p *int) {})
	c.Update(time.Now(), 4, func(p *int) { *p = 4 })
	_, ok := c.Refresh(time.Now(), 1)
	assert.Assert(t, ok == true)
	_, ok = c.Refresh(time.Now(), 2)
	assert.Assert(t, ok == false)
	_, ok = c.Refresh(time.Now(), 3)
	assert.Assert(t, ok == true)
}

func Test_ttl_cache2(t *testing.T) {
	c := NewSync(2, time.Second, Drop[int, int])

	c.Push(time.Now(), 1, func(p *int) { *p = 10 }, func(p *int) {})
	c.Refresh(time.Now(), 1)
	c.Push(time.Now(), 2, func(p *int) { *p = 20 }, func(p *int) {})
	c.Refresh(time.Now(), 1)
	c.Push(time.Now(), 3, func(p *int) { *p = 30 }, func(p *int) {})
	c.Update(time.Now(), 4, func(p *int) { *p = 40 })
	_, ok := c.Refresh(time.Now(), 1)
	assert.Assert(t, ok == true)
	_, ok = c.Refresh(time.Now(), 2)
	assert.Assert(t, ok == false)
	_, ok = c.Refresh(time.Now(), 3)
	assert.Assert(t, ok == true)
}

func Example_ttl_cache3() {
	c := NewSync(10, time.Second, Drop[int, int])

	c.Create(time.Now(), 1, func(p *int) { *p = 10 }, func(p *int) {})
	c.Push(time.Now(), 2, func(p *int) { *p = 20 }, func(p *int) {})
	c.Create(time.Now(), 3, func(p *int) { *p = 30 }, func(p *int) {})
	c.Push(time.Now(), 4, func(p *int) { *p = 40 }, func(p *int) {})
	c.Create(time.Now(), 5, func(p *int) { *p = 50 }, func(p *int) {})
	c.Push(time.Now(), 6, func(p *int) { *p = 60 }, func(p *int) {})

	c.Update(time.Now(), 1, func(p *int) { *p = *p + 100 })
	c.Update(time.Now(), 7, func(p *int) { *p = *p + 100 })

	c.Range(
		time.Now(),
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
