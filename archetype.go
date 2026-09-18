package volt

import (
	"slices"
)

// FNV-1a 64-bit parameters, used to hash a set of component ids into an
// archetype key.
const (
	fnv64Offset uint64 = 14695981039346656037
	fnv64Prime  uint64 = 1099511628211
)

// maxInlineComponents bounds the stack scratch used to sort a set of component
// ids before looking its archetype up. Larger sets fall back to a heap copy.
const maxInlineComponents = 16

// maxTransitionComponents is the largest number of components added in one
// call (AddComponents8): the ids of a transition fit inline in its cache entry.
const maxTransitionComponents = 8

// transitionCacheSize is the number of slots of the transition cache, a power
// of two so that a slot is a mask of the key.
const transitionCacheSize = 64

// transition remembers where adding a list of components to an archetype leads.
// The ids are kept as given by the caller, so a hit needs no sorting: two orders
// of the same set are two entries pointing at the same archetype. Archetypes are
// never destroyed, so an entry never goes stale.
type transition struct {
	from archetypeId
	ids  [maxTransitionComponents]ComponentId
	n    uint8
	dest archetypeId
}

// transitionSlot maps a transition to its slot in the cache.
func transitionSlot(fromId archetypeId, componentsIds []ComponentId) int {
	hash := (fnv64Offset ^ uint64(fromId)) * fnv64Prime
	for _, componentId := range componentsIds {
		hash ^= uint64(componentId)
		hash *= fnv64Prime
	}

	return int(hash & (transitionCacheSize - 1))
}

// archetypeKey hashes a sorted set of component ids. Two different sets may
// share a key: a lookup always confirms the archetype Type before trusting it.
func archetypeKey(sorted []ComponentId) uint64 {
	hash := fnv64Offset
	for _, componentId := range sorted {
		hash ^= uint64(componentId)
		hash *= fnv64Prime
	}

	return hash
}

// createArchetype registers a new archetype holding exactly the given set of
// components. The slice must be sorted; the archetype owns it from now on.
func (world *World) createArchetype(componentsIds componentsIds) *archetype {
	id := archetypeId(len(world.archetypes))
	world.archetypes = append(world.archetypes, archetype{
		Id:   id,
		Type: componentsIds,
	})

	// The first archetype registered under a key owns it; a later archetype
	// whose set collides on the same key is found by scan instead.
	key := archetypeKey(componentsIds)
	if _, taken := world.archetypesByKey[key]; !taken {
		world.archetypesByKey[key] = id
	}

	return &world.archetypes[id]
}

func (world *World) getArchetype(entityRecord entityRecord) *archetype {
	archetypeId := entityRecord.archetypeId

	if int(archetypeId) >= len(world.archetypes) {
		return nil
	}

	return &world.archetypes[archetypeId]
}

func (world *World) setArchetype(entityRecord entityRecord, archetype *archetype) {
	archetype.entities = append(archetype.entities, entityRecord.Id)

	entityRecord.key = len(archetype.entities) - 1
	entityRecord.archetypeId = archetype.Id
	world.entities[entityRecord.Id.Index()] = entityRecord
}

// getArchetypeForComponentsIds returns the archetype holding exactly the given
// set of components, whatever their order, creating it if needed.
func (world *World) getArchetypeForComponentsIds(componentsIds ...ComponentId) *archetype {
	var scratch [maxInlineComponents]ComponentId

	return world.archetypeForSet(append(scratch[:0], componentsIds...))
}

// archetypeForSet resolves the archetype holding exactly the given set of
// components, creating it if needed. The set is a scratch copy: it is sorted in
// place and never retained, the archetype created on a miss owns its own copy.
func (world *World) archetypeForSet(set []ComponentId) *archetype {
	slices.Sort(set)

	id, found := world.archetypesByKey[archetypeKey(set)]
	if !found {
		return world.createArchetype(slices.Clone(set))
	}
	if slices.Equal(world.archetypes[id].Type, set) {
		return &world.archetypes[id]
	}

	// Key collision: another set owns the key, scan for this one.
	for i := range world.archetypes {
		if slices.Equal(world.archetypes[i].Type, set) {
			return &world.archetypes[i]
		}
	}

	return world.createArchetype(slices.Clone(set))
}

// getNextArchetype returns the archetype reached by adding componentsIds to the
// archetype the entity lives in.
func (world *World) getNextArchetype(entityRecord entityRecord, componentsIds ...ComponentId) *archetype {
	// A single-component transition is resolved through the archetype graph.
	if len(componentsIds) == 1 {
		return world.archetypeAfterAdd(entityRecord.archetypeId, componentsIds[0])
	}
	if len(componentsIds) <= maxTransitionComponents {
		return world.archetypeAfterAddN(entityRecord.archetypeId, componentsIds)
	}

	return world.archetypeAfterAddSet(entityRecord.archetypeId, componentsIds)
}

// archetypeAfterAddN resolves a multi-component transition through the
// transition cache. A hit costs one key comparison; a miss resolves the set and
// fills the slot, which the next transition hashing to it overwrites. The whole
// key is compared, so a slot collision is a miss, never a wrong archetype.
func (world *World) archetypeAfterAddN(fromId archetypeId, componentsIds []ComponentId) *archetype {
	var key [maxTransitionComponents]ComponentId
	copy(key[:], componentsIds)
	n := uint8(len(componentsIds))

	t := &world.transitions[transitionSlot(fromId, componentsIds)]
	if t.from == fromId && t.n == n && t.ids == key {
		return &world.archetypes[t.dest]
	}

	dest := world.archetypeAfterAddSet(fromId, componentsIds)
	*t = transition{from: fromId, ids: key, n: n, dest: dest.Id}

	return dest
}

// archetypeAfterAddSet resolves the archetype holding the components of the
// archetype fromId plus componentsIds, through the canonical key.
func (world *World) archetypeAfterAddSet(fromId archetypeId, componentsIds []ComponentId) *archetype {
	var scratch [maxInlineComponents]ComponentId
	set := append(scratch[:0], componentsIds...)
	if int(fromId) < len(world.archetypes) {
		set = append(set, world.archetypes[fromId].Type...)
	}

	return world.archetypeForSet(set)
}

// archetypeAfterAdd returns the archetype obtained by adding componentId to the
// archetype fromId, using (and lazily populating) the cached archetype graph.
func (world *World) archetypeAfterAdd(fromId archetypeId, componentId ComponentId) *archetype {
	if destId, ok := world.archetypes[fromId].addEdges[componentId]; ok {
		return &world.archetypes[destId]
	}

	// Cache miss: compute the destination once. archetypeForSet may create a new
	// archetype and reallocate world.archetypes, so we resolve every archetype
	// by index afterwards rather than holding a stale pointer.
	newType := append(slices.Clone(world.archetypes[fromId].Type), componentId)
	destId := world.archetypeForSet(newType).Id
	world.linkArchetypes(fromId, destId, componentId)

	return &world.archetypes[destId]
}

// archetypeAfterRemove returns the archetype obtained by removing componentId
// from the archetype fromId, using (and lazily populating) the archetype graph.
func (world *World) archetypeAfterRemove(fromId archetypeId, componentId ComponentId) *archetype {
	if destId, ok := world.archetypes[fromId].removeEdges[componentId]; ok {
		return &world.archetypes[destId]
	}

	fromType := world.archetypes[fromId].Type
	newType := make(componentsIds, 0, len(fromType))
	for _, c := range fromType {
		if c != componentId {
			newType = append(newType, c)
		}
	}
	destId := world.archetypeForSet(newType).Id
	// dest --add componentId--> from, and from --remove componentId--> dest.
	world.linkArchetypes(destId, fromId, componentId)

	return &world.archetypes[destId]
}

// linkArchetypes records the bidirectional transition between two archetypes:
// fromId --add componentId--> destId and destId --remove componentId--> fromId.
func (world *World) linkArchetypes(fromId, destId archetypeId, componentId ComponentId) {
	from := &world.archetypes[fromId]
	if from.addEdges == nil {
		from.addEdges = make(map[ComponentId]archetypeId)
	}
	from.addEdges[componentId] = destId

	dest := &world.archetypes[destId]
	if dest.removeEdges == nil {
		dest.removeEdges = make(map[ComponentId]archetypeId)
	}
	dest.removeEdges[componentId] = fromId
}

// matchArchetypes appends, into buf, the id of every archetype whose Type
// contains all of componentsIds (the query's required components + tags). The
// caller passes a reused buffer (buf[:0]) to avoid per-call allocations.
func (world *World) matchArchetypes(buf []archetypeId, componentsIds []ComponentId) []archetypeId {
	for i := range world.archetypes {
		archetype := &world.archetypes[i]

		matched := true
		for _, componentId := range componentsIds {
			if !slices.Contains(archetype.Type, componentId) {
				matched = false
				break
			}
		}

		if matched {
			buf = append(buf, archetypeId(i))
		}
	}

	return buf
}
