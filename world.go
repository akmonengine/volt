// Package volt is an ECS for game development, based on the Archetype paradigm.
package volt

// uint16 identifier, for small scoped data.
type smallId uint16

// uint64 identifier, for big scoped data.
type id uint64

// Entity identifier in the world: a handle packing the entity's slot (index,
// low 32 bits) and the generation of that slot (high 32 bits). Generations
// start at 1 and are bumped every time the slot is freed by RemoveEntity, so a
// handle kept after its entity was removed never aliases the slot's next
// occupant: it is simply dead. The zero value is the null handle, which never
// refers to a live entity.
type EntityId id

const (
	entityIndexBits = 32
	entityIndexMask = EntityId(1)<<entityIndexBits - 1
)

// Index returns the slot of the entity in the World.
func (entityId EntityId) Index() uint32 {
	return uint32(entityId & entityIndexMask)
}

// Generation returns the generation of the handle. Live handles have a
// generation of at least 1; the null handle has generation 0.
func (entityId EntityId) Generation() uint32 {
	return uint32(entityId >> entityIndexBits)
}

// IsNull reports whether entityId is the null handle, which never refers to a
// live entity. It is the zero value of EntityId, so an unset field is null.
func (entityId EntityId) IsNull() bool {
	return entityId == 0
}

func newEntityId(index uint32, generation uint32) EntityId {
	return EntityId(generation)<<entityIndexBits | EntityId(index)
}

// nextGeneration bumps a slot generation, skipping 0 on wrap-around so that the
// null handle can never be minted.
func nextGeneration(generation uint32) uint32 {
	generation++
	if generation == 0 {
		generation = 1
	}

	return generation
}

// Component identifier in the register.
type ComponentId smallId

// archetype identifier in the world.
type archetypeId id

// List of ComponentId.
type componentsIds []ComponentId

// Implementation of an archetype with its identifier, componentsIds, and entitiesIds
type archetype struct {
	Id       archetypeId
	Type     componentsIds // sorted set of the components (and tags) held
	entities []EntityId

	// Archetype graph: cached transitions to neighbour archetypes.
	// addEdges[c] is the archetype reached by adding component c to this one;
	// removeEdges[c] the one reached by removing c. Archetypes are never
	// destroyed, so these edges never go stale. They turn the per-operation
	// archetype lookup from a linear scan into an O(1) hop after the first time.
	addEdges    map[ComponentId]archetypeId
	removeEdges map[ComponentId]archetypeId
}

// Record of an entity slot: the handle currently bound to the slot (its
// generation tells recycled slots apart), the archetype the entity lives in
// and its key (row) in that archetype. A negative key means the slot holds no
// live entity: either freed, or allocated but not yet placed in an archetype.
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
	archetypesByKey    map[uint64]archetypeId // archetype of a set of components, by archetypeKey
	transitions        [transitionCacheSize]transition
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
		entities:           make(entities, 0, initialCapacity),
		archetypes:         make([]archetype, 0, 1024),
		archetypesByKey:    make(map[uint64]archetypeId, 1024),
		storage:            make([]storage, TAGS_INDICES),
		entityAddedFn:      func(entityId EntityId) {},
		entityRemovedFn:    func(entityId EntityId) {},
		componentAddedFn:   func(entityId EntityId, componentId ComponentId) {},
		componentRemovedFn: func(entityId EntityId, componentId ComponentId) {},
	}

	world.createArchetype(nil)

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
	entityId := world.newEntity()
	archetype := world.getArchetypeForComponentsIds()

	world.setArchetype(world.entities[entityId.Index()], archetype)

	return entityId
}

// newEntity allocates a slot and returns its live handle. A brand new slot
// starts at generation 1; a recycled slot carries the generation bumped when it
// was freed. The record is reset and left unplaced (key -1) until setArchetype
// puts the entity in an archetype.
func (world *World) newEntity() EntityId {
	index, recycled := world.pool.Get()

	if recycled {
		entityId := world.entities[index].Id
		world.entities[index] = entityRecord{Id: entityId, key: -1}

		return entityId
	}

	entityId := newEntityId(index, 1)
	world.entities = append(world.entities, entityRecord{Id: entityId, key: -1})

	return entityId
}

// discardEntity gives back a slot allocated by newEntity whose entity could not
// be placed (creation failed). The handle was never handed out, so the
// generation is kept as is.
func (world *World) discardEntity(entityId EntityId) {
	world.pool.Recycle(entityId.Index())
}

// CreateEntityWithComponents2 creates an entity in World;
// It sets the components A, B to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents2[A, B ComponentInterface](world *World, a A, b B) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents2(world, world.entities[entityId.Index()], a, b)
	if err != nil {
		world.discardEntity(entityId)

		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents3 creates an entity in World;
//
// It sets the components A, B, C to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents3[A, B, C ComponentInterface](world *World, a A, b B, c C) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents3(world, world.entities[entityId.Index()], a, b, c)
	if err != nil {
		world.discardEntity(entityId)

		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents4 creates an entity in World;
//
// It sets the components A, B, C, D to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents4[A, B, C, D ComponentInterface](world *World, a A, b B, c C, d D) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents4(world, world.entities[entityId.Index()], a, b, c, d)
	if err != nil {
		world.discardEntity(entityId)

		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents5 creates an entity in World;
//
// It sets the components A, B, C, D, E to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents5[A, B, C, D, E ComponentInterface](world *World, a A, b B, c C, d D, e E) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents5(world, world.entities[entityId.Index()], a, b, c, d, e)
	if err != nil {
		world.discardEntity(entityId)

		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents6 creates an entity in World;
//
// It sets the components A, B, C, D, E, F to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents6[A, B, C, D, E, F ComponentInterface](world *World, a A, b B, c C, d D, e E, f F) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents6(world, world.entities[entityId.Index()], a, b, c, d, e, f)
	if err != nil {
		world.discardEntity(entityId)

		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents7 creates an entity in World;
//
// It sets the components A, B, C, D, E, F, G to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, a A, b B, c C, d D, e E, f F, g G) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents7(world, world.entities[entityId.Index()], a, b, c, d, e, f, g)
	if err != nil {
		world.discardEntity(entityId)

		return 0, err
	}

	return entityId, nil
}

// CreateEntityWithComponents8 creates an entity in World;
//
// It sets the components A, B, C, D, E, F, G, H to the entity, for faster performances than the atomic version.
func CreateEntityWithComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, a A, b B, c C, d D, e E, f F, g G, h H) (EntityId, error) {
	entityId := world.newEntity()

	err := addComponents8(world, world.entities[entityId.Index()], a, b, c, d, e, f, g, h)
	if err != nil {
		world.discardEntity(entityId)

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
//
// The slot is freed for reuse and its generation bumped: every handle to the
// removed entity is dead from now on (Exists reports false), even once the slot
// is recycled by a new entity. Removing a dead handle is a no-op.
func (world *World) RemoveEntity(entityId EntityId) {
	if !world.Exists(entityId) {
		return
	}

	world.entityRemovedFn(entityId)

	index := entityId.Index()
	record := world.entities[index]
	archetype := world.archetypes[record.archetypeId]

	lastEntityKey := len(archetype.entities) - 1
	for _, componentId := range archetype.Type {
		// Tags have no storage: their id lives outside the storage range,
		// so indexing world.storage[componentId] would overflow.
		if componentId >= TAGS_INDICES {
			continue
		}
		s := world.storage[componentId]
		if s != nil {
			s.moveLastToKey(archetype.Id, record.key)
		}
	}

	if lastEntityKey >= 0 {
		lastEntityId := world.archetypes[archetype.Id].entities[lastEntityKey]
		lastEntity := world.entities[lastEntityId.Index()]
		if lastEntity.key > record.key {
			lastEntity.key = record.key
			world.entities[lastEntityId.Index()] = lastEntity
			archetype.entities[record.key] = lastEntityId
		}

		archetype.entities = archetype.entities[:lastEntityKey]
		world.archetypes[archetype.Id] = archetype
	}

	// Free the slot: the tombstone (negative key) and the bumped generation
	// together guarantee that no handle, kept or forged, resolves to it until a
	// new entity is created there with the new generation.
	world.entities[index] = entityRecord{
		Id:  newEntityId(index, nextGeneration(entityId.Generation())),
		key: -1,
	}
	world.pool.Recycle(index)
}

// Exists reports whether entityId refers to a live entity of the World.
//
// It returns false for the null handle, for slots never opened, and for dead
// handles: entities removed, whether or not their slot has since been recycled
// by a new entity (the generation tells them apart).
func (world *World) Exists(entityId EntityId) bool {
	index := entityId.Index()
	if int(index) >= len(world.entities) {
		return false
	}

	record := &world.entities[index]

	return record.Id == entityId && record.key >= 0
}

// Count returns the number of entities in World.
func (world *World) Count() int {
	return len(world.entities) - world.pool.Count()
}
