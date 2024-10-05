package volt

import (
	"fmt"
	"slices"
)

type ComponentInterface interface {
	GetComponentId() ComponentId
}

func (world *World) getComponentsIds(components ...ComponentInterface) []ComponentId {
	componentsIds := make([]ComponentId, len(components))

	for i, component := range components {
		componentsIds[i] = component.GetComponentId()
	}

	return componentsIds
}

func ConfigureComponent[T ComponentInterface](world *World, conf any) T {
	var t T
	componentRegistry := world.ComponentsRegistry[t.GetComponentId()]

	componentRegistry.builderFn(&t, conf)

	return t
}

func AddComponent[T ComponentInterface](world *World, entityId EntityId, component T) error {
	componentId := component.GetComponentId()
	if world.HasComponents(entityId, componentId) {
		return fmt.Errorf("the entity %d already owns the component %d", entityId, componentId)
	}

	archetype := world.getNextArchetype(entityId, world.getComponentsIds(component)...)
	err := addComponentsToArchetype1(world, entityId, archetype, component)
	if err != nil {
		return fmt.Errorf("the component %d cannot be added to entity %d: %w", componentId, entityId, err)
	}

	world.componentAddedFn(entityId, componentId)

	return nil
}

func AddComponents2[A, B ComponentInterface](world *World, entityId EntityId, a A, b B) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId()})
	}

	err := addComponentsToArchetype2(world, entityRecord, archetype, a, b)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())

	return nil
}

func AddComponents3[A, B, C ComponentInterface](world *World, entityId EntityId, a A, b B, c C) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId(), c.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId()})
	}

	err := addComponentsToArchetype3(world, entityRecord, archetype, a, b, c)
	if err != nil {
		return fmt.Errorf("the components %v cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())

	return nil
}

func AddComponents4[A, B, C, D ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId()})
	}

	err := addComponentsToArchetype4(world, entityRecord, archetype, a, b, c, d)
	if err != nil {
		return fmt.Errorf("the components %v cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())

	return nil
}

func AddComponents5[A, B, C, D, E ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId()})
	}

	err := addComponentsToArchetype5(world, entityRecord, archetype, a, b, c, d, e)
	if err != nil {
		return fmt.Errorf("the components %v cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())
	world.componentAddedFn(entityId, e.GetComponentId())

	return nil
}

func AddComponents6[A, B, C, D, E, F ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E, f F) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId()})
	}

	err := addComponentsToArchetype6(world, entityRecord, archetype, a, b, c, d, e, f)
	if err != nil {
		return fmt.Errorf("the components %v cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())
	world.componentAddedFn(entityId, e.GetComponentId())
	world.componentAddedFn(entityId, f.GetComponentId())

	return nil
}

func AddComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E, f F, g G) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId()})
	}

	err := addComponentsToArchetype7(world, entityRecord, archetype, a, b, c, d, e, f, g)
	if err != nil {
		return fmt.Errorf("the components %v cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())
	world.componentAddedFn(entityId, e.GetComponentId())
	world.componentAddedFn(entityId, f.GetComponentId())
	world.componentAddedFn(entityId, g.GetComponentId())

	return nil
}

func AddComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E, f F, g G, h H) error {
	archetype := world.getArchetypeForComponentsIds(a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId())

	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}
	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId()})
	}

	err := addComponentsToArchetype8(world, entityRecord, archetype, a, b, c, d, e, f, g, h)
	if err != nil {
		return fmt.Errorf("the components %v cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())
	world.componentAddedFn(entityId, e.GetComponentId())
	world.componentAddedFn(entityId, f.GetComponentId())
	world.componentAddedFn(entityId, g.GetComponentId())
	world.componentAddedFn(entityId, h.GetComponentId())

	return nil
}

func (world *World) AddComponent(entityId EntityId, componentId ComponentId, conf any) error {
	if world.HasComponents(entityId, componentId) {
		return fmt.Errorf("the entity %d already owns the component %d", entityId, componentId)
	}

	componentRegistry, err := world.getConfigByComponentId(componentId)
	if err != nil {
		return err
	}
	err = componentRegistry.addComponent(world, entityId, conf)
	if err != nil {
		return err
	}

	world.componentAddedFn(entityId, componentId)

	return nil
}

func RemoveComponent[T ComponentInterface](world *World, entityId EntityId) error {
	var t T
	componentId := t.GetComponentId()
	entityRecord := world.Entities[entityId]

	if !world.hasComponents(entityRecord, componentId) {
		return fmt.Errorf("the entity %d doesn't own the component %d", entityId, componentId)
	}

	// Remove from the previous archetype
	s := getStorage[T](world)
	removeComponent(world, s, entityRecord, componentId)

	return nil
}

// Useful when generics are not available, but slower than the generic method
func (world *World) RemoveComponent(entityId EntityId, componentId ComponentId) error {
	entityRecord := world.Entities[entityId]

	if !world.hasComponents(entityRecord, componentId) {
		return fmt.Errorf("the entity %d doesn't own the component %d", entityId, componentId)
	}

	// Remove from the previous archetype
	s, err := world.getStorageForComponentId(componentId)
	if err != nil {
		return err
	}

	removeComponent(world, s, entityRecord, componentId)

	return nil
}

func removeComponent(world *World, s storage, entityRecord EntityRecord, componentId ComponentId) {
	world.componentRemovedFn(entityRecord.Id, componentId)

	oldArchetype := &world.archetypes[entityRecord.archetypeId]

	s.moveLastToKey(oldArchetype.Id, entityRecord.key)

	// Move every components to the new one, and set all the records
	componentKey := slices.Index(oldArchetype.Type, componentId)

	componentsIds := make([]ComponentId, len(oldArchetype.Type))
	copy(componentsIds, oldArchetype.Type)
	componentsIds = append(componentsIds[:componentKey], componentsIds[componentKey+1:]...)
	archetype := world.getArchetypeForComponentsIds(componentsIds...)
	moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)

	world.setArchetype(entityRecord, archetype)
}

func (world *World) HasComponents(entityId EntityId, componentsIds ...ComponentId) bool {
	entityRecord, ok := world.Entities[entityId]
	if !ok {
		return false
	}

	return world.hasComponents(entityRecord, componentsIds...)
}

func (world *World) hasComponents(entityRecord EntityRecord, componentsIds ...ComponentId) bool {
	archetype := world.archetypes[entityRecord.archetypeId]
	for _, componentId := range componentsIds {
		if !slices.Contains(archetype.Type, componentId) {
			return false
		}
	}

	return true
}

func GetComponent[T ComponentInterface](world *World, entityId EntityId) *T {
	s := getStorage[T](world)
	entityRecord := world.Entities[entityId]

	if !s.hasArchetype(entityRecord.archetypeId) {
		return nil
	}

	return &s.archetypesComponentsEntities[entityRecord.archetypeId][entityRecord.key]
}

// Useful when generics are not available, but slower than the generic method
func (world *World) GetComponent(entityId EntityId, componentId ComponentId) (any, error) {
	entityRecord := world.Entities[entityId]
	s, err := world.getStorageForComponentId(componentId)
	if err != nil {
		return nil, err
	}

	if !s.hasArchetype(entityRecord.archetypeId) {
		return nil, fmt.Errorf("the entity %d doesn't own the component %d", entityId, componentId)
	}

	return s.get(entityRecord.archetypeId, entityRecord.key), nil
}

func addComponentsToArchetype1[A ComponentInterface](world *World, entityId EntityId, archetype *Archetype, component A) error {
	storageA := getStorage[A](world)

	if storageA == nil {
		componentId := component.GetComponentId()
		return fmt.Errorf("no storage found for component %d", componentId)
	}

	// If the entity has no component, simply add it the archetype
	if entityRecord, ok := world.Entities[entityId]; !ok {
		world.setArchetype(entityRecord, archetype)
		setComponent(world, archetype.Id, component)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}

		setComponent(world, archetype.Id, component)
	}

	return nil
}

func addComponentsToArchetype2[A, B ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)

	if storageA == nil || storageB == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId()}
		return fmt.Errorf("no storage found for component %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func addComponentsToArchetype3[A, B, C ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B, componentC C) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)

	if storageA == nil || storageB == nil || storageC == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	setComponent(world, archetype.Id, componentC)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func addComponentsToArchetype4[A, B, C, D ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B, componentC C, componentD D) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	setComponent(world, archetype.Id, componentC)
	setComponent(world, archetype.Id, componentD)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func addComponentsToArchetype5[A, B, C, D, E ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B, componentC C, componentD D, componentE E) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)
	storageE := getStorage[E](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil || storageE == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId(), componentE.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	setComponent(world, archetype.Id, componentC)
	setComponent(world, archetype.Id, componentD)
	setComponent(world, archetype.Id, componentE)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func addComponentsToArchetype6[A, B, C, D, E, F ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B, componentC C, componentD D, componentE E, componentF F) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)
	storageE := getStorage[E](world)
	storageF := getStorage[F](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil || storageE == nil || storageF == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId(), componentE.GetComponentId(), componentF.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	setComponent(world, archetype.Id, componentC)
	setComponent(world, archetype.Id, componentD)
	setComponent(world, archetype.Id, componentE)
	setComponent(world, archetype.Id, componentF)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func addComponentsToArchetype7[A, B, C, D, E, F, G ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B, componentC C, componentD D, componentE E, componentF F, componentG G) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)
	storageE := getStorage[E](world)
	storageF := getStorage[F](world)
	storageG := getStorage[G](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil || storageE == nil || storageF == nil || storageG == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId(), componentE.GetComponentId(), componentF.GetComponentId(), componentG.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	setComponent(world, archetype.Id, componentC)
	setComponent(world, archetype.Id, componentD)
	setComponent(world, archetype.Id, componentE)
	setComponent(world, archetype.Id, componentF)
	setComponent(world, archetype.Id, componentG)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func addComponentsToArchetype8[A, B, C, D, E, F, G, H ComponentInterface](world *World, entityRecord EntityRecord, archetype *Archetype, componentA A, componentB B, componentC C, componentD D, componentE E, componentF F, componentG G, componentH H) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)
	storageE := getStorage[E](world)
	storageF := getStorage[F](world)
	storageG := getStorage[G](world)
	storageH := getStorage[H](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil || storageE == nil || storageF == nil || storageG == nil || storageH == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId(), componentE.GetComponentId(), componentF.GetComponentId(), componentG.GetComponentId(), componentH.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	setComponent(world, archetype.Id, componentA)
	setComponent(world, archetype.Id, componentB)
	setComponent(world, archetype.Id, componentC)
	setComponent(world, archetype.Id, componentD)
	setComponent(world, archetype.Id, componentE)
	setComponent(world, archetype.Id, componentF)
	setComponent(world, archetype.Id, componentG)
	setComponent(world, archetype.Id, componentH)
	world.setArchetype(entityRecord, archetype)

	return nil
}

func moveComponentsToArchetype(world *World, entityRecord EntityRecord, oldArchetype *Archetype, archetype *Archetype) int {
	var key, lastEntityKey int

	for _, componentId := range oldArchetype.Type {
		if !slices.Contains(archetype.Type, componentId) {
			continue
		}

		s, _ := world.getStorageForComponentId(componentId)

		if s != nil {
			key = s.copy(oldArchetype.Id, archetype.Id, entityRecord.key)
			s.moveLastToKey(oldArchetype.Id, entityRecord.key)
		}
	}

	lastEntityKey = len(oldArchetype.entities) - 1

	lastEntityId := oldArchetype.entities[lastEntityKey]
	lastEntity := world.Entities[lastEntityId]
	lastEntity.key = entityRecord.key
	world.Entities[lastEntityId] = lastEntity

	oldArchetype.entities[entityRecord.key] = lastEntityId
	oldArchetype.entities = oldArchetype.entities[:lastEntityKey]

	return key
}

func setComponent[T ComponentInterface](world *World, archetypeId ArchetypeId, component T) int {
	s := getStorage[T](world)

	return s.add(archetypeId, component)
}
