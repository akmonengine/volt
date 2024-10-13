package volt

import (
	"fmt"
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
		entities[i] = world.CreateEntity(fmt.Sprint(i))
	}

	// Check if the entities all exist in the world
	if len(world.Entities) != TEST_ENTITY_NUMBER {
		t.Errorf("Number of entities created invalid")
	}
	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		_, ok := world.Entities[entities[i]]

		if !ok {
			t.Errorf("Entity %d was not created properly", entities[i])
		}
	}
}

func TestCreateEntityWithComponents2(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	entityId, err := CreateEntityWithComponents2(world, "entity1", testComponent1{}, testComponent2{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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

	entityId, err := CreateEntityWithComponents3(world, "entity1", testComponent1{}, testComponent2{}, testComponent3{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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

	entityId, err := CreateEntityWithComponents4(world, "entity1", testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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

	entityId, err := CreateEntityWithComponents5(world, "entity1", testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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

	entityId, err := CreateEntityWithComponents6(world, "entity1", testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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

	entityId, err := CreateEntityWithComponents7(world, "entity1", testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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

	entityId, err := CreateEntityWithComponents8(world, "entity1", testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{}, testComponent8{})

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if id := world.SearchEntity("entity1"); id == 0 {
		t.Errorf("Could not find entityName %s", "entity1")
	}
	if _, ok := world.Entities[entityId]; !ok {
		t.Errorf("Could not find entityId %d", entityId)
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
		entities[i] = world.CreateEntity(fmt.Sprint(i))
	}
	// Remove first, last, and a middle one entity from the world
	world.RemoveEntity(entities[0])
	world.RemoveEntity(entities[TEST_ENTITY_NUMBER/2])
	world.RemoveEntity(entities[TEST_ENTITY_NUMBER-1])

	// Check the expected world size
	if len(world.Entities) != (TEST_ENTITY_NUMBER - 3) {
		t.Errorf("World size not valid after removal of entities")
	}

	// Check if the entities are correctly removed of the world
	for _, id := range []EntityId{0, TEST_ENTITY_NUMBER / 2, TEST_ENTITY_NUMBER - 1} {
		if world.SearchEntity(fmt.Sprint(0)) != 0 {
			t.Errorf("Entity %d was not removed", entities[id])
		}
	}
}

func TestWorld_SearchEntity(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
	}

	// Test searching for existing entities
	for entityName, entityId := range entities {
		if entityId != world.SearchEntity(fmt.Sprint(entityName)) {
			t.Errorf("SearchEntity does not return correct entityId for %s", fmt.Sprint(entityName))
		}
	}

	// Test searching for a non-existing entity
	if id := world.SearchEntity("nonexistent"); id != 0 {
		t.Errorf("world.SearchEntity returned id %d for a non existent entity", id)
	}
}

func TestWorld_GetEntityName(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
	}

	// Test the names for each entity
	for entityName, entityId := range entities {
		if fmt.Sprint(entityName) != world.GetEntityName(entityId) {
			t.Errorf("world.GetEntityName does not return correct value for id %d", entityId)
		}
	}

	// Test if none entity return an empty name
	if world.GetEntityName(0) != "" {
		t.Errorf("world.GetEntityName does not return empty string for entityId 0")
	}
}
