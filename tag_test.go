package volt

import (
	"fmt"
	"slices"
	"testing"
)

const (
	TAG_1 = iota + TAGS_INDICES
	TAG_2
)

func TestAddTag(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		AddComponent[testComponent1](world, entities[i], testComponent1{})
		world.AddTag(TAG_1, entities[i])
		world.AddTag(TAG_2, entities[i])
		world.RemoveTag(TAG_1, entities[i])
	}

	query := CreateQuery1[testComponent1](world, QueryConfiguration{Tags: []TagId{TAG_2}})
	if len(query.componentsIds) != 1 {
		t.Errorf("query should have 1 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)

	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestHasTag(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		AddComponent[testComponent1](world, entities[i], testComponent1{})
		world.AddTag(TAG_1, entities[i])
		world.AddTag(TAG_2, entities[i])
		world.RemoveTag(TAG_1, entities[i])
	}

	query := CreateQuery1[testComponent1](world, QueryConfiguration{Tags: []TagId{}})
	if len(query.componentsIds) != 1 {
		t.Errorf("query should have 1 component")
	}

	results := query.Foreach(nil)
	for result := range results {
		if !world.HasTag(TAG_2, result.EntityId) {
			t.Errorf("entity %d should have the tag %d", result.EntityId, TAG_2)
		}
	}
}

func TestRemoveTag(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		AddComponent[testComponent1](world, entities[i], testComponent1{})
		world.AddTag(TAG_1, entities[i])
		world.AddTag(TAG_2, entities[i])
		world.RemoveTag(TAG_1, entities[i])
	}

	query := CreateQuery1[testComponent1](world, QueryConfiguration{Tags: []TagId{TAG_1}})
	if len(query.componentsIds) != 1 {
		t.Errorf("query should have 1 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != 0 {
		t.Errorf("entities should not have the tag %d", TAG_1)
	}
}
