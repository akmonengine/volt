package volt

import (
	"slices"
	"testing"
)

const TEST_ENTITY_NUMBER = 10000

func TestCreateWorld(t *testing.T) {
	initialCapacity := 16
	world := CreateWorld(initialCapacity)
	if world == nil {
		t.Errorf("CreateWorld() returned nil value")

		return
	}
	if len(world.archetypes) != 1 || world.archetypes[0].Id != 0 {
		t.Errorf("CreateWorld() did not generate default archetype")
	}
}

func TestWorld_CreateEntity(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity()
	}

	// Check if the entities all exist in the world
	if len(world.entities) != TEST_ENTITY_NUMBER {
		t.Errorf("Number of entities created invalid")
	}
}

func TestCreateEntityWithComponents2(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	entityId, err := CreateEntityWithComponents2(world, testComponent1{}, testComponent2{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
}

func TestCreateEntityWithComponents3(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})

	entityId, err := CreateEntityWithComponents3(world, testComponent1{}, testComponent2{}, testComponent3{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent3](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent3 for entityId %d", entityId)
	}
}

func TestCreateEntityWithComponents4(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})

	entityId, err := CreateEntityWithComponents4(world, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent3](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent3 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent4](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent4 for entityId %d", entityId)
	}
}

func TestCreateEntityWithComponents5(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})

	entityId, err := CreateEntityWithComponents5(world, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent3](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent3 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent4](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent4 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent5](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent5 for entityId %d", entityId)
	}
}

func TestCreateEntityWithComponents6(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{})

	entityId, err := CreateEntityWithComponents6(world, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent3](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent3 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent4](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent4 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent5](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent5 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent6](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent6 for entityId %d", entityId)
	}
}

func TestCreateEntityWithComponents7(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{})

	entityId, err := CreateEntityWithComponents7(world, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent3](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent3 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent4](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent4 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent5](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent5 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent6](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent6 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent7](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent7 for entityId %d", entityId)
	}
}

func TestCreateEntityWithComponents8(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{})

	entityId, err := CreateEntityWithComponents8(world, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{}, testComponent8{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if component := GetComponent[testComponent1](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent1 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent2](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent2 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent3](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent3 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent4](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent4 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent5](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent5 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent6](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent6 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent7](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent7 for entityId %d", entityId)
	}
	if component := GetComponent[testComponent8](world, entityId); component == nil {
		t.Errorf("Could not find component testComponent8 for entityId %d", entityId)
	}
}

func TestWorld_RemoveEntity(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity()
	}
	// Remove first, last, and a middle one entity from the world
	world.RemoveEntity(entities[0])
	world.RemoveEntity(entities[TEST_ENTITY_NUMBER/2])
	world.RemoveEntity(entities[TEST_ENTITY_NUMBER-1])

	// Check if the entities are correctly removed of the world
	for _, id := range []EntityId{0, TEST_ENTITY_NUMBER / 2, TEST_ENTITY_NUMBER - 1} {
		if !slices.Contains(world.pool.ids, id) {
			t.Errorf("Entity %d was not removed", entities[id])
		}
	}
}

func TestWorld_Count(t *testing.T) {
	world := CreateWorld(1024)

	if world.Count() != 1024 {
		t.Errorf("world.Count should return 0 if the world is empty")
	}

	entities := make([]EntityId, TEST_ENTITY_NUMBER)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity()
	}

	if world.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("world.Count returned %d after inserting %d entities", world.Count(), TEST_ENTITY_NUMBER)
	}
}
