package volt

import (
	"fmt"
	"slices"
)

// ComponentInterface is the interface for all the Components.
//
// It wraps the GetComponentId method, that returns a Component identifier.
type ComponentInterface interface {
	GetComponentId() ComponentId
}

type ComponentIdConf struct {
	ComponentId
	conf any
}

func (world *World) getComponentsIds(components ...ComponentInterface) []ComponentId {
	componentsIds := make([]ComponentId, len(components))

	for i, component := range components {
		componentsIds[i] = component.GetComponentId()
	}

	return componentsIds
}

// ConfigureComponent configures a Component of type T using the build function related to it.
//
// The parameter conf contains all the data required for the configuration.
func ConfigureComponent[T ComponentInterface](world *World, conf any) T {
	var t T
	componentRegistry := world.componentsRegistry[t.GetComponentId()]

	componentRegistry.builderFn(&t, conf)

	return t
}

// AddComponent adds the component T to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has the component
//   - an internal error occurs
func AddComponent[T ComponentInterface](world *World, entityId EntityId, component T) error {
	componentId := component.GetComponentId()
	if world.HasComponents(entityId, componentId) {
		return fmt.Errorf("the entity %d already owns the component %d", entityId, componentId)
	}

	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	archetype := world.getNextArchetype(entityRecord, world.getComponentsIds(component)...)
	err := addComponentsToArchetype1(world, entityRecord, archetype, component)
	if err != nil {
		return fmt.Errorf("the component %d cannot be added to entity %d: %w", componentId, entityId, err)
	}

	world.componentAddedFn(entityId, componentId)

	return nil
}

// AddComponents2 adds the components A, B to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents2[A, B ComponentInterface](world *World, entityId EntityId, a A, b B) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents2(world, entityRecord, a, b)
}

func addComponents2[A, B ComponentInterface](world *World, entityRecord entityRecord, a A, b B) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId())

	entityId := entityRecord.Id
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

// AddComponents3 adds the components A, B, C to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents3[A, B, C ComponentInterface](world *World, entityId EntityId, a A, b B, c C) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents3(world, entityRecord, a, b, c)
}

func addComponents3[A, B, C ComponentInterface](world *World, entityRecord entityRecord, a A, b B, c C) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId())

	entityId := entityRecord.Id

	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId()})
	}

	err := addComponentsToArchetype3(world, entityRecord, archetype, a, b, c)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())

	return nil
}

// AddComponents4 adds the components A, B, C, D to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents4[A, B, C, D ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents4(world, entityRecord, a, b, c, d)
}

func addComponents4[A, B, C, D ComponentInterface](world *World, entityRecord entityRecord, a A, b B, c C, d D) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId())

	entityId := entityRecord.Id

	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId()})
	}

	err := addComponentsToArchetype4(world, entityRecord, archetype, a, b, c, d)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())

	return nil
}

// AddComponents5 adds the components A, B, C, D, E to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents5[A, B, C, D, E ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents5(world, entityRecord, a, b, c, d, e)
}

func addComponents5[A, B, C, D, E ComponentInterface](world *World, entityRecord entityRecord, a A, b B, c C, d D, e E) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId())

	entityId := entityRecord.Id

	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId()})
	}

	err := addComponentsToArchetype5(world, entityRecord, archetype, a, b, c, d, e)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())
	world.componentAddedFn(entityId, e.GetComponentId())

	return nil
}

// AddComponents6 adds the components A, B, C, D, E, F to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents6[A, B, C, D, E, F ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E, f F) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents6(world, entityRecord, a, b, c, d, e, f)
}

func addComponents6[A, B, C, D, E, F ComponentInterface](world *World, entityRecord entityRecord, a A, b B, c C, d D, e E, f F) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId())

	entityId := entityRecord.Id

	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId()})
	}

	err := addComponentsToArchetype6(world, entityRecord, archetype, a, b, c, d, e, f)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId()}, entityId, err)
	}

	world.componentAddedFn(entityId, a.GetComponentId())
	world.componentAddedFn(entityId, b.GetComponentId())
	world.componentAddedFn(entityId, c.GetComponentId())
	world.componentAddedFn(entityId, d.GetComponentId())
	world.componentAddedFn(entityId, e.GetComponentId())
	world.componentAddedFn(entityId, f.GetComponentId())

	return nil
}

// AddComponents7 adds the components A, B, C, D, E, F, G to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E, f F, g G) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents7(world, entityRecord, a, b, c, d, e, f, g)
}

func addComponents7[A, B, C, D, E, F, G ComponentInterface](world *World, entityRecord entityRecord, a A, b B, c C, d D, e E, f F, g G) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId())

	entityId := entityRecord.Id

	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId()})
	}

	err := addComponentsToArchetype7(world, entityRecord, archetype, a, b, c, d, e, f, g)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId()}, entityId, err)
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

// AddComponents8 adds the components A, B, C, D, E, F, G, H to the existing EntityId.
//
// It returns an error if:
//   - the entity does not exist
//   - the entity has one of the component
//   - an internal error occurs
//
// This solution is faster than an atomic solution.
func AddComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, entityId EntityId, a A, b B, c C, d D, e E, f F, g G, h H) error {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return fmt.Errorf("entity %v does not exist", entityId)
	}

	return addComponents8(world, entityRecord, a, b, c, d, e, f, g, h)
}

func addComponents8[A, B, C, D, E, F, G, H ComponentInterface](world *World, entityRecord entityRecord, a A, b B, c C, d D, e E, f F, g G, h H) error {
	archetype := world.getNextArchetype(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId())

	entityId := entityRecord.Id

	if world.hasComponents(entityRecord, a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId()) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId()})
	}

	err := addComponentsToArchetype8(world, entityRecord, archetype, a, b, c, d, e, f, g, h)
	if err != nil {
		return fmt.Errorf("the components %d cannot be added to entity %d: %w", []ComponentId{a.GetComponentId(), b.GetComponentId(), c.GetComponentId(), d.GetComponentId(), e.GetComponentId(), f.GetComponentId(), g.GetComponentId(), h.GetComponentId()}, entityId, err)
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

// AddComponent adds the component with ComponentId to the EntityId.
//
// This non-generic version is adapted for when generics are not available, though might be slower.
// It returns an error if:
//   - the entity already has the componentId
//   - the componentId is not registered in the World
//   - an internal error occurs
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

// AddComponents adds variadic components to the EntityId.
//
// This non-generic version is adapted for when generics are not available, though might be slower.
// It returns an error if:
//   - the entity already has the components Ids
//   - the componentsIds are not registered in the World
//   - an internal error occurs
func (world *World) AddComponents(entityId EntityId, componentsIdsConfs ...ComponentIdConf) error {
	var componentsIds []ComponentId
	for _, componentIdConf := range componentsIdsConfs {
		componentsIds = append(componentsIds, componentIdConf.ComponentId)
	}

	if world.HasComponents(entityId, componentsIds...) {
		return fmt.Errorf("the entity %d already owns the components %v", entityId, componentsIds)
	}

	for _, componentIdConf := range componentsIdsConfs {
		componentRegistry, err := world.getConfigByComponentId(componentIdConf.ComponentId)
		if err != nil {
			return err
		}

		err = componentRegistry.addComponent(world, entityId, componentIdConf.conf)
		if err != nil {
			return err
		}

		world.componentAddedFn(entityId, componentIdConf.ComponentId)
	}

	return nil
}

// RemoveComponent removes the component to EntityId.
//
// It returns an error if the EntityId does not have the component.
func RemoveComponent[T ComponentInterface](world *World, entityId EntityId) error {
	var t T
	componentId := t.GetComponentId()
	entityRecord := world.entities[entityId]

	if !world.hasComponents(entityRecord, componentId) {
		return fmt.Errorf("the entity %d doesn't own the component %d", entityId, componentId)
	}

	// Remove from the previous archetype
	s := getStorage[T](world)
	removeComponent(world, s, entityRecord, componentId)

	return nil
}

// RemoveComponent removes the component with ComponentId from the EntityId.
//
// This non-generic version is adapted for when generics are not available, though might be slower.
// It returns an error if:
//   - the entity does not have the component
//   - the ComponentId is not registered in the World
func (world *World) RemoveComponent(entityId EntityId, componentId ComponentId) error {
	entityRecord := world.entities[entityId]

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

func removeComponent(world *World, s storage, entityRecord entityRecord, componentId ComponentId) {
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

// HasComponents returns whether the entity has the given variadic list of ComponentId.
//
// It returns false if at least one ComponentId is not owned.
func (world *World) HasComponents(entityId EntityId, componentsIds ...ComponentId) bool {
	entityRecord, ok := world.entities[entityId]
	if !ok {
		return false
	}

	return world.hasComponents(entityRecord, componentsIds...)
}

func (world *World) hasComponents(entityRecord entityRecord, componentsIds ...ComponentId) bool {
	archetype := world.archetypes[entityRecord.archetypeId]
	for _, componentId := range componentsIds {
		if !slices.Contains(archetype.Type, componentId) {
			return false
		}
	}

	return true
}

// GetComponent returns a pointer to the component T owned by the entity.
//
// If the entity does not have the component, it returns nil
func GetComponent[T ComponentInterface](world *World, entityId EntityId) *T {
	s := getStorage[T](world)
	entityRecord := world.entities[entityId]

	if !s.hasArchetype(entityRecord.archetypeId) {
		return nil
	}

	return &s.archetypesComponentsEntities[entityRecord.archetypeId][entityRecord.key]
}

// GetComponent returns the component with ComponentId for EntityId.
//
// This non-generic version is adapted for when generics are not available,
// though might be slower and requires a type assertion.
// It returns an error if:
//   - the ComponentId is not registered in the World
//   - the entity does not have the component
func (world *World) GetComponent(entityId EntityId, componentId ComponentId) (any, error) {
	entityRecord := world.entities[entityId]
	s, err := world.getStorageForComponentId(componentId)
	if err != nil {
		return nil, err
	}

	if !s.hasArchetype(entityRecord.archetypeId) {
		return nil, fmt.Errorf("the entity %d doesn't own the component %d", entityId, componentId)
	}

	return s.get(entityRecord.archetypeId, entityRecord.key), nil
}

func addComponentsToArchetype1[A ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, component A) error {
	storageA := getStorage[A](world)

	if storageA == nil {
		componentId := component.GetComponentId()
		return fmt.Errorf("no storage found for component %d", componentId)
	}

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}
	storageA.add(archetype.Id, component)

	return nil
}

func addComponentsToArchetype2[A, B ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)

	if storageA == nil || storageB == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId()}
		return fmt.Errorf("no storage found for component %v", componentsIds)
	}

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)

	return nil
}

func addComponentsToArchetype3[A, B, C ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B, componentC C) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)

	if storageA == nil || storageB == nil || storageC == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)
	storageC.add(archetype.Id, componentC)

	return nil
}

func addComponentsToArchetype4[A, B, C, D ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B, componentC C, componentD D) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)
	storageC.add(archetype.Id, componentC)
	storageD.add(archetype.Id, componentD)

	return nil
}

func addComponentsToArchetype5[A, B, C, D, E ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B, componentC C, componentD D, componentE E) error {
	storageA := getStorage[A](world)
	storageB := getStorage[B](world)
	storageC := getStorage[C](world)
	storageD := getStorage[D](world)
	storageE := getStorage[E](world)

	if storageA == nil || storageB == nil || storageC == nil || storageD == nil || storageE == nil {
		componentsIds := []ComponentId{componentA.GetComponentId(), componentB.GetComponentId(), componentC.GetComponentId(), componentD.GetComponentId(), componentE.GetComponentId()}
		return fmt.Errorf("no storage found for components %v", componentsIds)
	}

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)
	storageC.add(archetype.Id, componentC)
	storageD.add(archetype.Id, componentD)
	storageE.add(archetype.Id, componentE)

	return nil
}

func addComponentsToArchetype6[A, B, C, D, E, F ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B, componentC C, componentD D, componentE E, componentF F) error {
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

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)
	storageC.add(archetype.Id, componentC)
	storageD.add(archetype.Id, componentD)
	storageE.add(archetype.Id, componentE)
	storageF.add(archetype.Id, componentF)

	return nil
}

func addComponentsToArchetype7[A, B, C, D, E, F, G ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B, componentC C, componentD D, componentE E, componentF F, componentG G) error {
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

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)
	storageC.add(archetype.Id, componentC)
	storageD.add(archetype.Id, componentD)
	storageE.add(archetype.Id, componentE)
	storageF.add(archetype.Id, componentF)
	storageG.add(archetype.Id, componentG)

	return nil
}

func addComponentsToArchetype8[A, B, C, D, E, F, G, H ComponentInterface](world *World, entityRecord entityRecord, archetype *archetype, componentA A, componentB B, componentC C, componentD D, componentE E, componentF F, componentG G, componentH H) error {
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

	// If the entity has no component, simply add it the archetype
	if entityRecord.archetypeId == 0 {
		world.setArchetype(entityRecord, archetype)
	} else {
		oldArchetype := world.getArchetype(entityRecord)
		if archetype.Id != oldArchetype.Id {
			moveComponentsToArchetype(world, entityRecord, oldArchetype, archetype)
			world.setArchetype(entityRecord, archetype)
		}
	}

	storageA.add(archetype.Id, componentA)
	storageB.add(archetype.Id, componentB)
	storageC.add(archetype.Id, componentC)
	storageD.add(archetype.Id, componentD)
	storageE.add(archetype.Id, componentE)
	storageF.add(archetype.Id, componentF)
	storageG.add(archetype.Id, componentG)
	storageH.add(archetype.Id, componentH)

	return nil
}

func moveComponentsToArchetype(world *World, entityRecord entityRecord, oldArchetype *archetype, archetype *archetype) int {
	var key, lastEntityKey int

	for _, componentId := range oldArchetype.Type {
		// tags are not movable
		if componentId >= TAGS_INDICES {
			continue
		}
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
	lastEntity := world.entities[lastEntityId]
	lastEntity.key = entityRecord.key
	world.entities[lastEntityId] = lastEntity

	oldArchetype.entities[entityRecord.key] = lastEntityId
	oldArchetype.entities = oldArchetype.entities[:lastEntityKey]

	return key
}
