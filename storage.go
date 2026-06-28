package volt

import (
	"fmt"
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
			archetypesComponentsEntities: make(ArchetypesComponentsEntities[T], 0),
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

// ArchetypesComponentsEntities stores, for each archetype, the column of T
// components (Structure of Arrays). It is indexed directly by archetypeId:
// archetypeId is dense (0, 1, 2, ...), so a slice avoids the hashing cost of a
// map on every storage access, in both read (queries) and write paths. A nil
// column means the archetype does not hold this component.
type ArchetypesComponentsEntities[T ComponentInterface] [][]T

type ComponentsStorage[T ComponentInterface] struct {
	componentId                  ComponentId
	archetypesComponentsEntities ArchetypesComponentsEntities[T]
}

func (c *ComponentsStorage[T]) getType() ComponentId {
	return c.componentId
}

func (c *ComponentsStorage[T]) getArchetypes() []archetypeId {
	var archetypes []archetypeId
	for id, column := range c.archetypesComponentsEntities {
		if column != nil {
			archetypes = append(archetypes, archetypeId(id))
		}
	}

	return archetypes
}

func (c *ComponentsStorage[T]) hasArchetype(archetypeId archetypeId) bool {
	return int(archetypeId) < len(c.archetypesComponentsEntities) && c.archetypesComponentsEntities[archetypeId] != nil
}

// getColumn returns the component column for archetypeId, or nil if this storage
// holds no data for it. It is bounds-safe: an archetype that does not contain
// this component (e.g. an optional component absent from the archetype) has an
// id beyond the columns slice, so a raw index would panic.
//
// Kept out of line on purpose: inlining it into the query hot loops perturbs
// their codegen and slows iteration. As a per-archetype call its cost is
// negligible, mirroring how the previous map access stayed out of line.
//
//go:noinline
func (c *ComponentsStorage[T]) getColumn(archetypeId archetypeId) []T {
	if int(archetypeId) >= len(c.archetypesComponentsEntities) {
		return nil
	}

	return c.archetypesComponentsEntities[archetypeId]
}

func (c *ComponentsStorage[T]) size(archetypeId archetypeId) int {
	if int(archetypeId) >= len(c.archetypesComponentsEntities) {
		return 0
	}

	return len(c.archetypesComponentsEntities[archetypeId])
}

// grow extends the columns slice so that archetypeId is a valid index.
// Growth is amortized through append, and new columns start as nil.
func (c *ComponentsStorage[T]) grow(archetypeId archetypeId) {
	for len(c.archetypesComponentsEntities) <= int(archetypeId) {
		c.archetypesComponentsEntities = append(c.archetypesComponentsEntities, nil)
	}
}

func (c *ComponentsStorage[T]) add(archetypeId archetypeId, component ComponentInterface) int {
	return c.addTyped(archetypeId, component.(T))
}

// addTyped appends a component without boxing it into ComponentInterface.
// The generic add/copy paths hold a concrete *ComponentsStorage[T], so they can
// call this directly and avoid one heap allocation per component added.
func (c *ComponentsStorage[T]) addTyped(archetypeId archetypeId, component T) int {
	c.grow(archetypeId)
	c.archetypesComponentsEntities[archetypeId] = append(c.archetypesComponentsEntities[archetypeId], component)

	return len(c.archetypesComponentsEntities[archetypeId]) - 1
}

func (c *ComponentsStorage[T]) copy(oldArchetypeId archetypeId, archetypeId archetypeId, recordKey int) int {
	return c.addTyped(archetypeId, c.archetypesComponentsEntities[oldArchetypeId][recordKey])
}

func (c *ComponentsStorage[T]) set(archetypeId archetypeId, key int, component ComponentInterface) {
	c.archetypesComponentsEntities[archetypeId][key] = component.(T)
}

func (c *ComponentsStorage[T]) get(archetypeId archetypeId, key int) any {
	return &c.archetypesComponentsEntities[archetypeId][key]
}

func (c *ComponentsStorage[T]) moveLastToKey(archetypeId archetypeId, recordKey int) {
	data := c.archetypesComponentsEntities[archetypeId]
	lastKey := len(data) - 1

	data[recordKey] = data[lastKey]
	c.archetypesComponentsEntities[archetypeId] = data[:lastKey]
}

func (c *ComponentsStorage[T]) delete(archetypeId archetypeId, key int) {
	if key < c.size(archetypeId) {
		data := c.archetypesComponentsEntities[archetypeId]
		c.archetypesComponentsEntities[archetypeId] = append(data[:key], data[key+1:]...)
	}
}
