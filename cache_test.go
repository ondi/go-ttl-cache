//
//
//

package cache

import "fmt"
import "time"
import "testing"

func ExampleTtlCache1() {
	var ok bool
	var e Evict_t
	c := NewSync(2, time.Second)
	
	c.Update(time.Now(), 1, 1, &e)
	c.Update(time.Now(), 2, 2, &e)
	c.Get(time.Now(), 1, &e)
	c.Update(time.Now(), 3, 3, &e)
	_, ok = c.Get(time.Now(), 1, &e)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 2, &e)
	fmt.Printf("%v\n", ok)
	_, ok = c.Get(time.Now(), 3, &e)
	fmt.Printf("%v\n", ok)
/* Output:
true
false
true
*/
}

func TtlCacheTest1(t * testing.T) {

}
