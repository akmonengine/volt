package volt

import (
	"fmt"
	"testing"
)

const (
	testComponent1Id = iota
	testComponent2Id
	testComponent3Id
	testComponent4Id
	testComponent5Id
	testComponent6Id
	testComponent7Id
	testComponent8Id
	testComponent9Id
)

type testComponent struct {
	x, y, z int
}

type testComponent1 struct {
	testComponent
}

type testComponent1Configuration struct {
	testComponent
}

func (t testComponent1) GetComponentId() ComponentId {
	return testComponent1Id
}

type testComponent2 struct {
	testComponent
}

func (t testComponent2) GetComponentId() ComponentId {
	return testComponent2Id
}

type testComponent3 struct {
}

func (t testComponent3) GetComponentId() ComponentId {
	return testComponent3Id
}

type testComponent4 struct {
}

func (t testComponent4) GetComponentId() ComponentId {
	return testComponent4Id
}

type testComponent5 struct {
}

func (t testComponent5) GetComponentId() ComponentId {
	return testComponent5Id
}

type testComponent6 struct {
}

func (t testComponent6) GetComponentId() ComponentId {
	return testComponent6Id
}

type testComponent7 struct {
}

func (t testComponent7) GetComponentId() ComponentId {
	return testComponent7Id
}

type testComponent8 struct {
}

func (t testComponent8) GetComponentId() ComponentId {
	return testComponent8Id
}

func TestAddComponent(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id, BuilderFn: func(component any, configuration any) {}})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id, BuilderFn: func(component any, configuration any) {}})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent(world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("could not add component testComponent1 to entity %d: %s", entities[i], err)
		}

		err = AddComponent(world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("could not add component testComponent2 to entity %d: %s", entities[i], err)
		}
	}

	for _, entityId := range entities {
		if !world.HasComponents(entityId, []ComponentId{testComponent1Id, testComponent2Id}...) {
			t.Errorf("Expected 2 components for entity %d", entityId)
		}
	}
}

func TestConfigureComponent(t *testing.T) {
}

func TestGetComponent(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id, BuilderFn: func(component any, configuration any) {}})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id, BuilderFn: func(component any, configuration any) {}})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
		err := AddComponent(world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		err = AddComponent(world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	for _, entityId := range entities {
		component := GetComponent[testComponent1](world, entityId)

		if component == nil {
			t.Errorf("Could not get created component for entityId %d", entityId)
		}
	}
}

func TestRemoveComponent(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id, BuilderFn: func(component any, configuration any) {}})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id, BuilderFn: func(component any, configuration any) {}})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
		err := AddComponent(world, entities[i], testComponent1{testComponent{x: i, y: i, z: i}})
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		err = AddComponent(world, entities[i], testComponent2{testComponent{x: i, y: i, z: i}})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	// Remove the component only on odd entities. Otherwise we empty the archetype, it would not prove the indices in storage are perfectly handled
	for i, entityId := range entities {
		if i%2 == 0 {
			err := RemoveComponent[testComponent1](world, entityId)
			if err != nil {
				t.Errorf("could not remove component testComponent1 to entity %d: %s", entityId, err)
			}
		}
	}

	for i, entityId := range entities {
		if i%2 == 0 {
			if world.HasComponents(entityId, []ComponentId{testComponent1Id}...) || !world.HasComponents(entityId, []ComponentId{testComponent2Id}...) {
				t.Errorf("Expected 1 components for entity %d", entityId)
			}

			// if the indices in storage are messed up, then i & x would not match anymore
			if i != GetComponent[testComponent2](world, entityId).x {
				t.Errorf("x value does not equals its i index in testComponent2, the storage indices are probably shuffled")
			}
		} else {
			if !world.HasComponents(entityId, []ComponentId{testComponent1Id, testComponent2Id}...) {
				t.Errorf("Expected 2 components for entity %d", entityId)
			}

			if i != GetComponent[testComponent1](world, entityId).x {
				t.Errorf("x value does not equals its i index in testComponent1, the storage indices are probably shuffled")
			}
			if i != GetComponent[testComponent2](world, entityId).x {
				t.Errorf("x value does not equals its i index in testComponent2, the storage indices are probably shuffled")
			}
		}
	}
}
