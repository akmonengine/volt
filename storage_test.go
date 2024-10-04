package volt

import (
	"testing"
)

func Test_getStorage(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})

	s1 := getStorage[testComponent1](world)
	if s1 == nil {
		t.Errorf("RegisterComponent did not generate storage for ComponentId %d", testComponent1Id)
	}
	s2 := getStorage[testComponent2](world)
	if s2 != nil {
		t.Errorf("storage found for not registered ComponentId %d", testComponent2Id)
	}
}

func TestWorld_getStorageForComponentId(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{ID: testComponent1Id})

	s, err := world.getStorageForComponentId(testComponent1Id)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s == nil {
		t.Errorf("getStorageForComponentId() did not find storage for %d", testComponent1Id)
	}
	s, err = world.getStorageForComponentId(testComponent2Id)
	if err == nil || s != nil {
		t.Errorf("getStorageForComponentId() returned a storage for the not registered ComponentId %d", testComponent1Id)
	}
}
