package volt

import (
	"fmt"
	"slices"
	"testing"
)

func TestCreateQuery1(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId})

	query := CreateQuery1[testTransform](world, []OptionalComponent{})

	if len(query.componentsIds) != 1 {
		t.Errorf("query should have 1 component")
	}
	if query.componentsIds[0] != testTransformId {
		t.Errorf("query should contain ComponentId %d", testTransformId)
	}
}

func TestQuery1_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId})
	RegisterComponent[testTag](world, &ComponentConfig[testTag]{ID: testTagId})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			AddComponent[testTransform](world, entityId, testTransform{})
		} else {
			AddComponents2[testTransform, testTag](world, entityId, testTransform{}, testTag{})
		}
	}

	query := CreateQuery1[testTransform](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery1_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		AddComponent[testTransform](world, entityId, testTransform{})
	}

	query := CreateQuery1[testTransform](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery1_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		AddComponent[testTransform](world, entityId, testTransform{})
	}

	query := CreateQuery1[testTransform](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}
