// Package volt is an ECS for game development, based on the Archetype paradigm.
package volt

import (
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
}

type entities []entityRecord

// World representation, container of all the data related to entities and their Components.
type World struct {
	componentsRegistry ComponentsRegister
	pool               pool
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
		pool:               pool{},
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
func (world *World) CreateEntity() EntityId {
	entityId := world.pool.Get()
	archetype := world.getArchetypeForComponentsIds()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)
	world.setArchetype(entityRecord, archetype)

	return entityId
}

func (world *World) addEntity(entityRecord entityRecord) {
	if int(entityRecord.Id) < len(world.entities) {
		world.entities[entityRecord.Id] = entityRecord
	} else {
		world.entities = append(world.entities, entityRecord)
	}
}

// CreateEntityWithComponents2 creates an entity in World;
// It sets the components A, B to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents2[A, B ComponentInterface](world *World, a A, b B) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

	err := addComponents2(world, entityRecord, a, b)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents3 creates an entity in World;
//
// It sets the components A, B, C to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents3[A, B, C ComponentInterface](world *World, a A, b B, c C) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

	err := addComponents3(world, entityRecord, a, b, c)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents4 creates an entity in World;
//
// It sets the components A, B, C, D to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents4[A, B, C, D ComponentInterface](world *World, a A, b B, c C, d D) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

	err := addComponents4(world, entityRecord, a, b, c, d)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents5 creates an entity in World;
//
// It sets the components A, B, C, D, E to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents5[A, B, C, D, E ComponentInterface](world *World, a A, b B, c C, d D, e E) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

	err := addComponents5(world, entityRecord, a, b, c, d, e)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents6 creates an entity in World;
//
// It sets the components A, B, C, D, E, F to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents6[A, B, C, D, E, F ComponentInterface](world *World, a A, b B, c C, d D, e E, f F) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

	err := addComponents6(world, entityRecord, a, b, c, d, e, f)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents7 creates an entity in World;
//
// It sets the components A, B, C, D, E, F, G to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, a A, b B, c C, d D, e E, f F, g G) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

	err := addComponents7(world, entityRecord, a, b, c, d, e, f, g)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents8 creates an entity in World;
//
// It sets the components A, B, C, D, E, F, G, H to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, a A, b B, c C, d D, e E, f F, g G, h H) (EntityId, error) {
	entityId := world.pool.Get()

	entityRecord := entityRecord{Id: entityId}
	world.addEntity(entityRecord)

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

	world.pool.Recycle(entityId)
}

// Count returns the number of entities in World.
func (world *World) Count() int {
	return len(world.entities) - world.pool.Count()
}
