//
//
//

package cache

import (
	"fmt"
	"testing"
	"time"
)

func Example_ttl_cache1() {
	var ok bool
	c := NewSync(2, time.Second, Drop[int, int])

	c.Create(time.Now(), 1, func() int { return 1 }, func(p *int) {})
	c.Get(time.Now(), 1)
	c.Create(time.Now(), 2, func() int { return 2 }, func(p *int) {})
	c.Get(time.Now(), 1)
	c.Create(time.Now(), 3, func() int { return 3 }, func(p *int) {})
	c.Update(time.Now(), 4, func(p *int) { *p = 4 })
	_, ok = c.Get(time.Now(), 1)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 2)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 3)
	fmt.Printf("%v\n", ok)
	// Output:
	// true
	// false
	// true
}

func Example_ttl_cache2() {
	var ok bool
	c := NewSync(2, time.Second, Drop[int, int])

	c.Push(time.Now(), 1, func() int { return 10 }, func(p *int) {})
	c.Get(time.Now(), 1)
	c.Push(time.Now(), 2, func() int { return 20 }, func(p *int) {})
	c.Get(time.Now(), 1)
	c.Push(time.Now(), 3, func() int { return 30 }, func(p *int) {})
	c.Update(time.Now(), 4, func(p *int) { *p = 40 })
	_, ok = c.Get(time.Now(), 1)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 2)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 3)
	fmt.Printf("%v\n", ok)
	// Output:
	// true
	// false
	// true
}

func Example_ttl_cache3() {
	c := NewSync(10, time.Second, Drop[int, int])

	c.Create(time.Now(), 1, func() int { return 10 }, func(p *int) {})
	c.Push(time.Now(), 2, func() int { return 20 }, func(p *int) {})
	c.Create(time.Now(), 3, func() int { return 30 }, func(p *int) {})
	c.Push(time.Now(), 4, func() int { return 40 }, func(p *int) {})
	c.Create(time.Now(), 5, func() int { return 50 }, func(p *int) {})
	c.Push(time.Now(), 6, func() int { return 60 }, func(p *int) {})

	c.Update(time.Now(), 1, func(p *int) { *p = *p + 100 })
	c.Update(time.Now(), 7, func(p *int) { *p = *p + 100 })

	c.Range(time.Now(), func(key int, value int) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
	// Output:
	// 2 20
	// 3 30
	// 4 40
	// 5 50
	// 6 60
	// 1 110
}

func TtlCacheTest1(t *testing.T) {

}
