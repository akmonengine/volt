package volt

type pool struct {
	ids  []EntityId
	next EntityId
}

func (pool *pool) Get() EntityId {
	var entityId EntityId
	if len(pool.ids) > 0 {
		entityId = pool.ids[len(pool.ids)-1]
		pool.ids = pool.ids[:len(pool.ids)-1]
	} else {
		entityId = pool.next
		pool.next++
	}

	return entityId
}

func (pool *pool) Recycle(id EntityId) {
	pool.ids = append(pool.ids, id)
}

func (pool *pool) Count() int {
	return len(pool.ids)
}
