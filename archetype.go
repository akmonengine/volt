package volt

import (
	"slices"
)

func (world *World) createArchetype(componentsIds ...ComponentId) *archetype {
	archetypeKey := archetypeId(len(world.archetypes))
	archetype := archetype{
		Id:   archetypeKey,
		Type: componentsIds,
	}
	world.archetypes = append(world.archetypes, archetype)

	return &world.archetypes[archetypeKey]
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
	world.entities[entityRecord.Id] = entityRecord
}

func (world *World) getArchetypeForComponentsIds(componentsIds ...ComponentId) *archetype {
	for i, archetype := range world.archetypes {
		if len(archetype.Type) != len(componentsIds) {
			continue
		}

		count := 0
		for _, componentId := range componentsIds {
			if slices.Contains(archetype.Type, componentId) {
				count++
			} else {
				break
			}
		}

		if count == len(archetype.Type) {
			return &world.archetypes[i]
		}
	}

	return world.createArchetype(componentsIds...)
}

func (world *World) getNextArchetype(entityRecord entityRecord, componentsIds ...ComponentId) *archetype {
	// Fast path: a single-component transition (AddComponent, AddTag, ...) is
	// resolved through the archetype graph, avoiding both the linear scan over
	// all archetypes and the slice rebuild done below.
	if len(componentsIds) == 1 {
		return world.archetypeAfterAdd(entityRecord.archetypeId, componentsIds[0])
	}

	var archetype *archetype
	if entityRecord.archetypeId == 0 {
		archetype = world.getArchetypeForComponentsIds(componentsIds...)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if oldArchetype != nil {
			archetype = world.getArchetypeForComponentsIds(append(componentsIds, oldArchetype.Type...)...)
		} else {
			archetype = world.getArchetypeForComponentsIds(componentsIds...)
		}
	}

	return archetype
}

// archetypeAfterAdd returns the archetype obtained by adding componentId to the
// archetype fromId, using (and lazily populating) the cached archetype graph.
func (world *World) archetypeAfterAdd(fromId archetypeId, componentId ComponentId) *archetype {
	if destId, ok := world.archetypes[fromId].addEdges[componentId]; ok {
		return &world.archetypes[destId]
	}

	// Cache miss: compute the destination once. getArchetypeForComponentsIds may
	// create a new archetype and reallocate world.archetypes, so we resolve every
	// archetype by index afterwards rather than holding a stale pointer.
	newType := append(slices.Clone(world.archetypes[fromId].Type), componentId)
	destId := world.getArchetypeForComponentsIds(newType...).Id
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
	destId := world.getArchetypeForComponentsIds(newType...).Id
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
