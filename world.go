package volt

import (
	"math/rand"
	"slices"
	"strings"
)

type smallID uint8
type ID uint64

type EntityId ID
type ComponentId smallID
type ArchetypeId ID
type Type []ComponentId

type Archetype struct {
	Id       ArchetypeId
	Type     Type
	entities []EntityId
}

type EntityRecord struct {
	Id          EntityId
	archetypeId ArchetypeId
	key         int
	name        EntityName
}

type EntityName [64]byte
type EntitiesNames map[EntityName]EntityId
type Entities map[EntityId]EntityRecord

type EntityComponentKey struct {
	entityId    EntityId
	componentId ComponentId
}

type World struct {
	ComponentsRegistry ComponentsRegister
	entitiesNames      EntitiesNames
	Entities           Entities
	archetypes         []Archetype
	storage            []storage

	entityAddedFn      func(entityId EntityId)
	entityRemovedFn    func(entityId EntityId)
	componentAddedFn   func(entityId EntityId, componentId ComponentId)
	componentRemovedFn func(entityId EntityId, componentId ComponentId)
}

func CreateWorld(initialCapacity int) *World {
	world := &World{
		entitiesNames:      make(EntitiesNames, initialCapacity),
		Entities:           make(Entities, initialCapacity),
		archetypes:         make([]Archetype, 0, 1024),
		storage:            make([]storage, 256),
		entityAddedFn:      func(entityId EntityId) {},
		entityRemovedFn:    func(entityId EntityId) {},
		componentAddedFn:   func(entityId EntityId, componentId ComponentId) {},
		componentRemovedFn: func(entityId EntityId, componentId ComponentId) {},
	}

	world.createArchetype()

	return world
}

func (world *World) SetEntityAddedFn(entityAddedFn func(entityId EntityId)) {
	world.entityAddedFn = entityAddedFn
}

func (world *World) SetEntityRemovedFn(entityRemovedFn func(entityId EntityId)) {
	world.entityRemovedFn = entityRemovedFn
}

func (world *World) SetComponentAddedFn(componentAddedFn func(entityId EntityId, componentId ComponentId)) {
	world.componentAddedFn = componentAddedFn
}

func (world *World) SetComponentRemovedFn(componentRemovedFn func(entityId EntityId, componentId ComponentId)) {
	world.componentRemovedFn = componentRemovedFn
}

func newEntityId() EntityId {
	return EntityId(rand.Uint64())
}

func (world *World) CreateEntity(name string) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()
	archetype := world.getArchetypeForComponentsIds()

	world.entitiesNames[entityName] = entityId
	entityRecord := EntityRecord{
		Id:   entityId,
		name: entityName,
	}
	world.Entities[entityId] = entityRecord
	world.setArchetype(entityRecord, archetype)

	return entityId
}

func CreateEntityWithComponents2[A, B ComponentInterface](world *World, name string, a A, b B) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents2(world, entityId, a, b)

	return entityId
}

func CreateEntityWithComponents3[A, B, C ComponentInterface](world *World, name string, a A, b B, c C) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents3(world, entityId, a, b, c)

	return entityId
}

func CreateEntityWithComponents4[A, B, C, D ComponentInterface](world *World, name string, a A, b B, c C, d D) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents4(world, entityId, a, b, c, d)

	return entityId
}

func CreateEntityWithComponents5[A, B, C, D, E ComponentInterface](world *World, name string, a A, b B, c C, d D, e E) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents5(world, entityId, a, b, c, d, e)

	return entityId
}

func CreateEntityWithComponents6[A, B, C, D, E, F ComponentInterface](world *World, name string, a A, b B, c C, d D, e E, f F) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents6(world, entityId, a, b, c, d, e, f)

	return entityId
}

func CreateEntityWithComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, name string, a A, b B, c C, d D, e E, f F, g G) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents7(world, entityId, a, b, c, d, e, f, g)

	return entityId
}

func CreateEntityWithComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, name string, a A, b B, c C, d D, e E, f F, g G, h H) EntityId {
	entityName := stringToEntityName(name)
	entityId := newEntityId()

	world.entitiesNames[entityName] = entityId
	world.Entities[entityId] = EntityRecord{name: entityName}

	AddComponents8(world, entityId, a, b, c, d, e, f, g, h)

	return entityId
}

func (world *World) PublishEntity(entityId EntityId) {
	world.entityAddedFn(entityId)
}

func (world *World) RemoveEntity(entityId EntityId) {
	world.entityRemovedFn(entityId)

	entityRecord := world.Entities[entityId]
	archetype := world.archetypes[entityRecord.archetypeId]

	lastEntityKey := len(archetype.entities) - 1
	for _, s := range world.storage {
		if s != nil && slices.Contains(archetype.Type, s.getType()) {
			s.moveLastToKey(archetype.Id, entityRecord.key)
		}
	}

	if lastEntityKey >= 0 {
		lastEntityId := world.archetypes[archetype.Id].entities[lastEntityKey]
		lastEntity := world.Entities[lastEntityId]
		if lastEntity.key > entityRecord.key {
			lastEntity.key = entityRecord.key
			world.Entities[lastEntityId] = lastEntity
			archetype.entities[entityRecord.key] = lastEntityId
		}

		archetype.entities = archetype.entities[:lastEntityKey]
		world.archetypes[archetype.Id] = archetype
	}

	delete(world.entitiesNames, world.Entities[entityId].name)
	delete(world.Entities, entityId)
}

func (world *World) SearchEntity(name string) EntityId {
	entityName := stringToEntityName(name)
	if entityId, ok := world.entitiesNames[entityName]; ok {
		return entityId
	}

	return 0
}

func (world *World) GetEntityName(entityId EntityId) string {
	if entity, ok := world.Entities[entityId]; ok {
		return entityNameToString(entity.name)
	}

	return ""
}

func (world *World) SetEntityName(entityId EntityId, name string) {
	entityName := stringToEntityName(name)

	entityRecord := world.Entities[entityId]
	entityRecord.name = entityName
	world.Entities[entityId] = entityRecord
	world.entitiesNames[entityName] = entityId
}

func stringToEntityName(name string) EntityName {
	var nameByte EntityName
	copy(nameByte[:], name)

	return nameByte
}

func entityNameToString(entityName EntityName) string {
	return strings.TrimRight(string(entityName[:]), "\x00")
}
