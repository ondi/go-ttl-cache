//
//
//

package cache

import "fmt"
import "time"
import "testing"

import "github.com/ondi/go-queue"

func ExampleTtlCache1() {
	var ok bool
	q := queue.New(1000)
	c := NewSync(2, time.Second, q.PushBackNoWait)
	
	c.Update(time.Now(), 1, func() interface{} {return 1})
	c.Get(time.Now(), 1)
	c.Update(time.Now(), 2, func() interface{} {return 2})
	c.Get(time.Now(), 1)
	c.Update(time.Now(), 3, func() interface{} {return 3})
	_, ok = c.Get(time.Now(), 1)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 2)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 3)
	fmt.Printf("%v\n", ok)
/* Output:
true
false
true
*/
}

func ExampleTtlCache2() {
	var ok bool
	c := NewSync(2, time.Second, Drop)
	
	c.Update(time.Now(), 1, func() interface{} {return 10})
	c.Get(time.Now(), 1)
	c.Update(time.Now(), 2, func() interface{} {return 20})
	c.Get(time.Now(), 1)
	c.Update(time.Now(), 3, func() interface{} {return 30})
	_, ok = c.Get(time.Now(), 1)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 2)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 3)
	fmt.Printf("%v\n", ok)
/* Output:
true
false
true
*/
}

func ExampleTtlCache3() {
	c := NewSync(5, time.Second, Drop)
	
	c.Update(time.Now(), 1, func() interface{} {return 10})
	c.Update(time.Now(), 2, func() interface{} {return 20})
	c.Update(time.Now(), 3, func() interface{} {return 30})
	c.Update(time.Now(), 4, func() interface{} {return 40})
	c.Update(time.Now(), 5, func() interface{} {return 50})
	c.Range(time.Now(), func(key interface{}, value interface{}) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
/* Output:
1 10
2 20
3 30
4 40
5 50
*/
}

func TtlCacheTest1(t * testing.T) {

}
