package volt

import (
	"fmt"
	"testing"
)

const (
	testTransformId = iota
	testTagId
)

type testTransform struct {
	x, y, z int
}

type testTransformConfiguration struct {
	x, y, z int
}

func (t testTransform) GetComponentId() ComponentId {
	return testTransformId
}

type testTag struct {
	x, y, z int
}

func (t testTag) GetComponentId() ComponentId {
	return testTagId
}

func TestAddComponent(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId, BuilderFn: func(component any, configuration any) {}})
	RegisterComponent[testTag](world, &ComponentConfig[testTag]{ID: testTagId, BuilderFn: func(component any, configuration any) {}})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent(world, entities[i], testTransform{})
		if err != nil {
			t.Errorf("could not add component testTransform to entity %d: %s", entities[i], err)
		}

		err = AddComponent(world, entities[i], testTag{})
		if err != nil {
			t.Errorf("could not add component testTag to entity %d: %s", entities[i], err)
		}
	}

	for _, entityId := range entities {
		if !world.HasComponents(entityId, []ComponentId{testTransformId, testTagId}...) {
			t.Errorf("Expected 2 components for entity %d", entityId)
		}
	}
}

func TestConfigureComponent(t *testing.T) {
}

func TestGetComponent(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId, BuilderFn: func(component any, configuration any) {}})
	RegisterComponent[testTag](world, &ComponentConfig[testTag]{ID: testTagId, BuilderFn: func(component any, configuration any) {}})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
		AddComponent(world, entities[i], testTransform{})
		AddComponent(world, entities[i], testTag{})
	}

	for _, entityId := range entities {
		component := GetComponent[testTransform](world, entityId)

		if component == nil {
			t.Errorf("Could not get created component for entityId %d", entityId)
		}
	}
}

func TestRemoveComponent(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId, BuilderFn: func(component any, configuration any) {}})
	RegisterComponent[testTag](world, &ComponentConfig[testTag]{ID: testTagId, BuilderFn: func(component any, configuration any) {}})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
		AddComponent(world, entities[i], testTransform{x: i, y: i, z: i})
		AddComponent(world, entities[i], testTag{x: i, y: i, z: i})
	}

	// Remove the component only on odd entities. Otherwise we empty the archetype, it would not prove the indices in storage are perfectly handled
	for i, entityId := range entities {
		if i%2 == 0 {
			err := RemoveComponent[testTransform](world, entityId)
			if err != nil {
				t.Errorf("could not remove component testTransform to entity %d: %s", entityId, err)
			}
		}
	}

	for i, entityId := range entities {
		if i%2 == 0 {
			if world.HasComponents(entityId, []ComponentId{testTransformId}...) || !world.HasComponents(entityId, []ComponentId{testTagId}...) {
				t.Errorf("Expected 1 components for entity %d", entityId)
			}

			// if the indices in storage are messed up, then i & x would not match anymore
			if i != GetComponent[testTag](world, entityId).x {
				t.Errorf("x value does not equals its i index in testTag, the storage indices are probably shuffled")
			}
		} else {
			if !world.HasComponents(entityId, []ComponentId{testTransformId, testTagId}...) {
				t.Errorf("Expected 2 components for entity %d", entityId)
			}

			if i != GetComponent[testTransform](world, entityId).x {
				t.Errorf("x value does not equals its i index in testTransform, the storage indices are probably shuffled")
			}
			if i != GetComponent[testTag](world, entityId).x {
				t.Errorf("x value does not equals its i index in testTag, the storage indices are probably shuffled")
			}
		}
	}
}
