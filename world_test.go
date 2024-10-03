package volt

import (
	"fmt"
	"testing"
)

const TEST_ENTITY_NUMBER = 10000

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
