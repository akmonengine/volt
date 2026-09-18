package volt

// pool hands out entity slots. A slot is the index part of an EntityId; the
// generation that keeps a recycled slot's handle unique lives in the entity
// record, not here. Freed slots are reused LIFO before any new slot is opened.
type pool struct {
	free []uint32
	next uint32
}

// Get returns a slot, and whether it is recycled (previously freed) or brand new.
func (pool *pool) Get() (index uint32, recycled bool) {
	if n := len(pool.free); n > 0 {
		index = pool.free[n-1]
		pool.free = pool.free[:n-1]

		return index, true
	}

	index = pool.next
	pool.next++

	return index, false
}

// Recycle gives a slot back, to be reused by the next Get.
func (pool *pool) Recycle(index uint32) {
	pool.free = append(pool.free, index)
}

// Count returns the number of freed slots waiting to be recycled.
func (pool *pool) Count() int {
	return len(pool.free)
}
