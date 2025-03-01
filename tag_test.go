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

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
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
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}

	err = world.AddTag(TAG_2, entities[0])
	if err == nil {
		t.Errorf("AddTag should return an error due to the entity %d already having the tag %d", entities[0], TAG_2)
	}

	err = world.AddTag(TAGS_INDICES-1, entities[0])
	if err == nil {
		t.Errorf("AddTag should return an error due to the TagId %d being lower than the range allowed", TAGS_INDICES-1)
	}
}

func TestHasTag(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
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

	if world.HasTag(TAG_2, 0) {
		t.Errorf("entity %d does not exist, it should not have the tag %d", 0, TAG_2)
	}
}

func TestRemoveTag(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery1[testComponent1](world, QueryConfiguration{Tags: []TagId{TAG_1}})
	if len(query.componentsIds) != 1 {
		t.Errorf("query should have 1 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != 0 {
		t.Errorf("entities should not have the tag %d", TAG_1)
	}

	err := world.RemoveTag(TAG_1, 0)
	if err == nil {
		t.Errorf("we should not be able to remove the tag %d from an unexisting entity %d", TAG_1, entities[0])
	}

	err = world.RemoveTag(TAG_1, entities[0])
	if err == nil {
		t.Errorf("we should not be able to remove the tag %d from the entity %d that does not own it", TAG_1, entities[0])
	}
}

func TestTag2(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery2[testComponent1, testComponent2](world, QueryConfiguration{Tags: []TagId{TAG_2}})
	if len(query.componentsIds) != 2 {
		t.Errorf("query should have 2 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)

	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestTag3(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent3](world, entities[i], testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery3[testComponent1, testComponent2, testComponent3](world, QueryConfiguration{Tags: []TagId{TAG_2}})
	if len(query.componentsIds) != 3 {
		t.Errorf("query should have 3 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)

	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestTag4(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent3](world, entities[i], testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent4](world, entities[i], testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery4[testComponent1, testComponent2, testComponent3, testComponent4](world, QueryConfiguration{Tags: []TagId{TAG_2}})

	if len(query.componentsIds) != 4 {
		t.Errorf("query should have 4 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)
	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestTag5(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent3](world, entities[i], testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent4](world, entities[i], testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent5](world, entities[i], testComponent5{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, QueryConfiguration{Tags: []TagId{TAG_2}})

	if len(query.componentsIds) != 5 {
		t.Errorf("query should have 5 component")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)
	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestTag6(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent3](world, entities[i], testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent4](world, entities[i], testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent5](world, entities[i], testComponent5{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent6](world, entities[i], testComponent6{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, QueryConfiguration{Tags: []TagId{TAG_2}})

	if len(query.componentsIds) != 6 {
		t.Errorf("query should have 6 components")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)
	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestTag7(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent3](world, entities[i], testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent4](world, entities[i], testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent5](world, entities[i], testComponent5{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent6](world, entities[i], testComponent6{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent7](world, entities[i], testComponent7{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, QueryConfiguration{Tags: []TagId{TAG_2}})

	if len(query.componentsIds) != 7 {
		t.Errorf("query should have 7 components")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)
	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}

func TestTag8(t *testing.T) {
	entities := make([]EntityId, TEST_ENTITY_NUMBER)
	world := CreateWorld(1024)

	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entities[i] = world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entities[i], testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent2](world, entities[i], testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent3](world, entities[i], testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent4](world, entities[i], testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent5](world, entities[i], testComponent5{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent6](world, entities[i], testComponent6{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent7](world, entities[i], testComponent7{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = AddComponent[testComponent8](world, entities[i], testComponent8{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.AddTag(TAG_2, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		err = world.RemoveTag(TAG_1, entities[i])
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, QueryConfiguration{Tags: []TagId{TAG_2}})

	if len(query.componentsIds) != 8 {
		t.Errorf("query should have 8 components")
	}

	results := query.Foreach(nil)
	if len(slices.Collect(results)) != len(entities) {
		t.Errorf("entities should have the tag %d", TAG_2)
	}

	_, err := world.GetComponent(entities[0], TAG_2)
	if err == nil {
		t.Errorf("The tag %d should not be returned as a Component from world.GetComponent", TAG_2)
	}
}
