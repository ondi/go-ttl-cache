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

func TtlCacheTest1(t * testing.T) {

}
