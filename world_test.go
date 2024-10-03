package volt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	assert.True(t, len(world.Entities) == TEST_ENTITY_NUMBER)
	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		_, ok := world.Entities[entities[i]]
		assert.True(t, ok)
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

	// Check if the entities all exist in the world
	assert.True(t, len(world.Entities) == (TEST_ENTITY_NUMBER-3))
	assert.True(t, world.SearchEntity(fmt.Sprint(0)) == 0)
	assert.True(t, world.SearchEntity(fmt.Sprint(TEST_ENTITY_NUMBER/2)) == 0)
	assert.True(t, world.SearchEntity(fmt.Sprint(TEST_ENTITY_NUMBER-1)) == 0)
}

func TestWorld_SearchEntity(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
	}

	// Test searching for existing entities
	for entityName, entityId := range entities {
		assert.Equal(t, entityId, world.SearchEntity(fmt.Sprint(entityName)))
	}

	// Test searching for a non-existing entity
	assert.Zero(t, world.SearchEntity("nonexistent"))
}

func TestWorld_GetEntityName(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))
	}

	// Test the names for each entity
	for entityName, entityId := range entities {
		assert.Equal(t, fmt.Sprint(entityName), world.GetEntityName(entityId))
	}

	// Test if none entity return an empty name
	assert.Empty(t, world.GetEntityName(0))
}
