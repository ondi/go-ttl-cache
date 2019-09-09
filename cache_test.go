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
	
	c.UpdateFront(time.Now(), 1, func() interface{} {return 1})
	c.FindFront(time.Now(), 1)
	c.UpdateFront(time.Now(), 2, func() interface{} {return 2})
	c.FindFront(time.Now(), 1)
	c.UpdateFront(time.Now(), 3, func() interface{} {return 3})
	_, ok = c.FindFront(time.Now(), 1)
	fmt.Printf("%v\n", ok)
	_, ok = c.FindFront(time.Now(), 2)
	fmt.Printf("%v\n", ok)
	_, ok = c.FindFront(time.Now(), 3)
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
	
	c.UpdateFront(time.Now(), 1, func() interface{} {return 10})
	c.FindFront(time.Now(), 1)
	c.UpdateFront(time.Now(), 2, func() interface{} {return 20})
	c.FindFront(time.Now(), 1)
	c.UpdateFront(time.Now(), 3, func() interface{} {return 30})
	_, ok = c.FindFront(time.Now(), 1)
	fmt.Printf("%v\n", ok)
	_, ok = c.FindFront(time.Now(), 2)
	fmt.Printf("%v\n", ok)
	_, ok = c.FindFront(time.Now(), 3)
	fmt.Printf("%v\n", ok)
/* Output:
true
false
true
*/
}

func ExampleTtlCache3() {
	c := NewSync(5, time.Second, Drop)
	
	c.UpdateBack(time.Now(), 1, func() interface{} {return 10})
	c.UpdateBack(time.Now(), 2, func() interface{} {return 20})
	c.UpdateBack(time.Now(), 3, func() interface{} {return 30})
	c.UpdateBack(time.Now(), 4, func() interface{} {return 40})
	c.UpdateBack(time.Now(), 5, func() interface{} {return 50})
	c.RangeFrontBack(time.Now(), func(key interface{}, value interface{}) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
	c.RangeBackFront(time.Now(), func(key interface{}, value interface{}) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
/* Output:
1 10
2 20
3 30
4 40
5 50
5 50
4 40
3 30
2 20
1 10
*/
}

func TtlCacheTest1(t * testing.T) {

}
