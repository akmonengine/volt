package volt

import (
	"testing"
)

func TestComponentConfig_addComponent(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{BuilderFn: func(component any, configuration any) {}})

	entityId := world.CreateEntity()
	componentRegistry, _ := world.getConfigByComponentId(testComponent1Id)
	err := componentRegistry.addComponent(world, entityId, testComponent1Configuration{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	component := GetComponent[testComponent1](world, entityId)
	if component == nil {
		t.Errorf("Component %d was not added to the entity %d", testComponent1Id, entityId)
		return
	}
}

func TestComponentConfig_builderFn(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{BuilderFn: func(component any, configuration any) {
		conf := configuration.(testComponent1Configuration)
		testTransformComponent := component.(*testComponent1)

		testTransformComponent.x = conf.x
		testTransformComponent.y = conf.y
		testTransformComponent.z = conf.z
	}})

	entityId := world.CreateEntity()
	componentRegistry, _ := world.getConfigByComponentId(testComponent1Id)
	err := componentRegistry.addComponent(world, entityId, testComponent1Configuration{
		testComponent{x: 1.0, y: 2.0, z: 3.0},
	})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	component := GetComponent[testComponent1](world, entityId)

	if component.x != 1.0 || component.y != 2.0 || component.z != 3.0 {
		t.Errorf("Component %d is not properly built", testComponent1Id)
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
