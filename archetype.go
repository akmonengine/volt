package volt

import (
	"slices"
)

func (world *World) createArchetype(componentsIds ...ComponentId) *Archetype {
	archetypeKey := ArchetypeId(len(world.archetypes))
	archetype := Archetype{
		Id:   archetypeKey,
		Type: componentsIds,
	}
	world.archetypes = append(world.archetypes, archetype)

	return &world.archetypes[archetypeKey]
}

func (world *World) getArchetype(entityRecord EntityRecord) *Archetype {
	archetypeId := entityRecord.archetypeId

	if int(archetypeId) >= len(world.archetypes) {
		return nil
	}

	return &world.archetypes[archetypeId]
}

func (world *World) setArchetype(entityRecord EntityRecord, archetype *Archetype) {
	archetype.entities = append(archetype.entities, entityRecord.Id)

	entityRecord.key = len(archetype.entities) - 1
	entityRecord.archetypeId = archetype.Id
	world.Entities[entityRecord.Id] = entityRecord
}

func (world *World) getArchetypeForComponentsIds(componentsIds ...ComponentId) *Archetype {
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

func (world *World) getArchetypesForComponentsIds(componentsIds ...ComponentId) []Archetype {
	var archetypes []Archetype

	for _, archetype := range world.archetypes {
		i := 0
		for _, componentId := range componentsIds {
			if slices.Contains(archetype.Type, componentId) {
				i++
			}
		}

		if i == len(componentsIds) {
			archetypes = append(archetypes, archetype)
		}
	}

	return archetypes
}

func (world *World) getNextArchetype(entityId EntityId, componentsIds ...ComponentId) *Archetype {
	var archetype *Archetype
	if entityRecord, ok := world.Entities[entityId]; !ok {
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
