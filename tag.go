package volt

import (
	"fmt"
	"slices"
)

const COMPONENTS_INDICES = 0
const TAGS_INDICES = 2048

type TagId = ComponentId

// AddTag adds a TagId to a given EntityId.
// This function returns an error if:
// - The id is lower than the valid range (< TAGS_INDICES)
// - The Tag is already owned
func (world *World) AddTag(tagId TagId, entityId EntityId) error {
	if tagId < TAGS_INDICES {
		return fmt.Errorf("the tagId %d is not allowed, it collides with Components Ids range [%d-%d]", tagId, COMPONENTS_INDICES, TAGS_INDICES)
	}

	if world.HasTag(tagId, entityId) {
		return fmt.Errorf("the entity %d already owns the tag %d", entityId, tagId)
	}

	entityRecord := world.entities[entityId]
	archetype := world.getNextArchetype(entityRecord, tagId)

	oldArchetype := world.getArchetype(entityRecord)
	if archetype.Id != oldArchetype.Id {
		moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
		world.setArchetype(entityRecord, archetype)
	}

	return nil
}

// HasTag returns a boolean, to check if an EntityId owns a Tag.
func (world *World) HasTag(tagId TagId, entityId EntityId) bool {
	if int(entityId) >= len(world.entities) {
		return false
	}
	entityRecord := world.entities[entityId]

	return world.hasComponents(entityRecord, tagId)
}

// RemoveTags removes a Tag for a given EntityId.
// It returns an error if:
// - The entity does not exists.
// - The entity already owns the Tag.
func (world *World) RemoveTag(tagId TagId, entityId EntityId) error {
	if int(entityId) >= len(world.entities) {
		return fmt.Errorf("the entity %d does not exist", entityId)
	}
	entityRecord := world.entities[entityId]

	if !world.HasTag(tagId, entityId) {
		return fmt.Errorf("the entity %d doesn't own the tag %d", entityId, tagId)
	}

	oldArchetype := &world.archetypes[entityRecord.archetypeId]

	// Move every components to the new one, and set all the records
	componentKey := slices.Index(oldArchetype.Type, tagId)

	componentsIds := make([]ComponentId, len(oldArchetype.Type))
	copy(componentsIds, oldArchetype.Type)
	componentsIds = append(componentsIds[:componentKey], componentsIds[componentKey+1:]...)
	archetype := world.getArchetypeForComponentsIds(componentsIds...)
	moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)

	world.setArchetype(entityRecord, archetype)

	return nil
}
