package volt

import (
	"fmt"
)

type ComponentConfigInterface interface {
	builderFn(component any, configuration any)
	getComponentId() ComponentId
	setComponent(component any)
	addComponent(world *World, entityId EntityId, configuration any) error
}
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
type ComponentBuilder func(component any, configuration any)

func RegisterComponent[T ComponentInterface](world *World, config ComponentConfigInterface) {
	var t T
	if world.ComponentsRegistry == nil {
		world.ComponentsRegistry = make(ComponentsRegister)
	}

	config.setComponent(t)
	world.ComponentsRegistry[t.GetComponentId()] = config
	getStorage[T](world)
}

func (world *World) getConfigByComponentId(componentId ComponentId) (ComponentConfigInterface, error) {
	for _, config := range world.ComponentsRegistry {
		if config.getComponentId() == componentId {
			return config, nil
		}
	}

	return nil, fmt.Errorf("componentConfiguration not found for %d", componentId)
}
