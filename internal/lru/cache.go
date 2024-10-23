package lru

import (
	"iter"
	"math"
)

// LRU cache. Not thread-safe.
type Cache[K comparable, V any] struct {
	keys    []K
	vals    []V
	lastUse []uint64
	tick    uint64
	evicted func(K, V)
}

func NewCache[K comparable, V any](capacity int, evicted ...func(key K, val V)) *Cache[K, V] {
	c := &Cache[K, V]{
		keys:    make([]K, 0, capacity),
		vals:    make([]V, 0, capacity),
		lastUse: make([]uint64, 0, capacity),
	}

	if len(evicted) > 0 {
		c.evicted = evicted[0]
	}

	return c
}

func (c *Cache[K, V]) Len() int {
	return len(c.keys)
}

func (c *Cache[K, V]) Cap() int {
	return len(c.keys)
}

func (c *Cache[K, V]) Resize(capacity int) {
	if cap(c.keys) == capacity {
		return
	}

	for capacity < c.Len() {
		c.removeOldest()
	}

	keys := append(make([]K, 0, capacity), c.keys...)
	vals := append(make([]V, 0, capacity), c.vals...)
	lastUse := append(make([]uint64, 0, capacity), c.lastUse...)

	c.Reset()

	c.keys = keys
	c.vals = vals
	c.lastUse = lastUse
}

// Clear cache without notice. To clear cache and notify each evict, use RemoveAll.
func (c *Cache[K, V]) Reset() {
	clear(c.keys)
	clear(c.vals)
	clear(c.lastUse)

	c.keys = c.keys[:0]
	c.vals = c.vals[:0]
	c.lastUse = c.lastUse[:0]
	c.tick = 0
}

func (c *Cache[K, V]) Has(key K) (ok bool) {
	for i := range c.keys {
		if c.keys[i] == key {
			return true
		}
	}

	return
}

func (c *Cache[K, V]) Get(key K) (val V, ok bool) {
	for i := range c.keys {
		if c.keys[i] == key {
			c.lastUse[i] = c.nextTick()
			return c.vals[i], true
		}
	}

	return
}

func (c *Cache[K, V]) GetOrSet(key K, setter func(K) (V, error)) (val V, err error) {
	var ok bool

	if val, ok = c.Get(key); ok {
		return
	}

	if val, err = setter(key); err == nil {
		c.append(key, val)
	}

	return
}

func (c *Cache[K, V]) Set(key K, val V) (ok bool) {
	if c.Has(key) {
		return
	}

	c.append(key, val)

	return true
}

func (c *Cache[K, V]) Replace(key K, val V) (existed bool) {
	for i := range c.keys {
		if c.keys[i] == key {
			c.vals[i], val = val, c.vals[i]
			c.lastUse[i] = c.nextTick()
			c.evict(key, val)
			return true
		}
	}

	c.append(key, val)

	return
}

func (c *Cache[K, V]) Remove(key K) (existed bool) {
	for i := range c.keys {
		if c.keys[i] == key {
			c.remove(i)
			return true
		}
	}

	return
}

// Clear cache and notify each evict. To clear cache without notice, use Reset.
func (c *Cache[K, V]) RemoveAll() {
	for i := range c.keys {
		c.evict(c.keys[i], c.vals[i])
	}

	c.Reset()
}

// Iterate all items in no particular order.
func (c *Cache[K, V]) Iterate() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := range c.keys {
			if !yield(c.keys[i], c.vals[i]) {
				return
			}
		}
	}
}

// Iterate all items in ascending order.
func (c *Cache[K, V]) IterateAsc() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		tick := c.oldestTick()

		for range c.keys {
			idx, ok := c.find(tick)

			if !ok || !yield(c.keys[idx], c.vals[idx]) {
				return
			}

			tick++
		}
	}
}

// Iterate all items in descending order.
func (c *Cache[K, V]) IterateDesc() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		tick := c.tick - 1

		for range c.keys {
			idx, ok := c.find(tick)

			if !ok || !yield(c.keys[idx], c.vals[idx]) {
				return
			}

			tick--
		}
	}
}

func (c *Cache[K, V]) append(key K, val V) {
	if len(c.keys) >= cap(c.keys) {
		c.removeOldest()
	}

	c.keys = append(c.keys, key)
	c.vals = append(c.vals, val)
	c.lastUse = append(c.lastUse, c.nextTick())
}

func (c *Cache[K, V]) removeOldest() {
	l := len(c.keys)

	if l == 0 {
		return
	}

	idx, ok := c.find(c.oldestTick())

	if ok {
		c.remove(idx)
	} else {
		c.repair()
		c.removeOldest()
	}
}

func (c *Cache[K, V]) oldestTick() uint64 {
	return c.tick - uint64(c.Len())
}

func (c *Cache[K, V]) find(tick uint64) (idx int, ok bool) {
	for i := range c.lastUse {
		if c.lastUse[i] == tick {
			return i, true
		}
	}

	return
}

func (c *Cache[K, V]) remove(idx int) {
	var (
		key K
		val V
	)

	end := len(c.keys) - 1

	// Swap with zero values
	key, c.keys[idx], c.keys[end] = c.keys[idx], c.keys[end], key
	val, c.vals[idx], c.vals[end] = c.vals[idx], c.vals[end], val
	c.lastUse[idx], c.lastUse[end] = c.lastUse[end], c.lastUse[idx]

	c.keys = c.keys[:end]
	c.vals = c.vals[:end]
	c.lastUse = c.lastUse[:end]

	c.evict(key, val)
}

func (c *Cache[K, V]) evict(key K, val V) {
	if c.evicted != nil {
		c.evicted(key, val)
	}
}

func (c *Cache[K, V]) nextTick() uint64 {
	c.preventTickOverflow()
	idx := c.tick
	c.tick++
	return idx
}

func (c *Cache[K, V]) preventTickOverflow() {
	if c.tick != math.MaxUint64 {
		return
	}

	tick := c.oldestTick()

	for i := range c.lastUse {
		if c.lastUse[i] < tick {
			c.repair()
			return
		}

		c.lastUse[i] += tick
	}

	c.tick = uint64(c.Len())
}

func (c *Cache[K, V]) repair() {
	for i := range c.lastUse {
		c.lastUse[i] = uint64(i)
	}

	c.tick = uint64(c.Len())
}
