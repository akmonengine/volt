package volt

import (
	"fmt"
	"maps"
	"slices"
)

func getStorage[T ComponentInterface](world *World) *ComponentsStorage[T] {
	var t T
	componentId := t.GetComponentId()

	if _, ok := world.ComponentsRegistry[componentId]; !ok {
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
	s := world.storage[componentId]

	if s == nil {
		return nil, fmt.Errorf("no storage for component id %v", componentId)
	}

	return s, nil
}

type storage interface {
	getType() ComponentId
	getArchetypes() []ArchetypeId
	hasArchetype(archetypeId ArchetypeId) bool
	add(archetypeId ArchetypeId, component ComponentInterface) int
	set(archetypeId ArchetypeId, key int, component ComponentInterface)
	get(archetypeId ArchetypeId, key int) any
	copy(oldArchetypeId ArchetypeId, archetypeId ArchetypeId, recordKey int) int
	size(archetypeId ArchetypeId) int
	moveLastToKey(archetypeId ArchetypeId, recordKey int)
	delete(archetypeId ArchetypeId, key int)
}

type ArchetypesComponentsEntities[T ComponentInterface] map[ArchetypeId][]T

type ComponentsStorage[T ComponentInterface] struct {
	componentId                  ComponentId
	archetypesComponentsEntities ArchetypesComponentsEntities[T]
}

func (c *ComponentsStorage[T]) getType() ComponentId {
	return c.componentId
}

func (c *ComponentsStorage[T]) getArchetypes() []ArchetypeId {
	return slices.Collect(maps.Keys(c.archetypesComponentsEntities))
}

func (c *ComponentsStorage[T]) hasArchetype(archetypeId ArchetypeId) bool {
	if _, ok := c.archetypesComponentsEntities[archetypeId]; !ok {
		return false
	}

	return true
}

func (c *ComponentsStorage[T]) size(archetypeId ArchetypeId) int {
	return len(c.archetypesComponentsEntities[archetypeId])
}

func (c *ComponentsStorage[T]) add(archetypeId ArchetypeId, component ComponentInterface) int {
	// this function could be simplified using:
	// c.size(archetypeId) - 1
	// but to reduce the usage of mapaccess we compute the size ourselves instead of calling c.size

	c.archetypesComponentsEntities[archetypeId] = append(c.archetypesComponentsEntities[archetypeId], component.(T))

	return len(c.archetypesComponentsEntities[archetypeId]) - 1

}

func (c *ComponentsStorage[T]) copy(oldArchetypeId ArchetypeId, archetypeId ArchetypeId, recordKey int) int {
	return c.add(archetypeId, c.archetypesComponentsEntities[oldArchetypeId][recordKey])
}

func (c *ComponentsStorage[T]) set(archetypeId ArchetypeId, key int, component ComponentInterface) {
	c.archetypesComponentsEntities[archetypeId][key] = component.(T)
}

func (c *ComponentsStorage[T]) get(archetypeId ArchetypeId, key int) any {
	return &c.archetypesComponentsEntities[archetypeId][key]
}

func (c *ComponentsStorage[T]) moveLastToKey(archetypeId ArchetypeId, recordKey int) {
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

func (c *ComponentsStorage[T]) delete(archetypeId ArchetypeId, key int) {
	if key < c.size(archetypeId) {
		data := c.archetypesComponentsEntities[archetypeId]
		c.archetypesComponentsEntities[archetypeId] = append(data[:key], data[key+1:]...)
	}
}
