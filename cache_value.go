//
//
//

package cache

type Evict[Key_t comparable, Mapped_t any] func(key Key_t, value Mapped_t)

func Drop[Key_t comparable, Mapped_t any](Key_t, Mapped_t) {}
