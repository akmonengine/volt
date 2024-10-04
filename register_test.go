package volt

import (
	"testing"
)

func TestComponentConfig_addComponent(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId, BuilderFn: func(component any, configuration any) {}})

	entityId := world.CreateEntity("entity")
	componentRegistry, _ := world.getConfigByComponentId(testTransformId)
	err := componentRegistry.addComponent(world, entityId, testTransformConfiguration{})

	if err != nil {
		t.Errorf(err.Error())
	}
	component := GetComponent[testTransform](world, entityId)
	if component == nil {
		t.Errorf("Component %d was not added to the entity %d", testTransformId, entityId)
		return
	}
}

func TestComponentConfig_builderFn(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId, BuilderFn: func(component any, configuration any) {
		conf := configuration.(testTransformConfiguration)
		testTransformComponent := component.(*testTransform)

		testTransformComponent.x = conf.x
		testTransformComponent.y = conf.y
		testTransformComponent.z = conf.z
	}})

	entityId := world.CreateEntity("entity")
	componentRegistry, _ := world.getConfigByComponentId(testTransformId)
	err := componentRegistry.addComponent(world, entityId, testTransformConfiguration{
		x: 1.0,
		y: 2.0,
		z: 3.0,
	})

	if err != nil {
		t.Errorf(err.Error())
	}
	component := GetComponent[testTransform](world, entityId)

	if component.x != 1.0 || component.y != 2.0 || component.z != 3.0 {
		t.Errorf("Component %d is not properly built", testTransformId)
	}
}

func TestComponentConfig_getComponentId(t *testing.T) {

}

func TestComponentConfig_setComponent(t *testing.T) {

}

func TestRegisterComponent(t *testing.T) {

}

func TestWorld_getConfigByComponentId(t *testing.T) {

}
