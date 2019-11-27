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
	
	c.Create(time.Now(), 1, func() interface{} {return 1})
	c.Get(time.Now(), 1)
	c.Create(time.Now(), 2, func() interface{} {return 2})
	c.Get(time.Now(), 1)
	c.Create(time.Now(), 3, func() interface{} {return 3})
	c.Update(time.Now(), 4, func(interface{}) interface{} {return 4})
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
	
	c.Push(time.Now(), 1, func() interface{} {return 10})
	c.Get(time.Now(), 1)
	c.Push(time.Now(), 2, func() interface{} {return 20})
	c.Get(time.Now(), 1)
	c.Push(time.Now(), 3, func() interface{} {return 30})
	c.Update(time.Now(), 4, func(interface{}) interface{} {return 40})
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
	c := NewSync(10, time.Second, Drop)
	
	c.Create(time.Now(), 1, func() interface{} {return 10})
	c.Push(time.Now(), 2, func() interface{} {return 20})
	c.Create(time.Now(), 3, func() interface{} {return 30})
	c.Push(time.Now(), 4, func() interface{} {return 40})
	c.Create(time.Now(), 5, func() interface{} {return 50})
	c.Push(time.Now(), 6, func() interface{} {return 60})
	
	c.Update(time.Now(), 1, func(p interface{}) interface{} {return p.(int) + 100})
	c.Update(time.Now(), 7, func(p interface{}) interface{} {return p.(int) + 100})

	c.Range(time.Now(), func(key interface{}, value interface{}) bool {
		fmt.Printf("%v %v\n", key, value)
		return true
	})
/* Output:
2 20
3 30
4 40
5 50
6 60
1 110
*/
}

func TtlCacheTest1(t * testing.T) {

}
