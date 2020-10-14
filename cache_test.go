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
	c := NewSync(2, time.Second, Drop)

	c.Create(time.Now(), 1, func() interface{} { return 1 })
	c.Get(time.Now(), 1)
	c.Create(time.Now(), 2, func() interface{} { return 2 })
	c.Get(time.Now(), 1)
	c.Create(time.Now(), 3, func() interface{} { return 3 })
	c.Update(time.Now(), 4, func(interface{}) interface{} { return 4 })
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
	c := NewSync(2, time.Second, Drop)

	c.Write(time.Now(), 1, func() interface{} { return 10 }, func(prev interface{}) interface{} { return prev })
	c.Get(time.Now(), 1)
	c.Write(time.Now(), 2, func() interface{} { return 20 }, func(prev interface{}) interface{} { return prev })
	c.Get(time.Now(), 1)
	c.Write(time.Now(), 3, func() interface{} { return 30 }, func(prev interface{}) interface{} { return prev })
	c.Update(time.Now(), 4, func(interface{}) interface{} { return 40 })
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
	c := NewSync(10, time.Second, Drop)

	c.Create(time.Now(), 1, func() interface{} { return 10 })
	c.Write(time.Now(), 2, func() interface{} { return 20 }, func(prev interface{}) interface{} { return prev })
	c.Create(time.Now(), 3, func() interface{} { return 30 })
	c.Write(time.Now(), 4, func() interface{} { return 40 }, func(prev interface{}) interface{} { return prev })
	c.Create(time.Now(), 5, func() interface{} { return 50 })
	c.Write(time.Now(), 6, func() interface{} { return 60 }, func(prev interface{}) interface{} { return prev })

	c.Update(time.Now(), 1, func(p interface{}) interface{} { return p.(int) + 100 })
	c.Update(time.Now(), 7, func(p interface{}) interface{} { return p.(int) + 100 })

	c.Range(time.Now(), func(key interface{}, value interface{}) bool {
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
