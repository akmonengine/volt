package volt

import (
	"fmt"
)

// ComponentConfigInterface is the interface
// defining the method required to create a new Component.
type ComponentConfigInterface interface {
	builderFn(component any, configuration any)
	getComponentId() ComponentId
	setComponent(component any)
	addComponent(world *World, entityId EntityId, configuration any) error
}

// Configuration for a component T.
//
// BuilderFn defines the function called to set a new component.
type ComponentConfig[T ComponentInterface] struct {
	id        ComponentId
	BuilderFn ComponentBuilder
	component T
}

func (componentConfig *ComponentConfig[T]) getComponentId() ComponentId {
	return componentConfig.id
}

func (componentConfig *ComponentConfig[T]) setComponent(component any) {
	componentConfig.component = component.(T)
	componentConfig.id = component.(T).GetComponentId()
}

func (componentConfig *ComponentConfig[T]) addComponent(world *World, entityId EntityId, configuration any) error {
	var t T
	componentConfig.builderFn(&t, configuration)

	archetype := world.getNextArchetype(entityId, componentConfig.id)
	err := addComponentsToArchetype1[T](world, entityId, archetype, t)

	return err
}

func (componentConfig *ComponentConfig[T]) builderFn(component any, configuration any) {
	if componentConfig.BuilderFn != nil {
		componentConfig.BuilderFn(component.(*T), configuration)
	}
}

type ComponentsRegister map[ComponentId]ComponentConfigInterface

// ComponentBuilder is the function called to set the properties of a given component.
//
// A type assertion is required on component and configuration parameters.
type ComponentBuilder func(component any, configuration any)

// RegisterComponent adds a component T to the registry of the given World.
//
// Once the component is registered, it can be added to an entity.
func RegisterComponent[T ComponentInterface](world *World, config ComponentConfigInterface) {
	var t T
	if world.componentsRegistry == nil {
		world.componentsRegistry = make(ComponentsRegister)
	}

	config.setComponent(t)
	world.componentsRegistry[t.GetComponentId()] = config
	getStorage[T](world)
}

func (world *World) getConfigByComponentId(componentId ComponentId) (ComponentConfigInterface, error) {
	for _, config := range world.componentsRegistry {
		if config.getComponentId() == componentId {
			return config, nil
		}
	}

	return nil, fmt.Errorf("componentConfiguration not found for %d", componentId)
}
