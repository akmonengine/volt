package volt

import (
	"fmt"
	"maps"
	"slices"
)

func getStorage[T ComponentInterface](world *World) *ComponentsStorage[T] {
	var t T
	componentId := t.GetComponentId()

	if world.componentsRegistry[componentId] == nil {
		return nil
	}

	if world.storage[componentId] == nil {
		s := &ComponentsStorage[T]{
			componentId:                  componentId,
			archetypesComponentsEntities: make(ArchetypesComponentsEntities[T]),
		}
		world.storage[componentId] = s
	}

	return world.storage[componentId].(*ComponentsStorage[T])
}

func (world *World) getStorageForComponentId(componentId ComponentId) (storage, error) {
	if componentId >= TAGS_INDICES {
		return nil, fmt.Errorf("componentId %d overflow the range allowed [%d;%d]", componentId, COMPONENTS_INDICES, TAGS_INDICES)
	}

	s := world.storage[componentId]

	if s == nil {
		return nil, fmt.Errorf("no storage for component id %v", componentId)
	}

	return s, nil
}

type storage interface {
	getType() ComponentId
	getArchetypes() []archetypeId
	hasArchetype(archetypeId archetypeId) bool
	add(archetypeId archetypeId, component ComponentInterface) int
	set(archetypeId archetypeId, key int, component ComponentInterface)
	get(archetypeId archetypeId, key int) any
	copy(oldArchetypeId archetypeId, archetypeId archetypeId, recordKey int) int
	size(archetypeId archetypeId) int
	moveLastToKey(archetypeId archetypeId, recordKey int)
	delete(archetypeId archetypeId, key int)
}

type ArchetypesComponentsEntities[T ComponentInterface] map[archetypeId][]T

type ComponentsStorage[T ComponentInterface] struct {
	componentId                  ComponentId
	archetypesComponentsEntities ArchetypesComponentsEntities[T]
}

func (c *ComponentsStorage[T]) getType() ComponentId {
	return c.componentId
}

func (c *ComponentsStorage[T]) getArchetypes() []archetypeId {
	return slices.Collect(maps.Keys(c.archetypesComponentsEntities))
}

func (c *ComponentsStorage[T]) hasArchetype(archetypeId archetypeId) bool {
	if _, ok := c.archetypesComponentsEntities[archetypeId]; !ok {
		return false
	}

	return true
}

func (c *ComponentsStorage[T]) size(archetypeId archetypeId) int {
	return len(c.archetypesComponentsEntities[archetypeId])
}

func (c *ComponentsStorage[T]) add(archetypeId archetypeId, component ComponentInterface) int {
	// this function could be simplified using:
	// c.size(archetypeId) - 1
	// but to reduce the usage of mapaccess we compute the size ourselves instead of calling c.size

	c.archetypesComponentsEntities[archetypeId] = append(c.archetypesComponentsEntities[archetypeId], component.(T))

	return len(c.archetypesComponentsEntities[archetypeId]) - 1

}

func (c *ComponentsStorage[T]) copy(oldArchetypeId archetypeId, archetypeId archetypeId, recordKey int) int {
	return c.add(archetypeId, c.archetypesComponentsEntities[oldArchetypeId][recordKey])
}

func (c *ComponentsStorage[T]) set(archetypeId archetypeId, key int, component ComponentInterface) {
	c.archetypesComponentsEntities[archetypeId][key] = component.(T)
}

func (c *ComponentsStorage[T]) get(archetypeId archetypeId, key int) any {
	return &c.archetypesComponentsEntities[archetypeId][key]
}

func (c *ComponentsStorage[T]) moveLastToKey(archetypeId archetypeId, recordKey int) {
	// this function could be simplified using:
	// 	lastKey := c.size(archetypeId) - 1
	// 	c.set(archetypeId, recordKey, c.archetypesComponentsEntities[archetypeId][lastKey])
	//	c.delete(archetypeId, lastKey)
	// but this would imply lot of map access, reducing the performances

	data := c.archetypesComponentsEntities[archetypeId]
	size := len(data)
	lastKey := size - 1

	data[recordKey] = data[lastKey]

	if lastKey < size {
		c.archetypesComponentsEntities[archetypeId] = append(data[:lastKey], data[lastKey+1:]...)
	}
}

func (c *ComponentsStorage[T]) delete(archetypeId archetypeId, key int) {
	if key < c.size(archetypeId) {
		data := c.archetypesComponentsEntities[archetypeId]
		c.archetypesComponentsEntities[archetypeId] = append(data[:key], data[key+1:]...)
	}
}
