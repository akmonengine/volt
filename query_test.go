package volt

import (
	"fmt"
	"slices"
	"testing"
)

func TestCreateQuery1(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})

	query := CreateQuery1[testComponent1](world, []OptionalComponent{})

	if len(query.componentsIds) != 1 {
		t.Errorf("query should have 1 component")
	}
	if query.componentsIds[0] != testComponent1Id {
		t.Errorf("query should contain ComponentId %d", testComponent1Id)
	}
}

func TestQuery1_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponent[testComponent1](world, entityId, testComponent1{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents2[testComponent1, testComponent2](world, entityId, testComponent1{}, testComponent2{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery1[testComponent1](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery1_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		err := AddComponent[testComponent1](world, entityId, testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery1[testComponent1](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery1_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponent[testComponent1](world, entityId, testComponent1{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery1[testComponent1](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery2(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})

	query := CreateQuery2[testComponent1, testComponent2](world, []OptionalComponent{})

	if len(query.componentsIds) != 2 {
		t.Errorf("query should have 2 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) {
		t.Errorf("query should contain ComponentIds %d and %d", testComponent1Id, testComponent2Id)
	}
}

func TestQuery2_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponents3[testComponent1, testComponent2, testComponent3](world, entityId, testComponent1{}, testComponent2{}, testComponent3{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents2[testComponent1, testComponent2](world, entityId, testComponent1{}, testComponent2{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery2[testComponent1, testComponent2](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery2_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents2[testComponent1, testComponent2](world, entityId, testComponent1{}, testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery2[testComponent1, testComponent2](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery2_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents2[testComponent1, testComponent2](world, entityId, testComponent1{}, testComponent2{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery2[testComponent1, testComponent2](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery3(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})

	query := CreateQuery3[testComponent1, testComponent2, testComponent3](world, []OptionalComponent{})

	if len(query.componentsIds) != 3 {
		t.Errorf("query should have 3 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) || !slices.Contains(query.componentsIds, testComponent3Id) {
		t.Errorf("query should contain ComponentIds %d, %d and %d", testComponent1Id, testComponent2Id, testComponent3Id)
	}
}

func TestQuery3_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponents4[testComponent1, testComponent2, testComponent3, testComponent4](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents3[testComponent1, testComponent2, testComponent3](world, entityId, testComponent1{}, testComponent2{}, testComponent3{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery3[testComponent1, testComponent2, testComponent3](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery3_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents3[testComponent1, testComponent2, testComponent3](world, entityId, testComponent1{}, testComponent2{}, testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery3[testComponent1, testComponent2, testComponent3](world, []OptionalComponent{})
	count := query.Count()
	if count != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery3_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents3[testComponent1, testComponent2, testComponent3](world, entityId, testComponent1{}, testComponent2{}, testComponent3{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery3[testComponent1, testComponent2, testComponent3](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery4(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})

	query := CreateQuery4[testComponent1, testComponent2, testComponent3, testComponent4](world, []OptionalComponent{})

	if len(query.componentsIds) != 4 {
		t.Errorf("query should have 4 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) || !slices.Contains(query.componentsIds, testComponent3Id) || !slices.Contains(query.componentsIds, testComponent4Id) {
		t.Errorf("query should contain ComponentIds %d, %d, %d and %d", testComponent1Id, testComponent2Id, testComponent3Id, testComponent4Id)
	}
}

func TestQuery4_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponents5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents4[testComponent1, testComponent2, testComponent3, testComponent4](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery4[testComponent1, testComponent2, testComponent3, testComponent4](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery4_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents4[testComponent1, testComponent2, testComponent3, testComponent4](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery4[testComponent1, testComponent2, testComponent3, testComponent4](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery4_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents4[testComponent1, testComponent2, testComponent3, testComponent4](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery4[testComponent1, testComponent2, testComponent3, testComponent4](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery5(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})

	query := CreateQuery5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, []OptionalComponent{})

	if len(query.componentsIds) != 5 {
		t.Errorf("query should have 5 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) || !slices.Contains(query.componentsIds, testComponent3Id) || !slices.Contains(query.componentsIds, testComponent4Id) || !slices.Contains(query.componentsIds, testComponent5Id) {
		t.Errorf("query should contain ComponentIds %d, %d, %d, %d and %d", testComponent1Id, testComponent2Id, testComponent3Id, testComponent4Id, testComponent5Id)
	}
}

func TestQuery5_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponents6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery5_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery5_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery5[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery6(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})

	query := CreateQuery6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, []OptionalComponent{})

	if len(query.componentsIds) != 6 {
		t.Errorf("query should have 6 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) || !slices.Contains(query.componentsIds, testComponent3Id) || !slices.Contains(query.componentsIds, testComponent4Id) || !slices.Contains(query.componentsIds, testComponent5Id) || !slices.Contains(query.componentsIds, testComponent6Id) {
		t.Errorf("query should contain ComponentIds %d, %d, %d, %d, %d and %d", testComponent1Id, testComponent2Id, testComponent3Id, testComponent4Id, testComponent5Id, testComponent6Id)
	}
}

func TestQuery6_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponents7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery6_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery6_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery6[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery7(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})

	query := CreateQuery7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, []OptionalComponent{})

	if len(query.componentsIds) != 7 {
		t.Errorf("query should have 7 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) || !slices.Contains(query.componentsIds, testComponent3Id) || !slices.Contains(query.componentsIds, testComponent4Id) || !slices.Contains(query.componentsIds, testComponent5Id) || !slices.Contains(query.componentsIds, testComponent6Id) || !slices.Contains(query.componentsIds, testComponent7Id) {
		t.Errorf("query should contain ComponentIds %d, %d, %d, %d, %d, %d and %d", testComponent1Id, testComponent2Id, testComponent3Id, testComponent4Id, testComponent5Id, testComponent6Id, testComponent7Id)
	}
}

func TestQuery7_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{ID: testComponent8Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		if i%2 == 0 {
			err := AddComponents8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{}, testComponent8{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		} else {
			err := AddComponents7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{})
			if err != nil {
				t.Errorf("%s", err.Error())
			}
		}
	}

	query := CreateQuery7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 2 {
		t.Errorf("query should have 2 archetypes")
	}
}

func TestQuery7_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery7_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery7[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7](world, []OptionalComponent{})
	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}

func TestCreateQuery8(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{ID: testComponent8Id})

	query := CreateQuery8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, []OptionalComponent{})

	if len(query.componentsIds) != 8 {
		t.Errorf("query should have 8 components")
	}
	if !slices.Contains(query.componentsIds, testComponent1Id) || !slices.Contains(query.componentsIds, testComponent2Id) || !slices.Contains(query.componentsIds, testComponent3Id) || !slices.Contains(query.componentsIds, testComponent4Id) || !slices.Contains(query.componentsIds, testComponent5Id) || !slices.Contains(query.componentsIds, testComponent6Id) || !slices.Contains(query.componentsIds, testComponent7Id) || !slices.Contains(query.componentsIds, testComponent8Id) {
		t.Errorf("query should contain ComponentIds %d, %d, %d, %d, %d, %d, %d and %d", testComponent1Id, testComponent2Id, testComponent3Id, testComponent4Id, testComponent5Id, testComponent6Id, testComponent7Id, testComponent8Id)
	}
}

func TestQuery8_filter(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{ID: testComponent8Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))

		err := AddComponents8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{}, testComponent8{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, []OptionalComponent{})
	archetypes := query.filter()
	if len(archetypes) != 1 {
		t.Errorf("query should have 1 archetype")
	}
}

func TestQuery8_Count(t *testing.T) {
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{ID: testComponent8Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		err := AddComponents8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{}, testComponent8{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, []OptionalComponent{})
	if query.Count() != TEST_ENTITY_NUMBER {
		t.Errorf("query should count %d entities", TEST_ENTITY_NUMBER)
	}
}

func TestQuery8_Foreach(t *testing.T) {
	var entities []EntityId
	world := CreateWorld(TEST_ENTITY_NUMBER)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{ID: testComponent2Id})
	RegisterComponent[testComponent3](world, &ComponentConfig[testComponent3]{ID: testComponent3Id})
	RegisterComponent[testComponent4](world, &ComponentConfig[testComponent4]{ID: testComponent4Id})
	RegisterComponent[testComponent5](world, &ComponentConfig[testComponent5]{ID: testComponent5Id})
	RegisterComponent[testComponent6](world, &ComponentConfig[testComponent6]{ID: testComponent6Id})
	RegisterComponent[testComponent7](world, &ComponentConfig[testComponent7]{ID: testComponent7Id})
	RegisterComponent[testComponent8](world, &ComponentConfig[testComponent8]{ID: testComponent8Id})

	for i := 0; i < TEST_ENTITY_NUMBER; i++ {
		entityId := world.CreateEntity(fmt.Sprint(i))
		entities = append(entities, entityId)

		err := AddComponents8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, entityId, testComponent1{}, testComponent2{}, testComponent3{}, testComponent4{}, testComponent5{}, testComponent6{}, testComponent7{}, testComponent8{})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	query := CreateQuery8[testComponent1, testComponent2, testComponent3, testComponent4, testComponent5, testComponent6, testComponent7, testComponent8](world, []OptionalComponent{})

	for result := range query.Foreach(nil) {
		if !slices.Contains(entities, result.EntityId) {
			t.Errorf("query should return EntityId %d in Foreach iterator", result.EntityId)
		}
	}
}
