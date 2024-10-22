package lru

// LRU cache. Not thread-safe.
type Cache[K comparable, V any] struct {
	keys    []K
	vals    []V
	lastUse []uint64
	evicted func(K, V)
	tick    uint64
}

func NewCache[K comparable, V any](capacity int, evicted ...func(K, V)) *Cache[K, V] {
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

func (c *Cache[K, V]) nextTick() uint64 {
	idx := c.tick
	c.tick++
	return idx
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

	idx := 0
	lastUse := c.lastUse[idx]

	for i := 1; i < l; i++ {
		if u := c.lastUse[i]; u < lastUse {
			idx = i
			lastUse = u
		}
	}

	c.remove(idx)
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
	c.vals[idx], c.vals[end] = c.vals[end], c.vals[idx]

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
