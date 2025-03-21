// Package volt is an ECS for game development, based on the Archetype paradigm.
package volt

import (
	"hash/fnv"
	"slices"
)

// uint16 identifier, for small scoped data.
type smallId uint16

// uint64 identifier, for big scoped data.
type id uint64

// Entity identifier in the world.
type EntityId id

// Component identifier in the register.
type ComponentId smallId

// archetype identifier in the world.
type archetypeId id

// List of ComponentId.
type componentsIds []ComponentId

// Implementation of an archetype with its identifier, componentsIds, and entitiesIds
type archetype struct {
	Id       archetypeId
	Type     componentsIds
	entities []EntityId
}

// Container of archetype and key position in storage, for a given EntityId
type entityRecord struct {
	Id          EntityId
	archetypeId archetypeId
	key         int
	name        entityName
}

// entityName is a string transformed to byte array.
//
// It avoids the garbage collector to analyze this data constantly,
// at the price of a fixed data size.
type entityName = string
type entities map[EntityId]entityRecord

// World representation, container of all the data related to entities and their Components.
type World struct {
	componentsRegistry ComponentsRegister
	entities           entities
	archetypes         []archetype
	storage            []storage

	entityAddedFn      func(entityId EntityId)
	entityRemovedFn    func(entityId EntityId)
	componentAddedFn   func(entityId EntityId, componentId ComponentId)
	componentRemovedFn func(entityId EntityId, componentId ComponentId)
}

// CreateWorld returns a pointer to a new World.
//
// It preallocates initialCapacity in memory.
func CreateWorld(initialCapacity int) *World {
	world := &World{
		entities:           make(entities, initialCapacity),
		archetypes:         make([]archetype, 0, 1024),
		storage:            make([]storage, TAGS_INDICES),
		entityAddedFn:      func(entityId EntityId) {},
		entityRemovedFn:    func(entityId EntityId) {},
		componentAddedFn:   func(entityId EntityId, componentId ComponentId) {},
		componentRemovedFn: func(entityId EntityId, componentId ComponentId) {},
	}

	world.createArchetype()

	return world
}

// SetEntityAddedFn sets a callback for when a new entity is added.
func (world *World) SetEntityAddedFn(entityAddedFn func(entityId EntityId)) {
	world.entityAddedFn = entityAddedFn
}

// SetEntityRemovedFn sets a callback for when an entity is removed.
func (world *World) SetEntityRemovedFn(entityRemovedFn func(entityId EntityId)) {
	world.entityRemovedFn = entityRemovedFn
}

// SetComponentAddedFn sets a callback for when a component is added to an entity.
func (world *World) SetComponentAddedFn(componentAddedFn func(entityId EntityId, componentId ComponentId)) {
	world.componentAddedFn = componentAddedFn
}

// SetComponentRemovedFn sets a callback for when a component is removed.
func (world *World) SetComponentRemovedFn(componentRemovedFn func(entityId EntityId, componentId ComponentId)) {
	world.componentRemovedFn = componentRemovedFn
}

// CreateEntity creates a new Entity in World;
// It is linked to no Component.
func (world *World) CreateEntity(name string) EntityId {
	if existingId := world.SearchEntity(name); existingId != 0 {
		return existingId
	}

	entityId := hashEntityName(name)
	archetype := world.getArchetypeForComponentsIds()

	entityRecord := entityRecord{
		Id:   entityId,
		name: name,
	}
	world.entities[entityId] = entityRecord
	world.setArchetype(entityRecord, archetype)

	return entityId
}

// CreateEntityWithComponents2 creates an entity in World;
// It sets the components A, B to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents2[A, B ComponentInterface](world *World, name string, a A, b B) (EntityId, error) {
	entityId := hashEntityName(name)

	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents2(world, entityRecord, a, b)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents3 creates an entity in World;
//
// It sets the components A, B, C to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents3[A, B, C ComponentInterface](world *World, name string, a A, b B, c C) (EntityId, error) {
	entityId := hashEntityName(name)

	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents3(world, entityRecord, a, b, c)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents4 creates an entity in World;
//
// It sets the components A, B, C, D to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents4[A, B, C, D ComponentInterface](world *World, name string, a A, b B, c C, d D) (EntityId, error) {
	entityId := hashEntityName(name)

	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents4(world, entityRecord, a, b, c, d)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents5 creates an entity in World;
//
// It sets the components A, B, C, D, E to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents5[A, B, C, D, E ComponentInterface](world *World, name string, a A, b B, c C, d D, e E) (EntityId, error) {
	entityId := hashEntityName(name)

	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents5(world, entityRecord, a, b, c, d, e)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents6 creates an entity in World;
//
// It sets the components A, B, C, D, E, F to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents6[A, B, C, D, E, F ComponentInterface](world *World, name string, a A, b B, c C, d D, e E, f F) (EntityId, error) {
	entityId := hashEntityName(name)
	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents6(world, entityRecord, a, b, c, d, e, f)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents7 creates an entity in World;
//
// It sets the components A, B, C, D, E, F, G to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, name string, a A, b B, c C, d D, e E, f F, g G) (EntityId, error) {
	entityId := hashEntityName(name)
	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents7(world, entityRecord, a, b, c, d, e, f, g)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents8 creates an entity in World;
//
// It sets the components A, B, C, D, E, F, G, H to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, name string, a A, b B, c C, d D, e E, f F, g G, h H) (EntityId, error) {
	entityId := hashEntityName(name)
	entityRecord := entityRecord{Id: entityId, name: name}
	world.entities[entityId] = entityRecord

	err := addComponents8(world, entityRecord, a, b, c, d, e, f, g, h)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// PublishEntity calls the callback setted in SetEntityAddedFn.
func (world *World) PublishEntity(entityId EntityId) {
	world.entityAddedFn(entityId)
}

// RemoveEntity removes all the data related to an Entity.
//
// It calls the callback setted in SetEntityRemovedFn beforehand, so that the callback still has access to the data.
func (world *World) RemoveEntity(entityId EntityId) {
	world.entityRemovedFn(entityId)

	entityRecord := world.entities[entityId]
	archetype := world.archetypes[entityRecord.archetypeId]

	lastEntityKey := len(archetype.entities) - 1
	for _, componentId := range archetype.Type {
		s := world.storage[componentId]
		if s != nil && slices.Contains(archetype.Type, s.getType()) {
			s.moveLastToKey(archetype.Id, entityRecord.key)
		}
	}

	if lastEntityKey >= 0 {
		lastEntityId := world.archetypes[archetype.Id].entities[lastEntityKey]
		lastEntity := world.entities[lastEntityId]
		if lastEntity.key > entityRecord.key {
			lastEntity.key = entityRecord.key
			world.entities[lastEntityId] = lastEntity
			archetype.entities[entityRecord.key] = lastEntityId
		}

		archetype.entities = archetype.entities[:lastEntityKey]
		world.archetypes[archetype.Id] = archetype
	}

	delete(world.entities, entityId)
}

// SearchEntity returns the EntityId named by name.
// If not found, returns 0.
func (world *World) SearchEntity(name string) EntityId {
	entityId := hashEntityName(name)
	if _, ok := world.entities[entityId]; ok {
		return entityId
	}

	return 0
}

// GetEntityName returns the name of an EntityId.
// If not found, returns an empty string.
func (world *World) GetEntityName(entityId EntityId) string {
	if entity, ok := world.entities[entityId]; ok {
		return entity.name
	}

	return ""
}

// SetEntityName sets the name for an EntityId.
func (world *World) SetEntityName(entityId EntityId, name string) {
	entityRecord := world.entities[entityId]
	entityRecord.name = name
	world.entities[entityId] = entityRecord
}

// Count returns the number of entities in World.
func (world *World) Count() int {
	return len(world.entities)
}

func hashEntityName(name entityName) EntityId {
	h := fnv.New64()
	_, err := h.Write([]byte(name))
	if err != nil {
		return EntityId(0)
	}

	return EntityId(h.Sum64())
}
