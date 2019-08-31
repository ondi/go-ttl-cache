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

func TtlCacheTest1(t * testing.T) {

}
